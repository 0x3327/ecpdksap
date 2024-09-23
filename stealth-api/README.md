# Stealth API

This document provides an overview of the `API`, `BlockchainService` and `GoHandler` classes, which manage blockchain interactions and WebAssembly-based cryptographic operations.

## Table of Contents
- [Introduction](#introduction)
- [API Overview](#api-overview)
- [BlockchainService](#blockchainservice)
  - [Methods](#blockchainservice-methods)
- [GoHandler](#gohandler)
  - [Methods](#gohandler-methods)

## Introduction
This project is designed to interact with Ethereum-based smart contracts and handle cryptographic operations with WebAssembly. The main components include:

- **API**: Provide endpoints that allows users to interact with the blockchain and perdorm cryptographic operations via WebAssembly. 
- **BlockchainService**: Manages interactions with the Ethereum blockchain, smart contracts, and event listeners.
- **GoHandler**: Interfaces with a WebAssembly (Wasm) module to perform cryptographic operations related to stealth addresses and meta-addresses.

---

## API Overview

The API is designed to provide endpoints that allow users to interact with the blockchain and perform cryptographic operations via WebAssembly. Below is a brief overview of the main API functionality.