pragma circom 2.1.6;

include "circomlib/poseidon.circom";
// include "https://github.com/0xPARC/circom-secp256k1/blob/master/circuits/bigint.circom";

template NullifierMaker () {
    signal input publicVar;
    signal input hashName;
    signal input pid;
    signal output nullifier;
    
    component hash = Poseidon(3);
    hash.inputs[0] <== hashName;
    hash.inputs[1] <== pid;
    hash.inputs[2] <== publicVar;
    
    nullifier <== hash.out;

}