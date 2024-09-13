import { Model, Sequelize } from 'sequelize-typescript';
import { DataTypes } from 'sequelize';

export const init = (sequelize: Sequelize) => sequelize.define<Model>('received_transactions', {
    id: {
        primaryKey: true,
        type: DataTypes.INTEGER,
    },
    transaction_hash: {
        type: DataTypes.STRING,
        allowNull: false,
        unique: true,
    },
    block_number: {
        type: DataTypes.INTEGER,
        allowNull: false,
    },
    amount: {
        type: DataTypes.INTEGER,
        allowNull: false,
    },
    stealth_address: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    ephemeral_key: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    view_tag: {
        type: DataTypes.STRING,
        allowNull: true,
    }
}, {
    timestamps: false
});