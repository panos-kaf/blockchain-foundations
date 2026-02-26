import { startClient } from "./client";
import { startServer } from "./server";
import { BOOTSTRAP_PEERS, selectRandomPeersPerSource} from "./peers";

startServer(18018);

let clientID = 1
let activeConnections = 0
let connectedPeers = new Set<string>()
let blacklistedPeers = new Map<string, number>()

const connect = (host: string, port: number) => {
    const peerId = `${host}:${port}`
    
    if (connectedPeers.has(peerId)) return // Already connected
    
    if (blacklistedPeers.has(peerId)) {
        const blacklistedUntil = blacklistedPeers.get(peerId)!
        if (Date.now() < blacklistedUntil) {
            console.log(`\x1b[34m[Node]\x1b[0m Skipping blacklisted peer ${peerId}`)
            return
        }
        blacklistedPeers.delete(peerId)
    }
    
    connectedPeers.add(peerId)
    activeConnections++
    
    const myId = clientID++ 
    
    const onClose = () => {
        activeConnections--
        connectedPeers.delete(peerId)
        console.log(`\x1b[34m[Node]\x1b[0m Connection to ${peerId} closed. Active: ${activeConnections}`)
        const cooldown = 15 * 60 * 1000 // 15mins
        blacklistedPeers.set(peerId, Date.now() + cooldown)
    }

    startClient(host, port, myId, onClose)
}

for (const peer of BOOTSTRAP_PEERS) {
    const host = peer.split(':')[0]
    const port = peer.split(':')[1]
    if (host === undefined || port === undefined) {
        console.error(`Invalid bootstrap peer ${peer}`)
        continue
    }
    connect(host, parseInt(port))
}

// Periodic Discovery
setInterval(() => {

    const nodes = selectRandomPeersPerSource(3)
    
    if (activeConnections < 20) {

        for (const peer of nodes) {

            if (BOOTSTRAP_PEERS.includes(peer)) continue;
            
            const [host, portStr] = peer.split(':')

            if (host && portStr) {
                console.log(`\x1b[34m[Node]\x1b[0m Discovery: attempting to connect to ${peer}`)
                connect(host, parseInt(portStr))
            }
        }
    }
    else console.log(`\x1b[34m[Node]\x1b[0m Capped Active connections: ${activeConnections}. Skipping discovery.`)
}, 5000)