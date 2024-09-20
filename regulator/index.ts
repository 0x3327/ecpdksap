import { IncrementalMerkleTree } from 'merkletreejs';
import * as fs from 'fs';
import { sha256 } from 'js-sha256';

import API from "./api";

async function testing() {
    const usersFile: string = "registered_users";
    const merkleTreeFile: string = "merkle_tree.json";
    const regulatorAPI: API = new API("localhost", 5555, usersFile, merkleTreeFile);
    await regulatorAPI.start();
    return regulatorAPI;
}

testing().then((res) => {
    console.log("Server sucessfully started"); 
});


