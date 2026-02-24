import { startClient } from "./client";
import { startServer } from "./server";
import { BOOTSTRAP_PEERS, getKnownPeers} from "./peers";

startServer(18018);

for (const peer of BOOTSTRAP_PEERS) {
    const host = peer.split(':')[0]
    const port = peer.split(':')[1]
    if (host === undefined || port === undefined) {
        console.error(`Invalid bootstrap peer ${peer}`)
        continue
    }
    startClient(host, parseInt(port))
}