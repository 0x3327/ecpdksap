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
    stealthAddress: {
        type: DataTypes.STRING,
        allowNull: false,
    },
    ephemeralKey: {
        type: DataTypes.STRING,
        allowNull: false,
    }
}, {
    timestamps: false
});