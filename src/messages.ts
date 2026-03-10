import { z } from 'zod'
import canonicalize from 'canonicalize'
import { hashObject } from './object'

export enum errorType {
    INTERNAL_ERROR = 'INTERNAL_ERROR',
    INVALID_FORMAT = 'INVALID_FORMAT',
    UNKNOWN_OBJECT = 'UNKNOWN_OBJECT',
    UNFINDABLE_OBJECT = 'UNFINDABLE_OBJECT',
    INVALID_HANDSHAKE = 'INVALID_HANDSHAKE',
    INVALID_TX_OUTPOINT = 'INVALID_TX_OUTPOINT',
    INVALID_TX_SIGNATURE = 'INVALID_TX_SIGNATURE',
    INVALID_TX_CONSERVATION = 'INVALID_TX_CONSERVATION',
    INVALID_BLOCK_COINBASE = 'INVALID_BLOCK_COINBASE',
    INVALID_BLOCK_TIMESTAMP = 'INVALID_BLOCK_TIMESTAMP',
    INVALID_BLOCK_POW = 'INVALID_BLOCK_POW',
    INVALID_GENESIS = 'INVALID_GENESIS'
}

export enum messageType {
    HELLO = 'hello',
    ERROR = 'error',
    GETPEERS = 'getpeers',
    PEERS = 'peers',
    GETOBJECT = 'getobject',
    IHAVEOBJECT = 'ihaveobject',
    OBJECT = 'object',
    GETMEMPOOL = 'getmempool',
    MEMPOOL = 'mempool',
    GETCHAINTIP = 'getchaintip',
    CHAINTIP = 'chaintip'
}

export let canonicalizeMessage = (msg: Message): string => {
    const res = canonicalize(msg)
    if (res === undefined) {
        throw new Error(`Error in canonicalizing message ${msg}`)
    }
    return res + '\n'
}

// Message schemas and types

export const HelloSchema = z.object({
    type: z.literal(messageType.HELLO),
    version: z.string().regex(/^0\.10\.[0-9]+$/),
    agent: z.string().max(1000).optional()
})

export const ErrorSchema = z.object({
    type: z.literal(messageType.ERROR),
    name: z.enum(errorType),
    description: z.string().max(1000)
})

export const GetPeersSchema = z.object({
    type: z.literal(messageType.GETPEERS),
})

export const PeersSchema = z.object({
    type: z.literal(messageType.PEERS),
    peers: z.array(z.string().regex(/^((?:\d{1,3}\.){3}\d{1,3}|\[[a-fA-F0-9:]+\]|[a-zA-Z0-9.-]+):[0-9]{1,5}$/).max(1000))
})

export const GetObjectSchema = z.object({
    type: z.literal(messageType.GETOBJECT),
    objectid: z.string().length(64)
})

export const IHaveObjectSchema = z.object({
    type: z.literal(messageType.IHAVEOBJECT),
    objectid: z.string().length(64)
})

export const TransactionSchema = z.object({
    type: z.literal('transaction'),
    inputs: z.array(z.object({
                        outpoint: z.object({
                            txid: z.string().length(64),
                            index: z.number().int().nonnegative()
                        }),
                        sig: z.string().length(128)
                        })
                    ).max(1000),
    outputs: z.array(z.object({
                        pubkey: z.string().length(128),
                        value: z.number().int().nonnegative()
                        })
                    ).max(1000),
})

export const BlockSchema = z.object({
    type: z.literal('block'),
    T: z.string(),
    created: z.number(),
    miner: z.string().max(128).optional(),
    nonce: z.string().length(64),
    note: z.string().max(1000).optional(),
    previd: z.string().length(64).nullable(),
    studentids: z.array(z.string().max(128)).max(10).optional(),
    txids: z.array(z.string().length(64)).max(1000),
})

export const ObjectSchema = z.object({
    type: z.literal(messageType.OBJECT),
    objectid: z.string().length(64),
    object: z.union([BlockSchema, TransactionSchema])
    })
    

export const GetMempoolSchema = z.object({
    type: z.literal(messageType.GETMEMPOOL),
})

export const MempoolSchema = z.object({
    type: z.literal(messageType.MEMPOOL),
    txids: z.array(z.string().length(64)).max(1000)
})

export const GetChaintipSchema = z.object({
    type: z.literal(messageType.GETCHAINTIP),
})

export const ChaintipSchema = z.object({
    type: z.literal(messageType.CHAINTIP),
    blockid: z.string().length(64),
})

export const MessageSchema = z.discriminatedUnion(
    'type', [
            HelloSchema, ErrorSchema, 
            GetPeersSchema, PeersSchema, 
            GetObjectSchema, IHaveObjectSchema, 
            ObjectSchema, GetMempoolSchema, 
            MempoolSchema, GetChaintipSchema, 
            ChaintipSchema
        ]
)

export type Message = z.infer<typeof MessageSchema>;
export type HelloMessage = z.infer<typeof HelloSchema>;
export type ErrorMessage = z.infer<typeof ErrorSchema>;
export type GetPeersMessage = z.infer<typeof GetPeersSchema>;
export type PeersMessage = z.infer<typeof PeersSchema>;
export type GetObjectMessage = z.infer<typeof GetObjectSchema>;
export type IHaveObjectMessage = z.infer<typeof IHaveObjectSchema>;
export type ObjectMessage = z.infer<typeof ObjectSchema>;
export type GetMempoolMessage = z.infer<typeof GetMempoolSchema>;
export type MempoolMessage = z.infer<typeof MempoolSchema>;
export type GetChaintipMessage = z.infer<typeof GetChaintipSchema>;
export type ChaintipMessage = z.infer<typeof ChaintipSchema>;

export type TransactionType = z.infer<typeof TransactionSchema>;
export type BlockType = z.infer<typeof BlockSchema>;


// Message constructors

export const makeHelloMessage = (version: string = '0.10.0', agent: string) => {
    const helloMessage: Message = {
        type: messageType.HELLO,
        version: version,
        agent: agent
    }
    try {
        HelloSchema.parse(helloMessage)
    } catch (error) {
        throw new Error(`Invalid hello message with version ${version} and agent ${agent}`)
    }
    return canonicalizeMessage(helloMessage)
}

export const makeErrorMessage = (name: errorType, description: string) => {
    const errorMessage: Message = {
        type: messageType.ERROR,
        name: name,
        description: description
    }
    try {
        ErrorSchema.parse(errorMessage)
    } catch (error) {
        throw new Error(`Invalid error message with name ${name} and description ${description}`)
    }
    return canonicalizeMessage(errorMessage)
}

export const makePeersMessage = (peers: string[]) => {
    const peersMessage: Message = {
        type: messageType.PEERS,
        peers
    }
    try {
        PeersSchema.parse(peersMessage)
    } catch (error) {
        throw new Error(`Invalid peers message with peers ${peers}`)
    }
    return canonicalizeMessage(peersMessage)
}

export const makeGetObjectMessage = (objectid: string) => {
    const getObjectMessage: Message = {
        type: messageType.GETOBJECT,
        objectid
    }
    try {
        GetObjectSchema.parse(getObjectMessage)
    } catch (error) {
        throw new Error(`Invalid getobject message with objectid ${objectid}`)
    }
    return canonicalizeMessage(getObjectMessage)
}

export const makeIHaveObjectMessage = (objectid: string) => {
    const ihaveObjectMessage: Message = {
        type: messageType.IHAVEOBJECT,
        objectid
    }
    try {
        IHaveObjectSchema.parse(ihaveObjectMessage)
    } catch (error) {
        throw new Error(`Invalid ihaveobject message with objectid ${objectid}`)
    }
    return canonicalizeMessage(ihaveObjectMessage)
}

export const makeObjectMessage = (object: TransactionType | BlockType) => {
    const objectid = hashObject(object)
    const objectMessage: Message = {
        type: messageType.OBJECT,
        objectid,
        object
    }
    try {
        ObjectSchema.parse(objectMessage)
    } catch (error) {
        throw new Error(`Invalid object message with id ${objectid}`)
    }
    return canonicalizeMessage(objectMessage)
}


// Static messages

export const StaticHello = canonicalizeMessage({
    type: messageType.HELLO,
    version: '0.10.0',
    agent: 'marabobos',
})

export const StaticGetPeers = canonicalizeMessage({
    type: messageType.GETPEERS,
})

// Error messages
export const StaticInvalidFormatError = canonicalizeMessage({
    type: messageType.ERROR,
    name: errorType.INVALID_FORMAT,
    description: 'The message format is invalid',
})

export const StaticInvalidHandshakeError = canonicalizeMessage({
    type: messageType.ERROR,
    name: errorType.INVALID_HANDSHAKE,
    description: 'Handshake not completed, expected hello message'
})
