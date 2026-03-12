import * as net from 'net';
import { makeIHaveObjectMessage, makeObjectMessage, StaticHello, type BlockType } from './messages';

const PORT = 18018;
const HOST = 'localhost';

const mockBlockId = "00000000522473196b73bc619a8b18472c4cb4c6caf785a13fa32aaae7222ff6";
const mockBlock : BlockType = {
    type: "block",
    T: "00000000abc00000000000000000000000000000000000000000000000000000",
    created: 1771159355,
    miner: "TestMiner",
    nonce: "00dd82159556175752d9ba7349df67bddd237b59183747383f7b720e85c32347",
    note: "Local test block",
    previd: null,
    txids: ["8265faf623dfbcb17528fcd2e67fdf78de791ed4c7c60480e8cd21c6cdc8bcd4"]
};

const client = net.createConnection({ port: PORT, host: HOST }, () => {
    console.log('Connected to local node');
    
    // 1. Send Hello
    client.write(StaticHello);

    // 2. Advertise the object
    client.write(makeIHaveObjectMessage(mockBlockId));
});

client.on('data', (data) => {
    const msgs = data.toString().trim().split('\n');
    for (const msg of msgs) {
        if (!msg) continue;
        const parsed = JSON.parse(msg);
        console.log('Received:', parsed);

        // 3. If node asks for it, send the object
        if (parsed.type === 'getobject' && parsed.objectid === mockBlockId) {
            console.log('Node requested object, sending it now...');
            client.write(makeObjectMessage(mockBlock));
        }
    }
});

client.on('error', (err) => console.error('Error:', err.message));
client.on('close', () => console.log('Connection closed'));

// object message for testing:
//{"type": "object", "objectid": "038eab4238d0ad86e9b8ebcedecd2abf4b71f8e2cd94a65c36d96f11a12e4ba7", "object": {"T": "00000000abc00000000000000000000000000000000000000000000000000000","created": 1671148800,"miner": "Marabu Bounty Hunter","nonce": "15551b5116783ace79cf19d95cca707a94f48e4cc69f3db32f41081dab3e6641","note": "First block on genesis, 50 bu reward","previd": "00000000522473196b73bc619a8b18472c4cb4c6caf785a13fa32aaae7222ff6","txids": ["8265faf623dfbcb17528fcd2e67fdf78de791ed4c7c60480e8cd21c6cdc8bcd4"],"type": "block"}}