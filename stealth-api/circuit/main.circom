pragma circom 2.1.6;

include "../node_modules/circomlib/circuits/poseidon.circom";
include "./merkleTree.circom";
include "./hashers.circom";
include "./nullifier.circom";
include "./privToPub.circom";

template inclusionProof (levels) {
    signal input hashName;
    signal input pid;
    signal input privKey;

    component privToPub = PrivToPubKey();
    privToPub.privKey <== privKey;
    signal pubKey[2] <== privToPub.pubKey;


    component hash = Poseidon(4);
    hash.inputs[0] <== hashName;
    hash.inputs[1] <== pid;
    hash.inputs[2] <== pubKey[0];
    hash.inputs[3] <== pubKey[1];

    signal leaf <== hash.out;

    signal input path_elements[levels];
    signal input path_index[levels];

    signal input publicVar;
    signal input root;

    signal output nullifier;
    
    component null = NullifierMaker();
    null.publicVar <== publicVar;
    null.hashName <== hashName;
    null.pid <== pid;
    null.privKey <== privKey;

    nullifier <== null.nullifier;

    component inclusionTest = LeafExists(levels);
    inclusionTest.leaf <== leaf;
    inclusionTest.path_elements <== path_elements;
    inclusionTest.path_index <== path_index;
    inclusionTest.root <== root;

}

component main { public [ publicVar, root ] } = inclusionProof(2);

/* INPUT = {
    "hashName": "12345",
    "pid": "123456",
    "privKey": "7",
    "path_elements": [
        "12",
        "34"
        ],
    "path_index": [
        "1", 
        "0"
        ],
    "publicVar": "10000",
    "root": "77"
} */
