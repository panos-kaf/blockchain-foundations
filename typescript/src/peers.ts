import * as fs from 'fs'
import { isIP } from 'net'
import * as path from 'path'

export const BOOTSTRAP_PEERS = 
[
    '95.179.158.137:18018', 
    '95.179.132.22:18018', 
    '45.32.235.245:18018',
]

const PEERS_FILE = path.join(process.cwd(), 'src/peers.csv')

const peersMap = new Map<string, string>()

const loadPeers = (): Map<string, string> => {
    
    BOOTSTRAP_PEERS.forEach(peer => peersMap.set(peer, 'bootstrap'))

    try {
        if (fs.existsSync(PEERS_FILE)) {
            const data = fs.readFileSync(PEERS_FILE, 'utf8')
            const lines = data.split('\n')
            // Skip header if present or handle empty lines
            lines.forEach(line => {
                const parts = line.trim().split(',')
                if (parts.length >= 2) {
                    const [peer, source] = parts
                    if (peer && peer !== 'Address') { // efficient header check
                         peersMap.set(peer, source || 'unknown')
                    }
                }
            })
        }
    } catch (e) {
        console.error('Failed to load peers file:', e)
    }
    return peersMap
}

const savePeers = () => {
    try {
        const lines = ['Address,Source']
        KNOWN_PEERS.forEach((source, peer) => {
            lines.push(`${peer},${source}`)
        })
        fs.writeFileSync(PEERS_FILE, lines.join('\n'))
    } catch (e) {
        console.error('Failed to save peers file:', e)
    }
}

export const KNOWN_PEERS = loadPeers()

export const getKnownPeers = (): string[] => Array.from(KNOWN_PEERS.keys())

// write bootstrap peers
if (!fs.existsSync(PEERS_FILE)) {
    savePeers()
}

// Regex for IPv4:Port
const IPV4_REGEX = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):([0-9]{1,5})$/

const IPV6_REGEX = /^\[([a-fA-F0-9:]+)\]:([0-9]{1,5})$/

// Regex for Domain:Port (e.g., node.example.com:18018 or localhost:18018)
// Matches alphanumeric parts separated by dots, ensuring at least one dot usually (unless localhost), ending with :port
const DOMAIN_REGEX = /^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}:([0-9]{1,5})$/

const sanitizePeer = (peer: string): string | null => {

    const trimmed = peer.trim()
    const isIPv4 = IPV4_REGEX.test(trimmed)
    const isIPv6 = IPV6_REGEX.test(trimmed)
    const isDomain = DOMAIN_REGEX.test(trimmed)

    if (!isIPv4 && !isIPv6 && !isDomain) {
        return null
    }

    const lastColonIndex = trimmed.lastIndexOf(':')
    const portStr = trimmed.substring(lastColonIndex + 1)
    const port = parseInt(portStr, 10)

    if (isNaN(port) || port <= 0 || port > 65535) {
        return null
    }

    const host = trimmed.substring(0, lastColonIndex)

    if (isIPv6) {
        if (host === '::1' || host === '[::1]' ||
            host.startsWith('[fe80:') || host.startsWith('[fc00:'))
        {
            return null
        }
    }

    if (host === 'localhost') return null

    // if (!trimmed.includes(':')) return null

    if (isIPv4) {

        if (host.startsWith('127.') || host.startsWith('0.') ||
            host.startsWith('192.168.') || host.startsWith('10.')){
            return null
        }

        const octets = host.split('.')
        if (octets.length !== 4) return null
        for (const octet of octets) {
            const num = parseInt(octet, 10)
            if (isNaN(num) || num < 0 || num > 255) {
                return null
            }
        }

        if (host.startsWith('172.')) {
            if (octets[1] === undefined) return null
            const secondOctet = parseInt(octets[1], 10)
            if (isNaN(secondOctet) || secondOctet < 16 || secondOctet > 31) {
                return null
            }
        }
    }
    
    return trimmed
}
export const appendPeers = (peers: string[], server: string) => {
    let changed = false
    
    peers.forEach(peer => {
        const sanitizedPeer = sanitizePeer(peer)

        // If sanitizedPeer is null, it was invalid, so we skip it
        if (sanitizedPeer && !KNOWN_PEERS.has(sanitizedPeer)) {
            KNOWN_PEERS.set(sanitizedPeer, server)
            console.log(`Added new peer: ${sanitizedPeer} from server ${server}`)
            changed = true
        }
    })
    
    if (changed) {
        console.log(`Saving ${KNOWN_PEERS.size} peers to disk...`)
        savePeers()
    } else {
        // console.log('No new peers to save.')
    }
}

export const selectRandomPeersPerSource = (count: number = 1): string[] => {
    const peersBySource = new Map<string, string[]>()
    const selectedPeers: string[] = []

    peersMap.forEach((source, peer) => {
        if (!peersBySource.has(source)) {
            peersBySource.set(source, [])
        }
        peersBySource.get(source)!.push(peer)
    })

    peersBySource.forEach((peers) => {
        if (peers.length <= count) {
            // If fewer peers than asked, take them all
            selectedPeers.push(...peers)
        } else {
            // Pick count unique random indexes
            const pickedIndexes = new Set<number>()
            while (pickedIndexes.size < count) {
                const index = Math.floor(Math.random() * peers.length)
                pickedIndexes.add(index)
            }
            
            pickedIndexes.forEach(index => {
                if (peers[index] !== undefined)
                    selectedPeers.push(peers[index])
            })
        }
    })

    return selectedPeers
}
