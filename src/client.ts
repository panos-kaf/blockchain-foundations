import { Socket } from 'net'
import Peer from './node'


export const startClient = (SERVER_HOST: string, SERVER_PORT: number, id: number, onClose?: () => void) => {

    const log = (msg: string, ...args: any[]) => console.log(`\x1b[36m[CLIENT][${id}]\x1b[0m ${msg}`, ...args)
    const logErr = (msg: string, ...args: any[]) => console.error(`\x1b[32m[CLIENT][${id}]\x1b[0m ${msg}`, ...args)

    const client = new Socket()
    
    client.connect(SERVER_PORT, SERVER_HOST, () => {
        const peer = new Peer(client, 'client', log, logErr)
        log(`Connected to server ${SERVER_HOST}:${SERVER_PORT}`)
        peer.greet()

        peer.onDisconnect = () => {
            log(`Disconnected from ${SERVER_HOST}:${SERVER_PORT}`)
            if (onClose) onClose()
        }
    })

    client.on('error', (error) => {
        logErr(`Error connecting to ${SERVER_HOST}:${SERVER_PORT}: ${error}`)
        if (onClose) onClose()
    })
}