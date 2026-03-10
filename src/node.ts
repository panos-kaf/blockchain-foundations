import type { Socket } from 'net'
import type {  Message, HelloMessage, ErrorMessage, 
            PeersMessage, GetPeersMessage, GetObjectMessage, 
            IHaveObjectMessage, ObjectMessage
            } from './messages'
import * as m from './messages'
import { hashObject, objectManager } from './object'
import { getKnownPeers, appendPeers } from './peers'

export const connectedPeers = new Set<Peer>()

export const broadcast = (msg: string) => {
    for (const peer of connectedPeers) {
        peer.sendMessage(msg)
    }
}

export default class Peer {
    socket: Socket
    id: string
    buffer: string = ''
    handshakeComplete: boolean = false
    onDisconnect?: () => void
    onLog?: (msg: string) => void
    onLogErr?: (msg: string) => void

    handlers: Record<string, (message: Message) => Promise<void>> = {
        hello: async (msg) => this.handleHello(msg as HelloMessage),
        getpeers: async (_m) => this.handleGetPeers(),
        peers: async (msg) => this.handlePeers(msg as PeersMessage),
        error: async (msg) => this.handleError(msg as ErrorMessage),
        getobject: async (msg) => this.handleGetObject(msg as GetObjectMessage),
        ihaveobject: async (msg) => this.handleIHaveObject(msg as IHaveObjectMessage),
        object: async (msg) => this.handleObject(msg as ObjectMessage),
    }

    constructor(socket: Socket, private role: 'server' | 'client', onLog?: (msg: string) => void, onLogErr?: (msg: string) => void) {
        this.socket = socket
        this.id = `${socket.remoteAddress}:${socket.remotePort}`
        this.onLog = onLog
        this.onLogErr = onLogErr
        connectedPeers.add(this) // Add to global set of peers
        this.initializeSocket()
    }

    private log(msg: string) {
        if (this.onLog) {
            this.onLog(msg)
        } else {
            console.log(`[${this.id}] ${msg}`)
        }
    }
    private logErr(msg: string) {
        if (this.onLogErr) {
            this.onLogErr(msg)
        } else {
            console.error(`[${this.id}] ${msg}`)
        }
    }

    /** Call after constructing to start the handshake */
    greet() {
        this.sendMessage(m.StaticHello)
        this.sendMessage(m.StaticGetPeers)
    }

    private initializeSocket() {
        this.socket.on('data', (data) => this.handleStream(data.toString()))
        this.socket.on('error', (err) => this.logErr(`Error: ${err}`))
        this.socket.on('close', () => {
            this.log(`Disconnected`)
            connectedPeers.delete(this)
            this.onDisconnect?.()
        })
    }

    private handleStream(data: string) {
        this.buffer += data
        const messages = this.buffer.split('\n')

        while (messages.length > 1) {
            const raw = messages.shift()?.trim() ?? ''
            if (!raw) continue
            this.handleMessage(raw)
        }

        this.buffer = messages[0] ?? ''
    }

    private async handleMessage(raw: string) {
        let parsed: unknown
        try {
            parsed = JSON.parse(raw)
        } catch {
            this.logErr(`Invalid JSON: ${raw}`)
            this.sendError('Could not parse message as JSON')
            return
        }

        const result = m.MessageSchema.safeParse(parsed)
        if (!result.success) {
            this.logErr(`Unknown message: ${raw}`)
            this.sendError('Unknown protocol message')
            return
        }

        const message = result.data

        // First message must be hello
        if (!this.handshakeComplete && message.type !== 'hello') {
            this.sendError('Expected hello as first message')
            this.socket.destroy()
            return
        }

        const handler = this.handlers[message.type]
        if (handler) {
            await handler(message)
        } else {
            this.logErr(`No handler for message type: ${message.type}`)
        }
    }

    sendMessage(msg: string) {
        this.socket.write(msg)
    }

    sendError(description: string) {
        this.sendMessage(m.makeErrorMessage(m.errorType.INVALID_FORMAT, description))
    }


    // ---- Protocol handlers ----

    private async handleHello(message: HelloMessage) {
        this.log(`Hello from ${message.agent} (${message.version})`)
        this.handshakeComplete = true
    }

    private async handleGetPeers() {
        this.log(`Requested peers`)
        const peers = getKnownPeers()
        this.sendMessage(m.makePeersMessage(peers))
    }

    private async handlePeers(message: PeersMessage) {
        this.log(`Received ${message.peers.length} peers`)
        appendPeers(message.peers, this.id)
    }

    private async handleError(message: ErrorMessage) {
        this.logErr(`${message.name} error: ${message.description}`)
    }

    private async handleGetObject(message: GetObjectMessage) {
        this.log(`Received getobject for ${message.objectid}`)
        if (await objectManager.exists(message.objectid)) {
            this.log(`We have object ${message.objectid}, sending it`)
            const object = await objectManager.get(message.objectid)
            this.sendMessage(m.makeObjectMessage(object))
        } else {
            this.log(`We don't have object ${message.objectid}`)
            this.sendMessage(m.makeErrorMessage(m.errorType.UNKNOWN_OBJECT, `Object ${message.objectid} not found`))
        }
    }

    private async handleIHaveObject(message: IHaveObjectMessage) {
        this.log(`Received ihaveobject for ${message.objectid}`)
        if (await objectManager.exists(message.objectid)) {
            this.log(`We already have object ${message.objectid}`)
            return
        }
        else {
        this.sendMessage(m.makeGetObjectMessage(message.objectid))
        this.log(`Requested object ${message.objectid}`)
        }
    }

    private async handleObject(message: ObjectMessage) {
        this.log(`Received object for ${message.objectid}`)
        if (await objectManager.exists(message.objectid)) {
            this.log(`We already have object ${message.objectid}`)
            return
        }
        try {
            const object = message.object
            const hash = hashObject(message.object)
            if (hash !== message.objectid) {
                this.logErr(`Object ID mismatch: has ${message.objectid}, expected hash: ${hash}`)
                return
            }
            const id = await objectManager.put(object)
            this.log(`Stored object ${id}`)
        } catch (error) {
            this.logErr(`Failed to store object ${message.objectid}: ${error}`)
        }
        // gossip!
        broadcast(m.makeIHaveObjectMessage(message.objectid))
    }

}