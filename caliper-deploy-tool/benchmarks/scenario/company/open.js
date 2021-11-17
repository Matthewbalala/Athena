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

module.exports.info  = 'opening accounts';

let account_array = [];
let bc, contx;
let companys = [];
let txnPerBatch;
let post = 0;
module.exports.init = function(blockchain, context, args) {
const fs = require('fs');
let rawdata = fs.readFileSync('/hyperledger/caliper/workspace/benchmarks/scenario/company/companydata/company2.json');
var company = rawdata.toString().split("\n")
for (let index = 0; index < company.length-1; index++) {
    const element = company[index];
    var info =JSON.parse(element)
    // console.log(JSON.parse(element))
    companys.push(info);
}
    txnPerBatch = args.txnPerBatch;
    bc = blockchain;
    contx = context;
    return Promise.resolve();
};


function getRandomCinicWorkload(n){
    var fewcompanys=[],id,leng=companys.length;
    //;console.log(leng,"11111111")
    for(var i=0;i<n;i++){
        id=Math.floor(Math.random()%leng)
        if (bc.getType() === 'fabric') {
        fewcompanys.push({
            chaincodeFunction: 'submit_company',
            chaincodeArguments: [companys[id].short_name.toString() + '_' + post.toString(), JSON.stringify(companys[id])],
        });
        post++;
        }
    //  console.log(fewcompanys,"33333333")
    }

    return fewcompanys;
}


module.exports.run = function() {
    let args = getRandomCinicWorkload(txnPerBatch);
    return bc.invokeSmartContract(contx, 'company', 'v0', args, 100);
};

module.exports.end = function() {
    return Promise.resolve();
};

module.exports.account_array = account_array;
