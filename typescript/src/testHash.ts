import { hashObject } from './object'
import { type BlockType } from './messages'

const genID = '00000000522473196b73bc619a8b18472c4cb4c6caf785a13fa32aaae7222ff6'

const genesisBlock: BlockType = {
    type: 'block',
    T: '00000000abc00000000000000000000000000000000000000000000000000000',
    created: 1771159355,
    miner: "Marabu",
    nonce: "00dd82159556175752d9ba7349df67bddd237b59183747383f7b720e85c32347",
    note: "Financial Times 2026-02-13: Crypto's battle with the banks is splitting Trump's base",
    previd: null,
    txids: [],
}

const generatedGenID = hashObject(genesisBlock)

console.log('Expected Genesis Block ID:\t', genID)
console.log('Generated Genesis Block ID:\t', generatedGenID)

if (genID === generatedGenID) {
    console.log('Test passed: Generated ID matches expected ID')
} else {
    console.error('Test failed: Generated ID does not match expected ID')
}