/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

module.exports.info  = 'small_bank_operations';

let bc, contx;
let no_accounts = 0;
let account_array = [];
let accounts, txnPerBatch;
const initial_balance = 1000000;
const operation_type = ['transact_savings','deposit_checking','send_payment','write_check', 'amalgamate'];
let prefix;

/**
 * Get account index
 * @return {Number} index
 */
function getAccount() {
    return Math.floor(Math.random()*Math.floor(account_array.length));
}

/**
 * Get two accounts
 * @return {Array} index of two accounts
 */
function get2Accounts() {
    let idx1 = getAccount();
    let idx2 = getAccount();
    if(idx2 === idx1) {
        idx2 = getAccount();
    }
    return [idx1, idx2];
}

/**
 * Generate unique account key for the transaction
 * @returns {Number} account key
 **/
function generateAccount() {
    // should be [a-z]{1,9}
    if(typeof prefix === 'undefined') {
        prefix = process.pid;
    }
    let count = account_array.length+1;
    let num = prefix.toString() + count.toString();
    return parseInt(num);
}

/**
 * Generates random string.
 * @returns {string} random string from possible characters
 **/
function random_string() {
    let text = '';
    const possible = 'ABCDEFGHIJKL MNOPQRSTUVWXYZ abcdefghij klmnopqrstuvwxyz';

    for (let i = 0; i < 12; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
}

/**
 * Generates small bank workload with specified number of accounts
 * and operations.
 * @returns {Object} array of json objects and each denotes
 * one operations
 **/
function generateWorkload() {
    let workload = [];
    for(let i= 0; (i < txnPerBatch && no_accounts < accounts); i++,no_accounts++) {
        let acc_id = generateAccount();
        account_array.push(acc_id);
        let acc = {
            'customer_id': acc_id,
            'customer_name': random_string(),
            'initial_checking_balance': initial_balance,
            'initial_savings_balance': initial_balance,
            'transaction_type': 'create_account'
        };
        workload.push(acc);
    }
    for(let j= workload.length; j<txnPerBatch; j++) {
        let op_index =  Math.floor(Math.random() * Math.floor(operation_type.length));
        let acc_index = getAccount();
        let random_op = operation_type[op_index];
        let random_acc = account_array[acc_index];
        let amount = Math.floor(Math.random() * 200);
        let op_payload;
        switch(random_op) {
        case 'transact_savings': {
            op_payload = {
                'amount': amount,
                'customer_id': random_acc,
                'transaction_type':random_op
            };
            break;
        }
        case 'deposit_checking': {
            op_payload = {
                'amount': amount,
                'customer_id': random_acc,
                'transaction_type':random_op
            };
            break;
        }
        case 'send_payment': {
            let accounts = get2Accounts();
            op_payload = {
                'amount': amount,
                'dest_customer_id': account_array[accounts[0]],
                'source_customer_id': account_array[accounts[1]],
                'transaction_type': random_op
            };
            break;
        }
        case 'write_check': {
            op_payload = {
                'amount': amount,
                'customer_id': random_acc,
                'transaction_type':random_op
            };
            break;
        }
        case 'amalgamate': {
            let accounts = get2Accounts();
            op_payload = {
                'dest_customer_id': account_array[accounts[0]],
                'source_customer_id': account_array[accounts[1]],
                'transaction_type': random_op
            };
            break;
        }
        default: {
            throw new Error('Invalid operation!!!');
        }
        }
        workload.push(op_payload);
    }
    return workload;
}

module.exports.init = function(blockchain, context, args) {
    if(!args.hasOwnProperty('accounts')) {
        return Promise.reject(new Error('smallbank.operations - \'accounts\' is missed in the arguments'));
    }
    if(!args.hasOwnProperty('txnPerBatch')) {
        return Promise.reject(new Error('smallbank.operations - \'txnPerBatch\' is missed in the arguments'));
    }
    accounts = args.accounts;
    if(accounts <= 3) {
        return Promise.reject(new Error('smallbank.operations - number accounts should be more than 3'));
    }
    txnPerBatch = args.txnPerBatch;
    bc = blockchain;
    contx = context;

    return Promise.resolve();
};

module.exports.run = function() {
    let args = generateWorkload();

    // rearrange arguments for the Fabric adapter
    if (bc.getType() === 'fabric') {
        args = args.map(arg => {
            return {
                chaincodeFunction: arg.transaction_type,
                // need to remove the key for the TX type
                chaincodeArguments: Object.keys(arg).filter(k => k !== 'transaction_type').map(k => arg[k].toString()),
            };
        });
    }

    return bc.invokeSmartContract(contx, 'smallbank', '1.0', args, 30);
};

module.exports.end = function() {
    return Promise.resolve();
};


module.exports.account_array = account_array;
