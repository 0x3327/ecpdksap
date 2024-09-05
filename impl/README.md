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
  - Format

    ```json
    {
      //Sender's private key
      "r": string,

      //Recipient's public spending key
      "K": string,

      //Recipient's public viewing key
      "V": "XAffineCoord.YAffineCoord", // note the `.` separator

      //Protocol Version
      "Version": string, // v0, v1, v2

      //View tag being used
      "ViewTagVersion": string // v0-1byte, v0-2bytes, v1-1byte
    }
    ```

  - For example:
    ```bash
    export SND_INPUT=$(cat ./gen_example/example/inputs/send.json) \
    && go run . send $SND_INPUT
    ```

- `receive-scan < jsonString >`

  - called on the recipient's side to check for incoming ETH transfers
  - `jsonString` is a text string containing all necessary parameters (see: [example-receive-input](./gen_example/example/inputs/receive.json))
  - Format:

    ```json
    {
      //Recipient's private spending key
      "k": string,

      //Recipient's private viewing key
      "v": string,

      //List of Senders' public keys
      "Rs": ["Rj_AffineXCoord.Rj_AffineYCoord"], // note the `.` separator

      //List of corresponding view tags
      "ViewTags": [] string, // hexadecimal string: 1 or 2 byte long

      //Protocol Version
      "Version": string, // v0, v1, v2

      //View tag being used
      "ViewTagVersion": string // v0-1byte, v0-2bytes, v1-1byte
    }
    ```

  - In example:
    ```bash
    export RCV_INPUT=$(cat ./gen_example/example/inputs/receive.json) \
    && go run . receive-scan $RCV_INPUT
    ```

- `gen-example < version: v0 | v1 | v2 > < sample-size: uint >`
  - generates input examples for the sender's recipient's side
  - `< version: v0 | v1 | v2 >` refers to the protocol versions
  - `< view-tag-version: v0-1byte | v0-2bytes | v1-1byte >` refers to the version of the view tag being used
  - `< sample-size: uint >` number of senders' public keys

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
