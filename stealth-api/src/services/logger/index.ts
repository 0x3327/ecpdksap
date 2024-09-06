import winston from 'winston';
import { Config } from '../../../types';

class LoggerService {
    logger: winston.Logger

    constructor(config: Config) {
        this.logger = winston.createLogger({
            transports: [
                new (winston.transports.Console)(),
            ],
            silent: !config.logging
        })
    }
}

export default LoggerService;