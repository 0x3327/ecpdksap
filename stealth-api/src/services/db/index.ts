import { DbConfig } from "../../../types";
import { Sequelize }  from 'sequelize-typescript';
import { init as receivedTransactionsModelInit } from "./models/received-transaction-model";
import { init as sentTransactionsModelInit } from "./models/sent-transaction-model";
import { init as registerAccountModelInit } from "./models/register-account-model";



class DB {
    config: DbConfig;
    sequelize: Sequelize;
    models!: {
        receivedTransactions: any,
        sentTransactions: any,
        registerAccount: any,
    };

    constructor(dbConfig: DbConfig) {
        this.config = dbConfig;
        this.sequelize = new Sequelize({
            username: this.config.username,
            password: this.config.password,
            dialect: 'sqlite',
            storage: `${dbConfig.database || 'db'}.sqlite`,
            logging: false,
        })

        this.registerModels();
    }

    registerModels(): void {
        const models = {
            receivedTransactions: receivedTransactionsModelInit(this.sequelize),
            sentTransactions: sentTransactionsModelInit(this.sequelize),
            registerAccount: registerAccountModelInit(this.sequelize),
        };
        
        this.models = models;
    }
}

export default DB;