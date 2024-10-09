[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/fomo/app?startapp=ref_JDLG5)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/EtherDrop/main/demo/demo.png)

# 🔥🔥 Ether Drop Bot With Support Cross Platform 🔥🔥

### Tested on Windows and Docker Alpine Os with a 4-core CPU using 5 threads.

**Go Version Tested 1.23.1**

# 🔥🔥 UPDATED 🔥🔥

**Using Session Storage Of Telegram Web Because Query Data Can Expired**

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) Must >= 1.23.
- #### Rename config.yml.example to config.yml.
- #### Place your browser session local storage .json file in the sessions folder.
- #### If you don’t have a local storage session, you can obtain it from [Telegram Web Tools](https://github.com/ehhramaaa/telegram-web-tools)
- #### If you want to use a custom browser, set the browser path in config.yml.
- #### Rename proxy.txt.example to proxy.txt and place your query data.
- #### It is recommended to use an IP info token to improve request efficiency when checking IPs.
- #### Auto Get Ref Code Will Generate File ref_code.json in main folder.
- #### Auto Register With Ref Code Required Input Your Ref Code.

## Features

|           Feature           | Supported |
| :-------------------------: | :-------: |
|        Get Ref Code         |    ✅     |
| Auto Register With Ref Code |    ✅     |
|   Auto Claim Daily Bonus    |    ✅     |
|       Auto Claim Task       |    ✅     |
|    Auto Claim Ref Reward    |    ✅     |
|      Proxy SOCKS5/HTTP      |    ✅     |
|     Use Session Storage     |    ✅     |
|       Multithreading        |    ✅     |
|      Random User Agent      |    ✅     |

## [Settings](https://github.com/ehhramaaa/EtherDrop/blob/main/configs/config.yml.example)

|         Settings         |                       Description                       |
| :----------------------: | :-----------------------------------------------------: |
|      **USE_PROXY**       |                For Using Proxy From File                |
|     **IPINFO_TOKEN**     | For Increase Check Ip Efficiency Put Your Token If Have |
|     **RANDOM_SLEEP**     |       Delay before the next lap (e.g. 1800, 3600)       |
| **RANDOM_REQUEST_DELAY** |       Delay before the next Request (e.g. 5, 10)        |
|      **MAX_THREAD**      |             Max Thread Worker Run Parallel              |

## Installation

```shell
git clone https://github.com/ehhramaaa/EtherDrop.git
cd EtherDrop
go mod tidy
```

## Usage

```shell
go run .
```

Or

```shell
go run main.go
```

## Or you can do build application by typing:

Windows:

```shell
go build -o EtherDrop.exe
```

Linux:

```shell
go build -o EtherDrop
chmod +x EtherDrop
./EtherDrop
```

**If You Want Auto Select Choice In Terminal or Apps**

For Option 1

```shell
go run . -action 1
```

**If Already Build To Apps**

```shell
./EtherDrop -action 1
```

For Option 2

```shell
go run . -action 2
```

**If Already Build To Apps**

```shell
./EtherDrop -action 2
```

For Option 3

```shell
go run . -action 3
```

**If Already Build To Apps**

```shell
./EtherDrop -action 3
```

## Mobile Guide

### Prerequisites 📚

> Download Termux, ISH Shell Or Another Terminal With Linux Base

> **Make Sure Golang Version >= 1.23**

### Installation

```shell
apt update && apt upgrade -y
apt instal golang chromium -y
go version -v
git clone https://github.com/ehhramaaa/EtherDrop.git
cd EtherDrop
go mod tidy
```

> For Usage and build application it's same like linux

> # Feel Free To Ask Any Question Join : [Skibidi Sigma Code Group](https://t.me/skibidi_sigma_code_chat)
