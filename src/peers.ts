export const BOOTSTRAP_PEERS = 
[
    '95.179.158.137:18018', 
    '95.179.132.22:18018', 
    '45.32.235.245:18018',
]

export const appendPeers = (peers: string[]) => {
    peers.forEach(peer => {
        if (!KNOWN_PEERS.has(peer)) {
            KNOWN_PEERS.add(peer)
            console.log(`Added new peer: ${peer}`)
        }
    })
}

export const KNOWN_PEERS = new Set<string>(BOOTSTRAP_PEERS)
