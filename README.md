# Ethereum Vanity Wallet Generator

<br>
<h3 align="center">
  
  <img src="https://user-images.githubusercontent.com/37617738/120125455-13de3a00-c1e3-11eb-9a51-707e2dcefdaa.png" alt="Ethereum Wallet Generator" height="70" />
   &nbsp&nbsp
  <img src="https://user-images.githubusercontent.com/37617738/120122724-aaefc580-c1d4-11eb-9343-234eb8fb3ab9.png" alt="Ethereum Wallet Generator" height="50" />
</h3>
<br/>
<h3 align="center">
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="50" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="60" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="70" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="80" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="90" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="100" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="110" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="120" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="130" />
  <img src="https://user-images.githubusercontent.com/37617738/120087436-1886ed80-c112-11eb-945f-8065957a1dd0.png" alt="go ethereum" height="140" />
</h3>
<br>

> **Blazing fast multiple Ethereum and Crypto Vanity Wallets Generator written in Go💰** <br>Generate a ten thousand crypto wallets (vanity address and mnemonic seed) in a sec ⚡️<br>Find beautiful and awesome vanity wallet addresses 🎨

[![Golang](https://badges.aleen42.com/src/golang.svg)](https://golang.org/)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1d765b63df4b4266bdcf653d5a024458)](https://www.codacy.com/gh/Planxnx/ethereum-wallet-generator/dashboard?utm_source=github.com&utm_medium=referral&utm_content=Planxnx/ethereum-wallet-generator&utm_campaign=Badge_Grade)
[![Docker Image Size (tag)](https://img.shields.io/docker/image-size/planxthanee/ethereum-wallet-generator/latest)](https://hub.docker.com/r/planxthanee/ethereum-wallet-generator)
[![Code Analysis & Tests](https://github.com/planxnx/ethereum-wallet-generator/actions/workflows/code-analysis.yml/badge.svg)](https://github.com/planxnx/ethereum-wallet-generator/actions/workflows/code-analysis.yml)
![GitHub issues](https://img.shields.io/github/issues/Planxnx/ethereum-wallet-generator)
[![DeepSource](https://deepsource.io/gh/Planxnx/ethereum-wallet-generator.svg/?label=active+issues)](https://deepsource.io/gh/Planxnx/ethereum-wallet-generator/?ref=repository-badge)
[![license](https://img.shields.io/badge/license-WTFPL%20--%20Do%20What%20the%20Fuck%20You%20Want%20to%20Public%20License-green.svg)](https://github.com/planxnx/ethereum-wallet-generator/blob/main/LICENSE)

## Easy & Fast Way to create a million beauty Ethereum Vanity Wallets 🎨⚡️

![ethereum and crypto wallets generated](https://user-images.githubusercontent.com/37617738/227807144-c1dc59ae-94fd-4fdf-9678-bf8c12e58cd4.png)

- Awesome and Beautiful vanity wallet addresses supported! [regex, prefix, suffix, contains] 🎨
- Blazing fast wallets generate. Speeding up to **+100k wallet/sec** (/w concurrency and only privatekey mode)⚡️
- Supports to generating until got the vanity wallet addresses you want 🤩 (using `-n -1` and `-limit <number>` flags)
- ∞ Infinite wallet generating! (set number to 0 to active infinite loop) ∞
- Generate word seed phrase with BIP-39 mnemonic (support 12, 24 Word Seed Phrase) (Default is 128 bits for 12 words).
- Embedded Database Supported! (with SQLite3). It's easiest to generate, manage, search a billion wallets without any pain.
- 🔐 **Encrypted Export** - Automatically export wallets to password-encrypted JSON files for secure sharing
- Tiny Sizes and Superior Speed with Golang 🚀 (required go 1.21 or higher)
- No Go? No Problem! [Docker images 🐳](https://hub.docker.com/r/planxthanee/ethereum-wallet-generator) or [exec files](https://github.com/Planxnx/ethereum-wallet-generator/releases/latest) are provided for you
- You can benchmark generating speed by setting the `isDryrun` flag 📈
- Default (HD Wallet)Hierarchical Deterministic Path - m/44'/60'/0'/0 .
- We recommend every user of this application audit and verify every source code in this repository and every imported dependecies for its validity and clearness. 👮🏻‍♂️

## What's a vanity address 🎨

A vanity address is an address which part of it is chosen by yourself. Adding vanity to an address is used to give it personality, to reinforce a brand, to send a message, or to make the owner(s) feel cool 🤩

Examples: `0x1111111254fb6c44bAC0beD2854e76F90643097d`, or `0x999999999aa3d5F44D48729b11c091Ed1f12f599`

## Installation

<img  align="right" src="https://user-images.githubusercontent.com/37617738/120122855-b1cb0800-c1d5-11eb-9502-8d64bb275337.png" height="140" alt="gopher" />

> **Homebrew is coming soon**

<be>

### Install from Source

```console
$ go install github.com/planxnx/ethereum-wallet-generator@latest
```

### Docker

```console
$ docker pull planxthanee/ethereum-wallet-generator:latest
```

### Download from latest release

> supports only Windows x86-64 and macOs

[Download](https://github.com/Planxnx/ethereum-wallet-generator/releases/latest)

## Modes

We've provided 2 modes for you to generate wallets.

- **[1] Normal Mode** - Generate wallets with mnemonic phrase. (default)
- **[2] Only Private Key Mode⚡️** - Generate wallets with private key only. **Increase speed up to 20x (+100k wallet/sec), but you will not get a mnemonic phrase.**

Engine options:
- **CPU** (default) - Built-in generator.
- **GPU (external worker)** - Use `-engine gpu -gpu-bin <path>` to read wallet stream from an external GPU worker process.

## Usage

```console
Usage of ethereum-wallet-generator:
  -n          int    set number of generate times (not number of result wallets) (set number to -1 for Infinite loop ∞, default 10)
  -limit      int    set limit number of result wallets. stop generate when result of vanity wallets reach the limit (set number to 0 for no limit, default 0)
  -db         string set sqlite output file name eg. wallets.db (db file will create in `/db` folder)
  -c          int    set concurrency value (default 1)
  -bit        int    set number of entropy bits [128 for 12 words, 256 for 24 words] (default 128)
  -mode       int    set mode of wallet generator [1: normal mode, 2: only private key mode]
  -strict     bool   strict contains mode, resolve only the addresses that contain all the given letters (required contains to use)
  -contains   string show only result that contained with the given letters (support for multiple characters)
  -prefix     string show only result that prefix was matched with the given letters  (support for single character)
  -suffix     string show only result that suffix was matched with the given letters (support for single character)
  -regex      string show only result that was matched with given regex (eg. ^0x99 or ^0x00)
  -dryrun     bool   generate wallet without a result (used for benchmark speed)
  -compatible bool   logging compatible mode (turn this on to fix logging glitch)
  -decrypt    string decrypt encrypted JSON file (e.g., output/wallets-20260302-101530.encrypted.json)
  -engine     string wallet generation engine [cpu|gpu] (default "cpu")
  -gpu-bin    string path to external GPU worker binary (optional; defaults to `vanity-eth-address`)
  -gpu-args   string arguments for GPU worker binary, space-delimited (used when -engine gpu)
  -output-mode string terminal output mode [safe|full|silent] (default "safe")
  -output-dir string directory for encrypted export file (in-memory mode) (default "output")
  -output-name string encrypted export filename; supports {timestamp} token (in-memory mode)
  -no-export  bool disable encrypted export on exit (in-memory mode)
  -contract-mode string gpu contract search mode [create|create2|create3] (optional)
  -contract-bytecode string path to contract bytecode file (required for create2/create3)
  -contract-address string origin/sender address (required for create2/create3)
  -contract-deployer string deployer address (required for create3)
```

GPU mode notes:
- Currently supports only `-mode 2` (private-key mode).
- Progress status shows: `resolved`, `gpu`, `elapsed`, and `eta(exp/p50/p90)` when prefix/suffix + GPU speed are available.
- External worker must print one wallet per stdout line:
  - JSON: `{"address":"0x...","privateKey":"..."}`
  - or plain text: `<address> <privateKey>`
  - or legacy format: `... Private Key: 0x<64hex> Address: 0x<40hex>`

Example with `l3wi/vanity-eth-address` compatible output:

```console
ethereum-wallet-generator -mode 2 -engine gpu \
  -gpu-bin /usr/local/bin/vanity-eth-address \
  -gpu-args \"--prefix cafe --device 0\"
```

Contract vanity examples (wrapper maps these flags into worker args automatically):

```console
# CREATE (nonce=0)
ethereum-wallet-generator -mode 2 -engine gpu \
  -contract-mode create \
  -prefix 0x7979 -suffix 7979 \
  -gpu-args "--device 0"

# CREATE2
ethereum-wallet-generator -mode 2 -engine gpu \
  -contract-mode create2 \
  -contract-bytecode ./bytecode.txt \
  -contract-address 0x0000000000000000000000000000000000000000 \
  -prefix 0x7979 \
  -gpu-args "--device 0"

# CREATE3
ethereum-wallet-generator -mode 2 -engine gpu \
  -contract-mode create3 \
  -contract-bytecode ./bytecode.txt \
  -contract-address 0x0000000000000000000000000000000000000000 \
  -contract-deployer 0x0000000000000000000000000000000000000001 \
  -prefix 0x7979 \
  -gpu-args "--device 0"
```

### GPU Docker image (bundled worker)

Build:

```console
docker build -f Dockerfile.gpu -t joyai/eth-vanity-address:gpu .
```

Run (NVIDIA runtime required):

```console
docker run --rm --gpus all joyai/eth-vanity-address:gpu \
  -mode 2 -engine gpu -gpu-args "--prefix cafe --device 0"
```

Notes:
- `-gpu-bin` is optional in GPU images; default worker binary is `vanity-eth-address`.
- The bundled worker source is AGPL-3.0 licensed (`l3wi/vanity-eth-address`).
- GPU image currently builds on CUDA 13 runtime.
- `-output-mode safe` hides private keys/mnemonics from terminal output (recommended).

## Benchmark

### Normal Mode

We've dryrun the generator on normal mode with 8 concurrents for 60,000 wallets on MacBook Air M1 2020 Memory 16 GB <br/>
and got speed up to 6,468.58 wallet/sec.

```console
ethereum-wallet-generator -n 60000 -dryrun -c 8 -mode 1
===============ETH Wallet Generator===============

60000 / 60000 | [██████████████████████████████████████████████████████] | 100.00% | 6469 p/s | resolved: 60000

Resolved Speed: 6468.58 w/s
Total Duration: 9.275597416s
Total Wallet Resolved: 60000 w

Copyright (C) 2023 Planxnx <planxthanee@gmail.com>
```

### Only Private Key Mode ⚡️

We've dryrun the generator on only private key mode with 8 concurrents for 1,000,000 wallets on MacBook Air M1 2020 Memory 16 GB <br/>
and got speed up to 111,778 wallet/sec.

```console
ethereum-wallet-generator -n 1000000 -dryrun -c 8 -mode 2
===============ETH Wallet Generator===============

1000000 / 1000000 | [███████████████████████████████████████████████] | 100.00% | 111778 p/s | resolved: 1000000

Resolved Speed: 111778.55 w/s
Total Duration: 8.94626016s
Total Wallet Resolved: 1000000 w

Copyright (C) 2023 Planxnx <planxthanee@gmail.com>
```

## Examples

### **Simple usgae:**

```console
$ ethereum-wallet-generator
# or
$ ethereum-wallet-generator -mode 1
===============ETH Wallet Generator===============

10 / 10 | [██████████████████████████████████████████████████████████████████████████████████████████] | 100.00% | 503 p/s | resovled: 10

Address                                    Seed
------------------------------------------ ------------------------------------------------------------------------------------------
0x239ffa10fcd89b2359a5bd8c27c866cfad8eb75a lecture edge end come west mountain van wing zebra trumpet size wool
0x3addecebd6c63be1730205d249681a179e3c768b need decide earth farm punch crush banana unfold income month bread unhappy
0xc4f55b38e6e586cf974eb005e07482fd40274a26 hundred hundred canvas casual staff candy sign travel sort chat travel space
0xe8df7efc452801dc7c75137136c76006bbc2e6d6 gospel father funny pair catalog today champion maple valid feed loop write
0xdf2809a480e29a883a69beb6dedff095984f09eb poet impulse can undo vital stadium tattoo labor trap now blanket assume
0xabc91fd93be63474c14699a1697533410115824c aisle almost miracle coach practice ostrich thing solution ask kiss idle object
0xc9af163bbac66c1c75f3c5f67f758eed1c6077ba funny shift guilt lucky fringe install sugar forget wagon famous inject evoke
0xcf959644c8ee3c20ac9fbecc85610de067cca890 cupboard analyst remove sausage frame engage visual crowd deny boy firm stick
0xa8904e447afb9e0d9b601669aeca53c9b66fe058 sentence skin april wool huge dad bitter loyal perfect again document boring
0x107a842b628b999827730e4543314c6681c72b93 turkey shove mountain yellow ugly shoot crouch donor topple girl polar pelican


Resolved Speed: 502.78 w/s
Total Duration: 24.832485ms
Total Wallet Resolved: 10 w

```

### **🎨️⚡ Generate until got expected number of vanity addresses and Speeding up with concurrency:**

```console
$ ethereum-wallet-generator -n -1 -limit 5 -contains 0x000,0x777 -c 8
===============ETH Wallet Generator===============

12435 | [██████████████████████████████████████████████████████████████████████████████████████████] | 100.00% | 5073 p/s | resovled: 5

Address                                    Seed
------------------------------------------ ------------------------------------------------------------------------------------------
0x0002e20a2528c4fe90af8dd10a38231d1d937531 jump skull butter topic bronze member feed wait flee oven deer rabbit
0x000ff528ae048f2cb71cea8cdeb0198ad45ff09f rotate change tooth design price milk derive olympic small sudden payment hover
0x000b98901463db593613e749d7a4803f24e3a7bb fish zone shift river sort visit begin hunt august trouble fatal easy
0x7772eef4d1f660d8cd0b89f4d6cdf90175b63b3a review today coil purity mouse lucky trip collect mail right weekend remove
0x77746fe800078d956b3176c367be88cdc39cd878 fiscal east exhibit arena expand indicate fury document vacuum mansion aisle garbage

Resolved Speed: 1.48 w/s
Total Duration: 2.4512123s
Total Wallet Resolved: 5 w
```

### **⚠⚡️ ️Extream speeding up with concurrency `Only Private Key mode` for generate vanity addresses:**

```console
$ ethereum-wallet-generator -n -1 -limit 5 -contains 0x00000,0x11111 -c 8 -mode 2
===============ETH Wallet Generator===============

252237 | [██████████████████████████████████████████████████████████████████████████████████████████] | ?% | 102903 p/s | resolved: 5
Address                                    Private Key
------------------------------------------ ------------------------------------------------------------------------------------------
0x0000005927adc84c599084c48f50525617a76cf6 aaf26b1e0d137813bd221e59b3e072a5ad8b58e36c5d30ae19e8d0e5e19f287e
0x111111285bf095c0fa68bc170f9c23a43af9ead0 2826c6596897074a3545fce2904e961a69291efce09585c81417587603a6ca55
0x11111235eebd9d28420aaae50ac243e0be980618 7b1993e90956d929036918ddad9f5a4de3cd946bee569e01a0226be57de94862
0x000008a6481c0396d318c04417b4cdbc053aef1f d306a6f9808f073e455e11bddb206403d175874e327835dc06d6e81f325231b0
0x11111022d16444795ba8c5ee348f2f24650888af ec4263f4879837afa1c380b8252ffd0ddbe468b49a18254b47b4b1f22ea900da

Resolved Speed: 1.48 w/s
Total Duration: 2.4512123s
Total Wallet Resolved: 5 w
```

### **24 word seed prhase and filter vanity addresses with contains and strict options:**

```console
$ ethereum-wallet-generator -n 50000 -limit 2 -contains 0x00,777,22 -strict -bit 256
===============ETH Wallet Generator===============

31099 / 50000 | [██████████████████████████████████████████████████████████████████████████████████████████] | 100.00% | 2277 p/s | resovled: 2

Address                                    Seed
------------------------------------------ ------------------------------------------------------------------------------------------
0x00325b7844a4c8612108f407c0ad722da3294777 delay pilot wall radio next uniform margin copper plunge kidney coil runway baby major token method arena brave expand route job raise budget buffalo
0x00ca8750521c2270e7776becd05d4a7d1b2ffdcb insect fashion original page stamp grow mean cinnamon embody favorite near useless relief crouch ranch nerve card captain situate truly cousin renew birth credit


Resolved Speed: 0.14 w/s
Total Duration: 13.857967333s
Total Wallet Resolved: 2 w
```

### **📚 Storing to embeded databse(SQLite3) to easily management:**

```console
$ ethereum-wallet-generator -n 50000 -c 12 -db 0x77.db -prefix 0x77
===============ETH Wallet Generator===============

50000 / 50000 | [██████████████████████████████████████████████████████████████████████████████████████████] | 100.00% | 5384 p/s | resovled: 178

Resolved Speed: 16.04 w/s
Total Duration: 11.09416725s
Total Wallet Resolved: 178 w
```

![image](https://user-images.githubusercontent.com/37617738/227806706-02a8a7fa-7d2b-43ca-b89b-c21cc51835ff.png)

### **🐳 Use Docker (recommend using concurrency for speed up):**

```console
$ docker run --rm -v $PWD:/db planxthanee/ethereum-wallet-generator -n 50000 -db wallet.db -c 8
===============ETH Wallet Generator===============

  100% |██████████████████████████████████████| (50000/50000, 4651 w/s) [10s:95ms]

Resolved Speed: 4563.50 w/s
Total Duration: 10.9545307s
Total Wallet Resolved: 50000 w

```

### **🔐 Encrypted Export & Decrypt:**

When using in-memory mode (without `-db` flag), wallets are automatically exported to an encrypted JSON file after generation:

```console
$ ethereum-wallet-generator -n 10
===============ETH Wallet Generator===============

10 / 10 | [██████████████████████████████████████████████████████████████████████████████████████████] | 100.00% | 503 p/s | resovled: 10

Enter password to encrypt wallets: [password hidden]
Successfully exported 10 wallets to output/wallets-20260302-101530.encrypted.json
```

Decrypt the encrypted file:

```console
$ ethereum-wallet-generator -decrypt output/wallets-20260302-101530.encrypted.json
Enter password to decrypt: [password hidden]
Successfully decrypted 10 wallets to output/wallets.json
```

**Security Features:**
- AES-256-GCM encryption with PBKDF2 key derivation (100,000 iterations)
- Password input is hidden (not visible in terminal)
- All encrypted/decrypted files are saved to `output/` directory
- Memory is cleared after encryption
- Output directory is automatically added to `.gitignore`
- **Addresses are visible in plain text** for easy identification of which wallets you have
- Only `privateKey` and `mnemonic` fields are encrypted (indicated by `hasPrivateKey` and `hasMnemonic` flags)

## Thanks to

- [conseweb/coinutil](https://github.com/conseweb/coinutil) - for BIP39 implementation in Go
- [tyler-smith/go-bip39](https://github.com/tyler-smith/go-bip39) - for BIP39 implementation in Go

## Contributing ![eth](https://user-images.githubusercontent.com/37617738/120125730-1d1bd680-c1e4-11eb-83ad-45664245cae9.png)

Pull requests are welcome!

For contributions please create a new branch and submit a pull request for review.

## License

This repo is released under the [WTFPL](http://www.wtfpl.net/) – Do What the Fuck You Want to Public License.
