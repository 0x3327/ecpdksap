import express, { Express, Request, Response } from "express";
import Regulator from "./regulator";
import bodyParser from "body-parser";

const sendResponse = (res: Response, status: number, message: string, data?: any) => {
    res.status(status);
    const response = {
        message,
        status,
        data,
    };
    res.send(response);
};

const sendResponseOK = (res: Response, message: string, data?: any) => {
    sendResponse(res, 200, message, data);
};

const sendResponseBadRequest = (res: Response, message: string, data?: any) => {
    sendResponse(res, 400, message, data);
};

class API {
    private host: string;
    private port: number;
    private server: Express;
    private regulator: Regulator;

    constructor(host: string, port: number, registeredUsersFilePath: string, merkleTreeFilePath: string) {
        this.host = host;
        this.port = port;
        this.regulator = new Regulator(registeredUsersFilePath, merkleTreeFilePath);

        this.server = express();
        this.server.use(bodyParser.json());

        this.exposeRoutes();
    }

    private exposeRoutes() {
        this.server.get('/', (req: Request, res: Response) => {
            sendResponseOK(res, "Service running", {timestamp: Date.now()});
        });

        this.server.post('/register-user', (req: Request, res: Response) => {
            const { name, pid, publicKeyX, publicKeyY } = req.body;

            // TODO: check if the tree is full already
            let index = this.regulator.registerUser(name, Number(pid), BigInt(publicKeyX), BigInt(publicKeyY));

            this.regulator.saveTreeToFile();

            let responseData = this.regulator.getProofForUser(index);
            // TODO: 
            //      - also add signed merkle root as response
            //      - send proper json format back
            
            let proof = {
                hashes: [] as {hash: string, index: number }[],
                root: null,
                signedRoot: null
            };
            for (let i = 0; i < responseData.length; i++) {
                proof["hashes"].push({ 
                    hash: responseData[i][0], 
                    index: responseData[i][1] 
                });
            }
            proof["root"] = this.regulator.tree.getRoot().toString();
            // console.log("Sending...");
            // console.log(proof);
            sendResponseOK(res, "Handling /register-user", proof);
        });

        // handling unspecified routes
        this.server.use((req: Request, res: Response) => {
            sendResponseBadRequest(res, "Specified path doesn't exist");
        });
    }

    public async start() {
        await this.regulator.init();
        this.regulator.loadTreeFromFile();
        this.server.listen(this.port, this.host, () => {
            console.log(`Started listening on ${this.host}:${this.port}`);
        }) 
    }
}

export default API;