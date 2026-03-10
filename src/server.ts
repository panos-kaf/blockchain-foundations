import { createServer } from 'net'
import Peer from './node'

let incomingClientID = 1;

const log = (msg: string, ...args: any[]) => console.log(`\x1b[35m[SERVER][${incomingClientID}]\x1b[0m ${msg}`, ...args)
const logErr = (msg: string, ...args: any[]) => console.error(`\x1b[31m[SERVER][${incomingClientID}]\x1b[0m ${msg}`, ...args)

export const startServer = (PORT: number = 18018) => {
    const server = createServer((socket) => {
        const clientID = incomingClientID++;
        const peer = new Peer(socket, 'server', log, logErr)
        log(`Client [${clientID}] connected from ${peer.id}`)
        peer.greet()

        peer.onDisconnect = () => {
            log(`Client [${clientID}] ${peer.id} disconnected`)
        }
    })

    server.listen(PORT, () => {
        log(`Server listening on port ${PORT}`)
    })
}