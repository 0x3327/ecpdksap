const { secp256k1 } = require("ethereum-cryptography/secp256k1");
const { keccak256 } = require("ethereum-cryptography/keccak");
const { toHex } = require("ethereum-cryptography/utils");
const EthCrypto = require("eth-crypto");

const pk = "0x9d1ffb1fc8b377ad4818725016113ca0572f37cf8a10ac9db18b19ceff8a7d02";

const publicKey = EthCrypto.publicKeyByPrivateKey(pk);

const address = EthCrypto.publicKey.toAddress(publicKey);

console.log({ publicKey, address });
