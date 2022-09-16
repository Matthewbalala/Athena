/*
 *  SPDX-License-Identifier: Apache-2.0
 */

'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
const path = require('path');

const cfgPath = path.resolve(process.env.FABRIC_CFG_PATH);

// A wallet stores a collection of identities
const wallet = new FileSystemWallet('wallet');

async function main() {

    // Main try/catch block
    try {

        // Identity to credentials to be stored in the wallet
        const credPath = path.join(cfgPath, '/crypto-config/peerOrganizations/'+process.env.PEER_DOMAIN+'/users/Admin@'+process.env.PEER_DOMAIN);
        const cert = fs.readFileSync(path.join(credPath, '/msp/signcerts/Admin@'+process.env.PEER_DOMAIN+'-cert.pem')).toString();
        const keystore = path.join(credPath, '/msp/keystore/')

        var files = fs.readdirSync(keystore)
        var key = fs.readFileSync(path.join(keystore,files[0])).toString()
        
        // Load credentials into wallet
        const identityLabel = 'Admin@'+process.env.PEER_DOMAIN;
        const identity = X509WalletMixin.createIdentity('Org1MSP', cert, key);

        await wallet.import(identityLabel, identity);

    } catch (error) {
        console.log(`Error adding to wallet. ${error}`);
        console.log(error.stack);
    }
}

main().then(() => {
    console.log('done');
}).catch((e) => {
    console.log(e);
    console.log(e.stack);
    process.exit(-1);
});