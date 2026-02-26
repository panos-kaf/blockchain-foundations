import { Socket } from 'net'
import { appendPeers } from './peers' 
import { ClientHelloMessage, GetPeersMessage, messageType, makePeersMessage } from './messages'
import { getKnownPeers } from './peers'


// const SERVER_PORT = 18018
// const SERVER_HOST = 'localhost'


export const startClient = (SERVER_HOST: string , SERVER_PORT: number, id: number, onClose?: () => void) => {
    
    const client = new Socket()

    let isDisconnected = false
    const handleClose = () => {
        if (isDisconnected) return
        isDisconnected = true
        log(`Connection to server ${SERVER_HOST}:${SERVER_PORT} closed`)
        if (onClose) onClose()
    }

    client.connect(SERVER_PORT, SERVER_HOST, () => {
        log(`Connected to server ${SERVER_HOST}:${SERVER_PORT}`)
    })
    
    const log = (msg: string, ...args: any[]) => console.log(`\x1b[36m[CLIENT ${id}]\x1b[0m ${msg}`, ...args)
    const logErr = (msg: string, ...args: any[]) => console.error(`\x1b[32m[CLIENT ${id}]\x1b[0m ${msg}`, ...args)

    client.write(ClientHelloMessage)
    client.write(GetPeersMessage)

    let buffer = ''
    client.on('data', (data) => {
        buffer += data
        const messages = buffer.split('\n')
        while (messages.length > 1) {
            let msg = messages.shift()
            // log(`Received message: ${msg}`)

            try {
                const parsedMessage = JSON.parse(msg || '')
                log(`Received '${parsedMessage.type}' message from server ${SERVER_HOST}:${SERVER_PORT}`)
                switch (parsedMessage.type) {
                    case messageType.HELLO:
                        // log(`Server ${SERVER_HOST}:${SERVER_PORT} says hello`)
                        break
                    case messageType.PEERS:
                        appendPeers(parsedMessage.peers, `${SERVER_HOST}:${SERVER_PORT}`)
                        log(`Updated known peers. Total known peers: ${getKnownPeers().length}`)
                        break
                    case messageType.GETPEERS:
                        client.write(makePeersMessage(getKnownPeers()))
                        break
                    default:
                        log(`${parsedMessage.type} messages not handled by client yet`)
                }
                
            } catch (e) {
                logErr(`Error parsing message: ${msg}`)
            }
        }
        if (messages[0] === undefined) {
            logErr(`Error in parsing messages`)
            return
        }
        buffer = messages[0]
    })


    client.on('error', (error) => {
        logErr(`Error: ${error}`)
        handleClose()
        // if (onClose) onClose()
    })

    client.on('close', () => {
        log(`Client ${id} disconnected`)
        handleClose()
        // if (onClose) onClose()
    })
}