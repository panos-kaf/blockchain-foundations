import { canonicalizeMessage, errorType, messageType } from './types'

export const ServerHello = canonicalizeMessage({
    type: messageType.HELLO,
    version: '0.10.0',
    agent: 'server-example',
})

export const ClientHello = canonicalizeMessage({
    type: messageType.HELLO,
    version: '0.10.0',
    agent: 'client-example',
})

export const GetPeers = canonicalizeMessage({
    type: messageType.GETPEERS,
})

// Error messages
export const InvalidFormatError = canonicalizeMessage({
    type: messageType.ERROR,
    name: errorType.INVALID_FORMAT,
    description: 'The message format is invalid',
})

export const InvalidHandshakeError = canonicalizeMessage({
    type: messageType.ERROR,
    name: errorType.INVALID_HANDSHAKE,
    description: 'Handshake not completed, expected hello message'
})
