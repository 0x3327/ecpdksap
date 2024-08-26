import { Model, Sequelize } from 'sequelize-typescript';
import { DataTypes } from 'sequelize';

export const init = (sequelize: Sequelize) => sequelize.define<Model>('sent_transactions', {
    id: {
        primaryKey: true,
        type: DataTypes.INTEGER,
    },
    transaction_hash: {
        type: DataTypes.STRING,
        allowNull: false,    
    },
    block_number: {
        type: DataTypes.INTEGER,
    },
    amount: {
        type: DataTypes.INTEGER,
        allowNull: false,
    },
    recipient_identifier: {
        type: DataTypes.STRING,
        allowNull: true,
    },
    recipient_identifier_type: {
        type: DataTypes.STRING,
        allowNull: true,
    },
    recipient_k: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    recipient_v: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    recipient_stealth_address: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    ephemeral_key: {
        type: DataTypes.STRING,
    }
}, {
    timestamps: false
});