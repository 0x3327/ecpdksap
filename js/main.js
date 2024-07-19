const { secp256k1 } = require("ethereum-cryptography/secp256k1");
const { keccak256 } = require("ethereum-cryptography/keccak");
const { toHex } = require("ethereum-cryptography/utils");
const EthCrypto = require("eth-crypto");

const pk = "0xbc126c60e95237fe11ecf3321a079e0b4c8686a08920dd7ee2d93cf9aa5c2bb6";

console.log(`pk.len`, pk.length);

const publicKey = EthCrypto.publicKeyByPrivateKey(pk);

const address = EthCrypto.publicKey.toAddress(publicKey);

console.log({ publicKey, address });
