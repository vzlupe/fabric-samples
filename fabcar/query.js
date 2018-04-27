'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Chaincode query
 */

var Fabric_Client = require('fabric-client');
var path = require('path');
var util = require('util');
var os = require('os');
var fs = require('fs');

//
var fabric_client = new Fabric_Client();

// setup the fabric network
var channel = fabric_client.newChannel('mychannel');
var peer = fabric_client.newPeer('grpc://localhost:7051');
channel.addPeer(peer);

//
var member_user = null;
var store_path = path.join(__dirname, 'hfc-key-store');
console.log('Store path:'+store_path);
var tx_id = null;

// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
Fabric_Client.newDefaultKeyValueStore({ path: store_path
}).then((state_store) => {
	// assign the store to the fabric client
	fabric_client.setStateStore(state_store);
	var crypto_suite = Fabric_Client.newCryptoSuite();
	// use the same location for the state store (where the users' certificate are kept)
	// and the crypto store (where the users' keys are kept)
	var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
	crypto_suite.setCryptoKeyStore(crypto_store);
	fabric_client.setCryptoSuite(crypto_suite);

	// get the enrolled user from persistence, this user will sign all requests
	return fabric_client.getUserContext('user1', true);
}).then((user_from_store) => {
	if (user_from_store && user_from_store.isEnrolled()) {
		console.log('Successfully loaded user1 from persistence');
		member_user = user_from_store;
	} else {
		throw new Error('Failed to get user1.... run registerUser.js');
	}

	// queryTemp chaincode function - requires 1 argument, ex: args: ['TEMP4'],
	// queryAllTemps chaincode function - requires no arguments , ex: args: [''],
	const request = {
		//targets : --- letting this default to the peers assigned to the channel
		chaincodeId: 'fabcar',
		fcn: 'queryAllTemps',
		args: ['']
	};

	// send the query proposal to the peer
	return channel.queryByChaincode(request);
}).then((query_responses) => {
	console.log("Query has completed, checking results");
	// query_responses could have more than one  results if there multiple peers were used as targets
	if (query_responses && query_responses.length == 1) {
		if (query_responses[0] instanceof Error) {
			console.error("error from query = ", query_responses[0]);
		} else {
			processQuery(query_responses[0]);
			//console.log("Response is ", query_responses[0].toString());
		}
	} else {
		console.log("No payloads were returned from query");
	}
}).catch((err) => {
	console.error('Failed to query successfully :: ' + err);
});

function processQuery(query_responses) {
	var qr = query_responses.toString();
	var qrar = qr.split("},{"); //split response into individual transactions
	var san = qrar[0].substring(2); //remove characters from first transaction
	qrar[0] = san;
	san = qrar[qrar.length - 1].substring(0, qrar[qrar.length - 1].length - 2); //remove characters from last transaction
	qrar[qrar.length -1] = san;
	
	for (var i = 0; i < qrar.length; i++) { //process each transaction to split it into keys and records
		var qrar2 = processLine(qrar[i]);
		qrar[i] = qrar2;
	}
	//remove first transaction
	qrar.splice(0,1);

	writeResults(qrar);
}

function processLine(line) {
	var kv = line.split(", ");
	var k = kv[0].substring(6);
	var k2 = k.substring(1, k.length - 1);
	var v = kv[1].substring(10, kv[1].length - 1);
	var v2 = v.split(",");
	for (var i = 0; i < v2.length; i++) {
		var j = v2[i].split(":");
		var j2 = j[0].substring(1, j[0].length - 1);
		j[0] = j2;
		j2 = j[1].substring(1, j[1].length - 1);
		j[1] = j2;
		v2[i] = j;
	}
	return [k2,v2];
}

function writeResults(qrar) {
	var ws = fs.createWriteStream("testResults/results.txt"); //file to output results to
	var count = 0;
	var fcount = 0;
	for (var i = 0; i < qrar.length; i++) {
		var resp = qrar[i][1][3][1]
		if (resp === "unsuccessful-peer0" || resp === "unsuccessful-peer1" || resp == "unsuccessful-peer2") {
			fcount += 1;
		}
		count += 1;
	}
	ws.write("Found " + fcount + " faults out of " + count + " records");
}
