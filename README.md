# BIP-39 Splitter

`bip39-splitter` is a command line tool **to split (and reconstruct) a BIP-39 mnemonic (seed phrase)**. If done
correctly, this improves the **availability** and **security** of your mnemonic by adding **fault tolerance** and
**decentralization** in the storage of your mnemonic.

`bip39-splitter` is:
* **Cross-platform:** it will compile to all OSes/architectures supported by Golang
* **Self-contained:** a single binary with no dependency to make sure it will always work
* **Small:** <400 lines of code excluding the Shamir Secret Sharing implementation (~600 lines including it)

All languages
[officially supported by BIP-39](https://github.com/bitcoin/bips/blob/master/bip-0039/bip-0039-wordlists.md) (as of
March 27, 2025) are supported:
* English
* Japanese
* Korean
* Spanish
* Chinese (Simplified)
* Chinese (Traditional)
* French
* Italian
* Czech
* Portuguese

Dictionaries were copied from the [official BIPs repository](https://github.com/bitcoin/bips/tree/master/bip-0039).

## Usage
> [!IMPORTANT]  
> For best security, this program should be compiled and executed on an
> **[air-gapped](https://en.wikipedia.org/wiki/Air_gap_(networking)) computer**.

First, build the executable:
```bash
go build -o bip39splitter cmd/bip39splitter/main.go
```

Then, split your mnemonic:
```text
$ ./bip39splitter split -l english -p 5 -t 3
Enter your BIP-39 mnemonic (seed phrase): icon thunder tube include demise charge valley odor asthma volcano cost invest again hidden fiction goose bulk twin retire chair tooth genius flight erase

Parts:
==================================================================================================
caff12146c4ba97dd69ed64752d44a16903a01d8591d5d1cd47a653c89a47f8d3044259c0f993122786178f58f7d434cd9
9c5102654e850ea4c847ae14fc0ef3f85de3553d277affb369b17ca7152398b858b0fc276d3669908b4f0e07cd98be581b
a036298481395175e340563af8ec29a15327c71a0a05967dbb927a0b629fc1369b75a3686302e44487ffaeb632dd7eec6b
6bef4382cf502d0cc672227243ea71ef83661438736b092c1b118d9b6e1e6ac20abade05ecb22cf0e4a5dab7704fe3ede5
c6f475900b1c147918d5f1a4660a74d6f9bfa9aa88aeaafdb674354741a7cf89411d0e596c0723394d1a96dc8d47da8bcb
==================================================================================================
```
In this example, language is set to English (`-l english`), number of parts is set to 5 (`-p 5`) and number of required
parts to be able to reconstruct the mnemonic is set to 3 (`-t 3`).

> [!TIP]
> It is expected that the same mnemonic will give different parts each time as Shamir Secret Sharing is based on
> randomly generated numbers.

Later, you can reconstruct the mnemonic by gathering any 3 of the 5 parts (order doesn't matter):
```text
$ ./bip39splitter combine -l english
Enter part 1 and press ENTER (leave empty to finish): 9c5102654e850ea4c847ae14fc0ef3f85de3553d277affb369b17ca7152398b858b0fc276d3669908b4f0e07cd98be581b
Enter part 2 and press ENTER (leave empty to finish): 6bef4382cf502d0cc672227243ea71ef83661438736b092c1b118d9b6e1e6ac20abade05ecb22cf0e4a5dab7704fe3ede5
Enter part 3 and press ENTER (leave empty to finish): caff12146c4ba97dd69ed64752d44a16903a01d8591d5d1cd47a653c89a47f8d3044259c0f993122786178f58f7d434cd9
Enter part 4 and press ENTER (leave empty to finish):

Reconstructed mnemonic (seed phrase):
icon thunder tube include demise charge valley odor asthma volcano cost invest again hidden fiction goose bulk twin retire chair tooth genius flight erase
```

## Motivation
Splitting your mnemonic into multiple parts (stored in different places) brings two major availability and security
advantages over storing it in one single place:
* **Fault tolerance:** depending on the splitting parameters you choose, you can afford to lose one or more parts of
  your mnemonic without compromising your ability to recover the mnemonic.
* **Decentralization:** an attacker needs to gather several parts of your mnemonic (requiring access to different
  places) to steal it, making it much more difficult than if it is stored wholly in a single place.

In other words, splitting your mnemonic makes it more secure against attacks and _your own_ mistakes.

## Algorithm
There are some steps involved before and after splitting the secret with Shamir Secret Sharing. Here is the complete
algorithm implemented in `bip39-splitter`:
1. Turn each word of the mnemonic into its corresponding index in the
   [BIP-39 dictionary](https://github.com/bitcoin/bips/tree/master/bip-0039) of the chosen language (e.g. word `ability`
   of the English dictionary becomes `2`, as it's the second word of the list).
2. Convert word indices to an array of bytes (keep the order!). Indices are **2 bytes long** and encoded in **Big
   Endian** (e.g. index `2` from the word `ability` is represented as `0x00 0x02`).
3. Apply Shamir Secret Sharing to the resulting array of bytes. The SSS implementation used is [the one present in
   Hashicorp Vault v1.19.0](https://github.com/hashicorp/vault/blob/v1.19.0/shamir/shamir.go).
4. Encode each part in its hexadecimal string representation (see Golang's
   [`hex.EncodeToString()` function](https://pkg.go.dev/encoding/hex#EncodeToString) implementation for details).

### Shamir Secret Sharing (SSS) implementation
As per the concept of _Don't roll your own cryptography_, the SSS implementation (which is the core of this tool) was
copied as-is from [Hashicorp Vault](https://www.vaultproject.io/), an enterprise-grade secrets management software. The
version used in `bip39-splitter` is extracted from
[Vault v1.19.0](https://github.com/hashicorp/vault/blob/v1.19.0/shamir/shamir.go).

The code was copied (instead of imported as a Go module) in order to avoid any build dependency, as it is critical to be
able to (re)build this software many years from now.

> [!IMPORTANT]
> As all implementations of Shamir Secret Sharing tend to be slightly different (plus the fact that we add extra steps
> besides SSS), this software is not compatible with other SSS implementations such as the `sss` tool on Linux. This
> means YOU NEED TO KEEP A COPY OF THIS PROGRAM TO RECONSTRUCT YOUR MNEMONIC IN THE FUTURE.

## Licensing
This software is distributed under two licenses:
* The Shamir Secret Sharing implementation, being taken from Hashicorp Vault, remains under their
  [Mozilla Public License version 2.0](pkg/shamir/LICENSE).
* The rest of the program is under [MIT License](LICENSE), as permitted by the Mozilla Public License version 2.0.
