# `cryptctl` CLI Tool for managing [EncryptedSecrets](https://github.com/opensecrecy/encrypted-secrets)


[![Build Status](https://github.com/opensecrecy/cryptctl/workflows/CI/badge.svg)](https://github.com/opensecrecy/cryptctl/actions?query=workflow%3ACI+branch%3Amain)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

## Note: `cryptctl` is currently a work in progress and is in the alpha stage. Please use it with caution in production environments.
## What is `cryptctl`?

**cryptctl** is a simple command-line interface (CLI) tool designed to facilitate the management of EncryptedSecrets.
With Cryptctl, you can easily update encrypted secrets within your Kubernetes cluster, ensuring the secure handling of sensitive information.

### Features
- **Effortless Encryption:** Cryptctl simplifies the process of encrypting secrets by providing a straightforward command-line interface. It handles the encryption and decryption operations seamlessly, making it easy to work with encrypted secrets in your Kubernetes environments.

- **Simplified Management:** Since, the secrets are encrypted, they can be easily stored in a repository. Once, the `EncryptedSecret` object is applied, `encrypted-secrets` controller takes care of decrypting the provided secrets and creates a k8s secret with decrpted values. Which can be access by the pod as required.


Here's a **`cryptctl`** demo:

- `cryptctl edit <filename>`
![cryptctl edit demo GIF](img/cryptctl-edit-demo.gif)

- `cryptctl create -f <filename> -p <provider>`
![cryptctl edit demo GIF](img/cryptctl-create-demo.gif)


## Installation

- Linux: Download the latest release from [here](https://github.com/opensecrecy/cryptctl/releases)

- Mac (you can use [Homebrew](https://brew.sh/) to install `cryptctl` on Mac):
```bash
brew tap opensecrecy/cryptctl
brew install opensecrecy/cryptctl/cryptctl
```

- Windows: Download the latest release from [here](https://github.com/opensecrecy/cryptctl/releases)

# Usage
## Pre-requisites
- Make sure you have EDITOR environment variable set to your preferred editor. (e.g. `export EDITOR=vim`)
- Make sure you have a kubernetes cluster running and `kubectl` is configured to access the cluster in case of `k8s` provider and if you are using `aws-kms` provider make sure you have aws credentials configured.

## Commands
**1. init:** Initialize the encryption provider. This command will create an encryption certificate in the namespace provided if you are using `k8s` provider and will create a kms key if you are using `aws-kms` provider.

**2. create:** This command will create an `EncryptedSecret` object manifest with the provided inputs and provider.

**3. edit:** This command will open the provided `EncryptedSecret` object manifest in the editor provided by the EDITOR environment variable. Once the file is saved, it will create the `EncryptedSecret` object manifest which you can apply.

## Comparison with some other popular alternatives

| Tool     | Namespaced Encryption Keys | Multiple Providers Support | Auto Key Renewal | Client Side Decryption Support |
|:--------:|:------------------------:|:--------------------------:|:----------------:|:------------------------------:|
| cryptctl | ✅                        | ✅                          | ❌                | ✅                              |
| kubeseal | ❌                        | ❌                          | ✅                | ❌                              |
