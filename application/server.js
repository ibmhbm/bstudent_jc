// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { Wallets, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

// Constants
const PORT = 8080;
const HOST = "0.0.0.0";

// use static file
app.use(express.static(path.join(__dirname, "views")));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
    res.sendFile(__dirname + "/index.html");
});

async function cc_call(fn_name, args){
    // chaincode 호출 파트
    // load the network configuration
    const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
    let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
    // Create a new file system based wallet for managing identities.
    const walletPath = path.join(process.cwd(), "wallet");
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);
    // Check to see if we've already enrolled the user.
    const identity = await wallet.get("appUser");
    if (!identity) {
        console.log(
            'An identity for the user "appUser" does not exist in the wallet'
        );
        console.log("Run the registerUser.js application before retrying");
        return;
    }
    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: "appUser",
        discovery: { enabled: true, asLocalhost: true },
    });
    // Get the network (channel) our contract is deployed to.
    const network = await gateway.getNetwork("mychannel");
    // Get the contract from the network.
    const contract = network.getContract("mJongchin");
    // Submit the specified transaction.  
    if (fn_name == "addAsset") {
        await contract.submitTransaction("addAsset", args[0],args[1],args[2],args[3],args[4]);       
    } else if (fn_name == "readAsset"){
        result = await contract.evaluateTransaction("readAsset", args);
    } else if (fn_name== "deleteAsset"){
        await contract.submitTransaction("deleteAsset", args);
    } else if (fn_name == "addOwner") {
        await contract.submitTransaction("addOwner", args[0],args[1]);     
    }
    else result = "not supported fuction!"   
    await gateway.disconnect();
    return result
}

app.get("/readAsset", async(req, res) => {
    const name = req.query.name
    console.log("request :" + name)

    result = await cc_call("readAsset",name)  
    console.log("reqest : " + name )
    console.log("response : " + result.toString())

    const result_str =`{"result":"Done /readAsset","msg":${result.toString()}}`
    res.status(200).json(JSON.parse(result_str))
});

app.post("/read", (req, res) => { 
    const id=req.body.assetid
    console.log("request :"+id) 
});

app.post("/addAsset", async(req, res) => { 
    const name=req.body.assetname
    const typeString=req.body.type
    const category=req.body.category 
    const owner=req.body.owner
    const stake=req.body.stake
    var type = parseInt(typeString)
    console.log(typeString)
    console.log("request :"+name+" "+type+" "+category+" "+owner+" "+stake)  
    var args = [name,type,category,owner,stake]
    result = await cc_call("addAsset", args)

    // const resultPath = path.join(process.cwd(), "views/result.html")
    // var resultHTML = fs.readFileSync(resultPath,"utf-8")

});

app.post("/deleteAsset", async(req, res) => { 
    const id=req.body.assetid
    console.log("request :"+id) 
    result = await cc_call("deleteAsset", id)
});

app.post("/addOwner", async(req, res) => { 
    const name=req.body.ownername
    const stake=req.body.stake
    console.log("request :"+name+" "+stake)  
    var args = [name,stake]
    result = await cc_call("addOwner", args)

    // const resultPath = path.join(process.cwd(), "views/result.html")
    // var resultHTML = fs.readFileSync(resultPath,"utf-8")

});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
