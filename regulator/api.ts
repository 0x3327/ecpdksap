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
            const { name, pid, pub_x, pub_y } = req.body;

            // TODO: check if the tree is full already
            let index = this.regulator.registerUser(name, Number(pid), BigInt(pub_x), BigInt(pub_y));

            this.regulator.saveTreeToFile();

            let responseData = this.regulator.getProofForUser(index);
            responseData.push([this.regulator.tree.getRoot().toString()]);
            // TODO: also add signed merkle root as response
            sendResponseOK(res, "Handling /register-user", responseData);
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