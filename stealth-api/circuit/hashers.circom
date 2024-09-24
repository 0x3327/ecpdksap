pragma circom 2.1.6;

include "../node_modules/circomlib/circuits/poseidon.circom";
// include "https://github.com/0xPARC/circom-secp256k1/blob/master/circuits/bigint.circom";

template HashLeftRight() {
    signal input left;
    signal input right;

    signal output hash;

    component hasher = Poseidon(2);
    left ==> hasher.inputs[0];
    right ==> hasher.inputs[1];

    hash <== hasher.out;
}


