[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/fomo/app?startapp=ref_JDLG5)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/EtherDrop/main/demo/demo.png)

# 🔥🔥 Ether Drop Bot With Support Cross Platform 🔥🔥

### Tested on Windows and Docker Alpine Os with a 4-core CPU using 5 threads.

**Go Version Tested 1.23.1**

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) Must >= 1.23.
- #### Rename config.yml.example to config.yml.
- #### Rename query.txt.example to query.txt and place your query data.
- #### Rename proxy.txt.example to proxy.txt and place your query data.
- #### If you don’t have a query data, you can obtain it from [Telegram Web Tools](https://github.com/ehhramaaa/telegram-web-tools)
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
|       Use Query Data        |    ✅     |
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

## Mobile Guide

### Prerequisites 📚

> Download Termux, ISH Shell Or Another Terminal With Linux Base
> Make Sure Golang Version >= 1.23

### Installation

```shell
apt update && apt upgrade -y
apt instal golang -y
go version -v
git clone https://github.com/ehhramaaa/EtherDrop.git
cd EtherDrop
go mod tidy
```

### Usage

```shell
cd EtherDrop
go run .
```

## Or you can do build application by typing:

```shell
cd EtherDrop
go build -o EtherDrop
chmod +x EtherDrop
./EtherDrop
```

# Feel Free To Ask Any Question Join : [Skibidi Sigma Code Group](https://t.me/skibidi_sigma_chat)
