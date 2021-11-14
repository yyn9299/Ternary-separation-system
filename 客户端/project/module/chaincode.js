const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const path = require('path');
const router = require('koa-router')();
const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('./AppUtil');

const channelName = 'mychannel'; //
const chaincodeName = 'basic';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');

async function getcontract(username){
	const ccp = buildCCPOrg1();
	const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');
	const wallet = await buildWallet(Wallets, walletPath);
	//await enrollAdmin(caClient, wallet, mspOrg1);
	//await registerAndEnrollUser(caClient, wallet, mspOrg1, org1UserId, 'org1.department1');
	try{
		var gateway = new Gateway();
		await gateway.connect(ccp, {
			wallet,
			identity: username,
			discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
		});
		var network = await gateway.getNetwork(channelName);
		var contract = network.getContract(chaincodeName);
		return contract;
	}catch(error){
		console.error(`Error in connecting to gateway for Org1: ${error}`);
	}
}

async function getgateway(username){
	const ccp = buildCCPOrg1();
	const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com');
	const wallet = await buildWallet(Wallets, walletPath);
	try{
		const gateway = new Gateway();
		await gateway.connect(ccp, {
			wallet,
			identity: username,
			discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
		});
		return gateway;
	}catch(error){
		console.error(`Error in connecting to gateway for Org1: ${error}`);
	}
}

module.exports = getcontract;
//module.exports = getgateway;