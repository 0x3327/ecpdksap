## Overview

This subdirectory represents reference **ECPDKSAP** implementation.

This is a standard `go` module that makes use of a custom CLI command subsystem.

## Usage Details

Using either `go run .` command prefix or the binary file(located in `./builds`), there exist following subcommands:

- `bench < only-bn254 | all-curves >`

  - with `only-bn254` benchmarking the optimized code version for the BN254 curve
  - and `all-curves` benchmarking general implementation across 6 different curves (BLS12-377, BLS12-381, BLS24-315, BN254, BW6-633, BW6-761)

- `send < jsonString >`

  - called before sending ETH to a stealth address
  - `jsonString` is a text string containing all necessary parameters (see: [example-send-input](./gen_example/example/inputs/send.json))
  - For example:
    ```
    export SEND_EXAMPLE_INPUT=$(cat ./gen_example/example/inputs/send.json) && go run . send $SEND_EXAMPLE_INPUT
    ```

- `recive-scan < jsonString >`

  - called on the recipient's side to check for incoming ETH transfers
  - `jsonString` is a text string containing all necessary parameters (see: [example-receive-input](./gen_example/example/inputs/receive.json))
  - In example:
    ```
    export RECEIVE_EXAMPLE_INPUT=$(cat ./gen_example/example/inputs/receive.json) && go run . receive-scan $RECEIVE_EXAMPLE_INPUT
    ```

- `gen-example < version: v0 | v1 | v2 > < sample-size: 1...1000 >`
  - generates input examples for the sender's recipient's side
  - `< version: v0 | v1 | v2 >` refers to the protocol versions
  - `< sample-size: 1...1000 >` number of senders' public keys

## Directory structure

- `./benchmark`:
  - used for benchmarking results
    - BLS12-377, BLS12-381, BLS24-315, BN254, BW6-633, BW6-761 curves comparison
    - optimized code version for the best **BN254** curve
- `./builds`:
  - contains different binary code versions of the entire module
- `./gen_example`:
  - helper submodule that generates example inputs to be used via CLI
- `./gnark-crypto-fork`:
  - forked version of [consensys/gnark-crypto]() with added specialized methods required by ECPDKSAP
- `./recipient`:
  - contains code for the recipient's side (triggered via CLI)
- `./sender`:
  - contains code for the sender's side (triggered via CLI)
- `./versions`:
  - implementations of three different protocol versions (v0..v2)
