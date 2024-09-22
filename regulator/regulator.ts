import { IncrementalMerkleTree } from 'merkletreejs';
import * as fs from 'fs';
import { sha256 } from 'js-sha256';

const { buildPoseidon } = require('circomlibjs');

class Regulator {
    private _poseidon: any;
    private registeredUsersFilePath: string;
    private merkleTreeFilePath: string;
    public tree!: IncrementalMerkleTree;


    constructor(registeredUsersFilePath: string, merkleTreeFilePath: string) {
        this.registeredUsersFilePath = registeredUsersFilePath;
        this.merkleTreeFilePath = merkleTreeFilePath;
    }

    /**
     * Hash function used in Merkle tree 
     * 
     * @param inputs  - array of objects we want to hash (objects must be convertible to BigInt)
     * @returns {BigInt} - the resulting hash
     */
    private poseidon(inputs: any[]): BigInt {
        const bigIntInputs = inputs.map(IncrementalMerkleTree.bigNumberify);
        const hash = this._poseidon(bigIntInputs);
        const bn = IncrementalMerkleTree.bigNumberify(this._poseidon.F.toString(hash))
        return bn
    }

    /**
     * Generates hash of the user 
     * 
     * @param {string} name - name of the user
     * @param {number} pid - user's personal identification number (PID)
     * @param {BigInt} pub_x - x coordinate of user's public key
     * @param {BigInt} pub_y - y coordinate of user's public key
     * @returns {BigInt} - hash of the user
     */
    private generateUserHash(name: string, pid: number, pub_x: BigInt, pub_y: BigInt): BigInt {
        const name_hash_hex = sha256(name);
        const userHash = this.poseidon(['0x' + name_hash_hex, pid, pub_x, pub_y]);
        return userHash;
    }

    /**
     * Checks if the user is already registered
     * 
     * @returns {number} - index of the user in the tree if registered, -1 otherwise
     */
    private checkIfRegistered(name: string, pid: number, pub_x: BigInt, pub_y: BigInt): number {
        if (!fs.existsSync(this.registeredUsersFilePath)) {
            return -1;
        }

        const fileData = fs.readFileSync(this.registeredUsersFilePath, 'utf-8');
        const lines = fileData.split('\n');

        for (const line of lines) {
            const data = line.split(' ');
            if (data.length > 1 && Number(data[1]) === pid) {
                return parseInt(data[2], 10);
            }
        }
        return -1;
    }

    /**
     * Initialization of poseidon hash function and Merkle Tree
     */
    public async init() {
        this._poseidon = await buildPoseidon();
        this.tree = new IncrementalMerkleTree(this.poseidon.bind(this), {
            depth: 2,
            arity: 2,
            zeroValue: BigInt(0)
        });
    }

    /**
     * Registers the user
     * @returns {number} - index of the user in the tree
     */
    public registerUser(name: string, pid: number, pub_x: BigInt, pub_y: BigInt): number {

        // user is already registered
        let index = this.checkIfRegistered(name, pid, pub_x, pub_y);
        if (index !== -1)
            return index; 
    
        const userHash = this.generateUserHash(name, pid, pub_x, pub_y);
        this.tree.insert(userHash);

        // Write user into file
        index = this.tree.indexOf(userHash);
        fs.appendFileSync(this.registeredUsersFilePath, `${name} ${pid} ${index}\n`);

        // console.log('User hash:', userHash.toString(16)); // Debug output
        // console.log('User inserted at index:', index); // Debug output

        return index;
    }

    /**
     * Gets hashes needed to construct Merkle proof for the user at leafIndex
     * @param {number} leafIndex - index of the user in the Merkle tree
     * @returns {any[]} - array of hashes needed for the proof and their indexes, 
     * each entry of this array is another array of the form [hash, index]
     */
    public getProofForUser(leafIndex: number): any[] {
        // TODO: check if the leafIndex argument is given
        let hashes = [];
        if (this.tree !== undefined) {
            const proof = this.tree.getProof(leafIndex);
            for (let i = 0; i < proof.siblings.length; i++) {
                hashes.push([proof.siblings[i][0].toString(), proof.pathIndices[i]]);
            }
        }
        return hashes;
    }

    /**
     * Saves merkle tree into the given file 
     * TODO: update function to work with trees of arbitrary depth
     */
    public saveTreeToFile(): void {
        if (this.tree !== undefined) {
            const leaves = this.tree.getLeaves().map((leaf) => leaf.toString(16));
            const root = this.tree.getRoot().toString(16);
            const parents: string[] = [];
    
            // Calculating parents hash
            for (let i = 0; i < leaves.length; i += 2) {
                if (i + 1 < leaves.length) {
                    const leftLeaf = BigInt('0x' + leaves[i]);
                    const rightLeaf = BigInt('0x' + leaves[i + 1]);
                    const parentHash = this.poseidon([leftLeaf, rightLeaf]).toString(16);
                    parents.push(parentHash);
                }
            }
    
            const treeData = {
                leaves: leaves,
                parents: parents,
                root: root
            };
    
            // Convert to JSON and save
            fs.writeFileSync(this.merkleTreeFilePath, JSON.stringify(treeData, null, 2));
            console.log(`Merkle tree with parents saved to ${this.merkleTreeFilePath}`);
        }
    }

    /**
     * Constructs merkle tree from the given file 
     * TODO: update function to work with trees of arbitrary depth
     */
    public loadTreeFromFile(): void {
        if (fs.existsSync(this.merkleTreeFilePath)) {
            const fileData = fs.readFileSync(this.merkleTreeFilePath, 'utf-8');
            const treeData = JSON.parse(fileData);

            this.tree = new IncrementalMerkleTree(this.poseidon.bind(this), {
                depth: 2,
                arity: 2,
                zeroValue: BigInt(0)
            });
    
            // Add leaves
            treeData.leaves.forEach((leaf: string) => {
                if (leaf !== '0')
                    this.tree?.insert(BigInt('0x' + leaf));
            });
        } else {
            console.log(`File ${this.merkleTreeFilePath} does not exist.`);
        }
    }
}

export default Regulator;