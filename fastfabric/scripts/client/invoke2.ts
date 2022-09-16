/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

import { FileSystemWallet, Gateway} from 'fabric-network';
import {ProposalResponseObject,ProposalResponse} from 'fabric-network/node_modules/fabric-client'
import fs from 'fs';
import path from 'path';
import exec from 'child_process'

async function main() {
    var conflictPercentage = parseInt(process.argv[6])
    var endorserIdx = parseInt(process.argv[5])
    var repetitions = parseInt(process.argv[4]);
    var transferCount = parseInt(process.argv[3]);
    var iteration = parseInt(process.argv[2]);
    var bcResps: Array<ProposalResponseObject> = [];

    var ccpPath = path.resolve(__dirname, 'connection.json');
    const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
    const ccp = JSON.parse(ccpJSON);

    var endorserAddresses =exec.execSync("bash printEndorsers.sh").toString().replace("\n","").split(" ")

    ccp.name = ccp.name.replace("CHANNEL", process.env.CHANNEL)
    ccp.organizations.Org1.signedCert.path = ccp.organizations.Org1.signedCert.path.split("DOMAIN").join(process.env.PEER_DOMAIN)
    ccp.orderers.address.url = ccp.orderers.address.url.replace("ADDRESS", process.env.ORDERER_ADDRESS)
    ccp.peers.address.url = ccp.peers.address.url.replace("ADDRESS", endorserAddresses[endorserIdx])
    const user = "Admin@" + process.env.PEER_DOMAIN
    console.log(`Thread no.: ${iteration},\tTx per process: ${transferCount / 2},\trepetitions: ${repetitions}`);
    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(__dirname, './wallet');
    const wallet = new FileSystemWallet(walletPath);
    var tIdx = 0
    try {

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(user);
        if (!userExists) {
            console.log(`An identity for the user "${user}" does not exist in the wallet`);
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        var gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: false } });

        // Get the network (channel) our contract is deployed to.
        var network = await gateway.getNetwork(String(process.env.CHANNEL));
        var channel  = network.getChannel()
        var client = gateway.getClient()


        // Submit the specified transaction.

        for (var r = 0; r < repetitions; r++) {
            for (var i = iteration * transferCount; i < iteration * transferCount + transferCount - 1; i += 2) {
                var a
                if (Math.random() <= conflictPercentage / 100) {
                    a = ["account0", "account1", "1"]
                } else {
                    a = ["account".concat(i.toString()), "account".concat((i + 1).toString()), "1"]
                }
                
                var resp = await channel.sendTransactionProposal({
                    chaincodeId: String(process.env.CHAINCODE),
                    fcn: "transfer",
                    args: a,
                    txId: client.newTransactionID()
                });
                bcResps.push(resp);
            }
        }
        console.log(`Thread ${iteration}: Endorsement done, got ${bcResps.length} endorsments for ${transferCount / 2 * repetitions} transactions`)

        console.log(`Thread ${iteration}: Start sending transactions`);

        for (let i = 0; i < bcResps.length; i++) {
            var propResp = <ProposalResponse[]>(bcResps[i][0]);
            var res = await channel.sendTransaction({ proposalResponses: propResp, proposal: bcResps[i][1] });
            if (res.status != "SUCCESS") {
                console.log(`Thread ${iteration}: ${res.status}`);
            }
            tIdx++
        }

        // Disconnect from the gateway.
        await gateway.disconnect();
        console.log(`Thread ${iteration} is done!`);
    } catch (error) {
        console.error(`Thread ${iteration}: Failed to submit transaction ${tIdx}: ${error}`);
        process.exit(1);
    }
}

main();
