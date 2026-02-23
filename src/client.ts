import { Socket } from 'net'
import * as messages from './messages'
import  { appendPeers, KNOWN_PEERS, BOOTSTRAP_PEERS } from './peers' 
import * as types from './types'

const SERVER_PORT = 18018
// const SERVER_HOST = '95.179.158.137'
const SERVER_HOST = 'localhost'

const client = new Socket()
client.connect(SERVER_PORT, SERVER_HOST, () => {
    console.log(`Connected to server ${SERVER_HOST}:${SERVER_PORT}`)
})

const getPeers = () => {
    const peersMessage = messages.GetPeers
    client.write(peersMessage)

    // wait for response
    let buffer = ''
    client.on('data', (data) => {
        buffer += data
        const messages = buffer.split('\n')
        while (messages.length > 1) {
            let message = messages.shift()
            console.log(`Received message: ${message}`)
        }
        if (messages[0] === undefined) {
            console.error(`Error in parsing messages`)
            return
        }
        buffer = messages[0]
    }) 

}

client.write(messages.ClientHello)
client.write(messages.GetPeers)

let buffer = ''
client.on('data', (data) => {
    buffer += data
    const messages = buffer.split('\n')
    while (messages.length > 1) {
        let msg = messages.shift()
        console.log(`Received message: ${msg}`)

        try {
            const parsedMessage = JSON.parse(msg || '')
            if (parsedMessage.type === 'peers') {
                console.log(`Received peers: ${parsedMessage.peers}`)
                appendPeers(parsedMessage.peers)
            }
        } catch (e) {
            console.error(`Error parsing message: ${msg}`)
        }
    }
    if (messages[0] === undefined) {
        console.error(`Error in parsing messages`)
        return
    }
    buffer = messages[0]
})


client.on('error', (error) => {
    console.error(`Received error ${error}`)
})

client.on('close', () => {
    console.log(`Client disconnected`)
})
