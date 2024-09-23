import { Model, Sequelize } from 'sequelize-typescript';
import { DataTypes } from 'sequelize';

export const init = (sequelize: Sequelize) => sequelize.define<Model>('register_account', {
    pid: {
        primaryKey: true,
        type: DataTypes.INTEGER,
        allowNull: false,
    },
    name: {
        type: DataTypes.STRING,
        allowNull: false,    
    },
    privateKey: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    publicKeyX: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    publicKeyY: {
        type: DataTypes.STRING,
        allowNull: false,
    }
}, {
    timestamps: false
});