import { createServer } from 'net'
import { makeErrorMessage, makePeersMessage, MessageSchema } from './messages'
import { messageType, errorType } from './messages'
import { ServerHelloMessage, InvalidHandshakeError, GetPeersMessage} from './messages'
import { getKnownPeers } from './peers'

// const PORT = 18018

const log = (msg: string, ...args: any[]) => console.log(`\x1b[35m[SERVER]\x1b[0m ${msg}`, ...args)
const logErr = (msg: string, ...args: any[]) => console.error(`\x1b[31m[SERVER]\x1b[0m ${msg}`, ...args)

export const startServer = (PORT: number = 18018) => {
    const server = createServer(async (socket) => {

        let handshaked = false

        const id = `${socket.remoteAddress}:${socket.remotePort}`
        log(`Client connected from ${id}`)
        
        socket.write(ServerHelloMessage)
        socket.write(GetPeersMessage)
        
        socket.on('error', (error) => {
            logErr(`Received error ${error}`)
        })

        socket.on('close', () => {
            log(`Client ${id} disconnected`)
        })

        let buffer = ''
        socket.on('data', (data) => {
            buffer += data

            const messages = buffer.split('\n')
            while (messages.length > 1) {
                let msg = messages.shift()
                if (msg === undefined) {
                    logErr(`Error defragmenting messages`)
                    continue
                }

                let message
                try {
                    message = JSON.parse(msg)
                } catch (error) {
                    logErr(`Error parsing JSON from client ${id}: ${error}`)
                    const err = makeErrorMessage(errorType.INVALID_FORMAT, 'Invalid JSON')
                    socket.write(err)
                    continue
                }

                log(`Received '${message.type}' message from ${id}`)
                
                try {
                    message = MessageSchema.parse(message)
                } catch (_) {
                    logErr(`Unknown protocol message`, message)
                    const err = makeErrorMessage(errorType.INVALID_FORMAT, 'Unknown protocol message')
                    socket.write(err)
					continue
                }

                if (message.type === messageType.HELLO) {
                    handshaked = true
                }
                else if (!handshaked) {
                    logErr(`Handshake not completed, expected hello message but received ${message.type}`)
                    socket.write(InvalidHandshakeError)
                    socket.end()
                    return
                }

                switch (message.type) {
                    case messageType.HELLO:
                        log(`Client ${id}`, message.agent ? `(${message.agent})` : '', `says hello`)
                        break
                    case messageType.GETPEERS:
                        socket.write(makePeersMessage(getKnownPeers()))
                        log(`Sent peers to client ${id}`)
                        break
                    default:
                        log(`${message.type} messages not handled by server yet`)
                }

                if (messages[0] === undefined) {
                    logErr(`Error in parsing messages`)
                    return
                }
                buffer = messages.join('\n')
            }
        })
    })

    server.listen(PORT, () => {
        log(`Server listening on port ${PORT}`)
    })
}
