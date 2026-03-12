import { Level } from 'level'
import canonicalize from 'canonicalize'
import { blake2s } from '@noble/hashes/blake2'
import { bytesToHex } from '@noble/hashes/utils'
import { verify } from '@noble/ed25519'
import type { TransactionType, CoinbaseTransactionType, BlockType } from './messages'

type ObjectType = TransactionType | CoinbaseTransactionType | BlockType

const FIND_TIMEOUT_MS = 5000

export const validateObject = async (object: ObjectType, objectid: string) =>  {
    const hash = hashObject(object)
    if (hash !== objectid) {
        throw new Error(`Object ID mismatch: has ${objectid}, expected hash: ${hash}`)
    }
    switch (object.type) {
        case 'block': {
            const block = object as BlockType
            return
        }
        case 'transaction': {
            const tx = object as CoinbaseTransactionType
            if (tx.height !== undefined) {
                // handle coinbase txs
                return
            }
            else { // handle normal txs
                let totalInputValue = 0, totalOutputValue = 0
                const tx = object as TransactionType
                const inputs = tx.inputs
                if (inputs.length === 0) {
                    throw new Error('Transaction must have at least one input')
                }
                for (const input of inputs) {
                    const outpoint = input.outpoint
                    const incomingTx = await objectManager.get(outpoint.txid)
                    if (incomingTx === undefined || incomingTx.type !== 'transaction') {
                        throw new Error(`Input references unknown transaction ${outpoint.txid}`)
                    }
                    const incomingTxTyped = incomingTx as TransactionType
                    if (outpoint.index < 0 || outpoint.index >= incomingTxTyped.outputs.length) {
                        throw new Error(`Input references invalid output index ${outpoint.index} in transaction ${outpoint.txid}`)
                    }

                    const output = incomingTxTyped.outputs[outpoint.index]
                    if (output === undefined) {
                        throw new Error(`Input references non-existent output index ${outpoint.index} in transaction ${outpoint.txid}`)
                    }

                    totalInputValue += output.value

                    const sig = input.sig
                    if (sig.length !== 128) {
                        throw new Error(`Invalid signature length for input referencing transaction ${outpoint.txid}`)
                    }
                    const pubkey = output.pubkey
                    if (!verify(Buffer.from(sig, 'hex'), Buffer.from(objectid, 'hex'), Buffer.from(pubkey, 'hex'))) {
                        throw new Error(`Invalid signature for input referencing transaction ${outpoint.txid}`)
                    }
                    
                }
                const outputs = tx.outputs
                if (outputs.length > 0) {
                    for (const output of outputs) {
                        if (output.value < 0) {
                            throw new Error('Output value cannot be negative')
                        }
                        totalOutputValue += output.value
                        if (output.pubkey.length !== 64) {
                            throw new Error('Invalid pubkey length in output')
                        }
                    }
                }
                if (totalOutputValue > totalInputValue) {
                    throw new Error('Output value exceeds input value')
                } else {
                    const crumbs = totalInputValue - totalOutputValue
                    console.log(`Transaction ${objectid} is valid with ${crumbs} crumbs of fee`)
                    return crumbs
                }
            }
        }
        default:
            const type = (object as any).type
            throw new Error(`Unknown object type: ${type}`)
    }
}

export const hashObject = (object: ObjectType): string => {
    const canonicalized = canonicalize(object)
    if (canonicalized === undefined) {
        throw new Error('Failed to canonicalize object')
    }
    return bytesToHex(blake2s(canonicalized))
}

interface PendingWaiter {
    resolve: (object: ObjectType) => void
    reject: (error: Error) => void
}

class ObjectManager {
    db = new Level('./db', { valueEncoding: 'json'})
    pendingFinds: Map<string, PendingWaiter[]> = new Map()

    id(object: ObjectType): string {
        const canonicalized = canonicalize(object)
        if (canonicalized === undefined) {
            throw new Error('Failed to canonicalize object')
        }
        return bytesToHex(blake2s(canonicalized))
    }

    async exists(id: string): Promise<boolean> {
        return await this.db.has(id)
    }

    async get(id: string): Promise<ObjectType> {
        const object = await this.db.get(id)
        if (object === undefined) {
            throw new Error(`Object ${id} not found`)
        }
        return object as unknown as ObjectType
    }

    async put(object: ObjectType): Promise<string> {
        const objectId = this.id(object)
        await this.db.put(objectId, object as any)

        const waiters = this.pendingFinds.get(objectId)
        if (waiters) {
            for (const waiter of waiters) {
                waiter.resolve(object)
            }
            this.pendingFinds.delete(objectId)
        }

        return objectId
    }

    async findObject(objectId: string, sendGetObject: (id : string) => void): Promise<ObjectType> {
        try {
            return await this.get(objectId)
        } catch {}

        sendGetObject(objectId)

        const waitPromise = new Promise<ObjectType>((resolve) => {
            const existing = this.pendingFinds.get(objectId)
            if (existing) {
                existing.push({ resolve, reject: () => {} })
            } else {
                this.pendingFinds.set(objectId, [{ resolve, reject: () => {} }])
            }
            
        })

        const timeoutPromise = new Promise<never>((_, reject) => {
            setTimeout(() => {
                this.pendingFinds.delete(objectId)
                reject(new Error(`Timeout waiting for object ${objectId}`))
            }, FIND_TIMEOUT_MS)
        })

        return Promise.race([waitPromise, timeoutPromise])
    }
}

export const objectManager = new ObjectManager()
