# Curvy Stealth Address Protocol (formerly ECPDKSAP)

![](./docs/assets/demo-banner.png)

Website with DEMO available at: https://0xcurvy.io

## Overview

The repository contains implementation code for all [ECPDKSAP protocol versions](./docs).

All protocol variations are implemented in Go programming language and use [consensys/gnark-crypto](https://github.com/Consensys/gnark-crypto) library.

Forked version: [0x3327/gnark-crypto](https://github.com/0x3327/gnark-crypto) containing custom function implementations needed by ECPDKSAP protocols.

Smart contracts are developed using Foundry framework and follow the [EIP-5564](https://eips.ethereum.org/EIPS/eip-5564) standard.

## Project structure

Each of the following sub-directories contain their own documentation specifications.

The project's dir. structure is as following:

- `./docs`:  General documentation, with detailed results

- `./impl`: Off-chain protocol code (used by sender and recipient to generate/scan data)
- `./stealth-api`: On-chain contracts (used for sender - recipient "communication")
- `./ft`: front-end, client - oriented

## Additional resources



 **Elliptic Curve Pairing Stealth Address Protocols**
Marija Mikic, Mihajlo Srbakoski https://arxiv.org/abs/2312.12131
