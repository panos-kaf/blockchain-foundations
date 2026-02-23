import { createServer } from 'net'
import { makePeersMessage, MessageSchema } from './types'
import * as m from './messages'
import { BOOTSTRAP_PEERS } from './peers'

const PORT = 18018

const server = createServer(async (socket) => {

    let handshaked = false

    const id = `${socket.remoteAddress}:${socket.remotePort}`
    console.log(`Client connected from ${id}`)
    
    socket.write(m.ServerHello)
    
    socket.on('error', (error) => {
        console.error(`[${id}]: Received error ${error}`)
    })

    socket.on('close', () => {
        console.log(`[${id}]: Client disconnected`)
    })

    let buffer = ''
    socket.on('data', (data) => {
        buffer += data

        const messages = buffer.split('\n')
        while (messages.length > 1) {
            let msg = messages.shift()
            if (msg === undefined) {
                console.error(`Error defragmenting messages`)
                return
            }

            let message
            try {
                message = JSON.parse(msg)
            } catch (error) {
                console.error(`Error parsing message as JSON`, message)
                // socket.write(`Received invalid message that could not parse as json` + msg)
                socket.write(m.InvalidFormatError)
                socket.end()
                return
            }

            try {
                message = MessageSchema.parse(message)
            } catch (_) {
                console.error(`Unknown protocol message`, message)
                // socket.write(`Received invalid protocol message` + message)
                socket.write(m.InvalidFormatError)
                // socket.end()
                // return
            }

            console.log(`[${id}]: Received message`, message)

            if (message.type === 'hello') {
                handshaked = true
            }
            else if (!handshaked) {
                console.error(`Handshake not completed, expected hello message but received ${message.type}`)
                socket.write(m.InvalidHandshakeError)
                socket.end()
                return
            }

            if (message.type === 'getpeers') {
                socket.write(makePeersMessage(BOOTSTRAP_PEERS))
            }
            
            if (messages[0] === undefined) {
                console.error(`Error in parsing messages`)
                return
            }

            buffer = messages.join('\n')
        }

    })

})

server.listen(PORT, () => {
    console.log(`Server listening on port ${PORT}`)
})