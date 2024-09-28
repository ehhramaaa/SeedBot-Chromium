[![Static Badge](https://img.shields.io/badge/Telegram-Bot%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](http://t.me/seed_coin_bot/app?startapp=5024522783)
[![Static Badge](https://img.shields.io/badge/Telegram-Channel%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_code)
[![Static Badge](https://img.shields.io/badge/Telegram-Chat%20Link-Link?style=for-the-badge&logo=Telegram&logoColor=white&logoSize=auto&color=blue)](https://t.me/skibidi_sigma_chat)

![demo](https://raw.githubusercontent.com/ehhramaaa/SeedBot-Chromium/main/demo/demo.png)

# 🔥🔥 Seed Bot Auto Get Query Data, Auto Claim And Auto Completing Task 🔥🔥

## Recommendation before use

### Tested in windows and linux
#### Go Version >= 1.23
#### The thread usage depends on your system specifications, as this uses a browser and may feel heavy on lower-spec systems.
#### Rename config.yml.example to config.yml.
#### Place your browser session local storage .json file in the sessions folder.
#### If you want to use a custom browser, set the browser path in config.yml.
#### It is recommended to use an IP info token to improve request efficiency when checking IPs.

## Features

|             Feature             | Supported |
| :-----------------------------: | :-------: |
|          Auto Farming           |    ✅     |
|         Multithreading          |    ✅     |
|         Auto Bird Hunt          |    ✅     |
|         Auto Feed Bird          |    ✅     |
|         Auto Hatch Egg          |    ✅     |
|         Use Query Data          |    ✅     |
|         Auto Catch Worm         |    ✅     |
|        Auto Upgrade Tree        |    ✅     |
|       Auto Play Spin Egg        |    ✅     |
|        Random User Agent        |    ✅     |
|      Auto Claim First Egg       |    ✅     |
|      Auto Check Inventory       |    ✅     |
|      Auto Upgrade Storage       |    ✅     |
|     Auto Claim Login Bonus      |    ✅     |
|     Auto Upgrade Holy Water     |    ✅     |
|    Auto Claim Streak Reward     |    ✅     |
|    Auto Completing Main Task    |    ✅     |
|  Auto Increase Bird Happiness   |    ✅     |
| Auto Completing Holy Water Task |    ✅     |
|        Auto Fusion Piece        |    ⏳     |
|       Auto Connect Wallet       |    ⏳     |
|       Proxy Socks5 / HTTP       |    ✅     |

## [Settings](https://github.com/ehhramaaa/SeedBot-Chromium/blob/main/config.yml)

|               Settings                |                             Description                             |
| :-----------------------------------: | :-----------------------------------------------------------------: |
|          **AUTO_FEED_BIRD**           |                  Auto Feed Bird If Worm Available                   |
|          **AUTO_BIRD_HUNT**           |                 Auto Bird Hunt If Energy Available                  |
|          **AUTO_HATCH_EGG**           |              Auto Hatch Egg If Available In Inventory               |
|        **AUTO_PLAY_SPIN_EGG**         |               Auto Play Spin Egg If Ticket Available                |
|        **AUTO_UPGRADE.SPEED**         |                 Auto Upgrade Tree If Balance Enough                 |
|       **AUTO_UPGRADE.STORAGE**        |               Auto Upgrade Storage If Balance Enough                |
|     **CLAIM_FARMING_SEED_AFTER**      | Delay Before Claim Farming Seed (e.g. MIN:3600, MAX:7200) In Second |
|      **AUTO_UPGRADE.HOLY_WATER**      |        Auto Upgrade Holy Water If Task Holy Water Completed         |
|   **AUTO_UPGRADE.MAX_LEVEL.SPEED**    |                  Max Upgrade Level Amount For Tree                  |
|  **AUTO_UPGRADE.MAX_LEVEL.STORAGE**   |                Max Upgrade Level Amount For Storage                 |
| **AUTO_UPGRADE.MAX_LEVEL.HOLY_WATER** |               Max Upgrade Level Amount For Holy Water               |
|           **RANDOM_SLEEP**            |    Delay Before The Next Lap (e.g. MIN:3600, MAX:7200) In Second    |

## Prerequisites 📚

Before you begin, make sure you have the following installed:

- [Golang](https://go.dev/doc/install) **Go Version Tested 1.23.1**

## Installation

You can download the [**repository**](https://github.com/ehhramaaa/SeedBot-Chromium.git) by cloning it to your system and installing the necessary dependencies:

```shell
git clone https://github.com/ehhramaaa/SeedBot-Chromium.git
cd SeedBot-Chromium
go run .
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
go build -o SeedBot.exe
```

Linux:

```shell
go build -o SeedBot
```
