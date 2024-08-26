import { DbConfig } from "../../../types";
import { ModelCtor, Sequelize }  from 'sequelize-typescript';
import { init as receivedTransactionsModelInit } from "./models/received-transaction-model";
import { init as sentTransactionsModelInit } from "./models/sent-transaction-model";


class DB {
    config: DbConfig;
    sequelize: Sequelize;
    models!: {
        receivedTransactions: any,
        sentTransactions: any,
    };

    constructor(dbConfig: DbConfig) {
        this.config = dbConfig;
        this.sequelize = new Sequelize({
            username: this.config.username,
            password: this.config.password,
            dialect: 'sqlite',
            storage: 'db.sqlite',
            logging: false,
        })

        this.registerModels();
    }

    registerModels(): void {
        const models = {
            receivedTransactions: receivedTransactionsModelInit(this.sequelize),
            sentTransactions: sentTransactionsModelInit(this.sequelize),
        };
        
        this.models = models;
    }
}

export default DB;