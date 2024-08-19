import configLoader from  './config-loader';
import App from './app';

const config = configLoader.load();
const app = new App(config);

app.start().then(() => {
    console.log('Application started successfully.');
}).catch((err) => {
    console.error('Failed to start the application:', err);
});