import express, { Express, Request, Response } from 'express';
import dotenv from "dotenv";

dotenv.config();

const app: Express = express();
const port = process.env.REGULATOR_PORT;

app.get("/", (req: Request, res: Response) => {
    res.send("Regulatory body");
});

app.use((req: Request, res: Response) => {
    res.status(404).send("Error 404: Not Found");
});

app.listen(port, () => {
    console.log(`[Regulatory Body]: server started listening on http://localhost:${port}`);
});