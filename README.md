# My Savings &mdash; Telegram Bot

[![License](https://img.shields.io/badge/License-Unlicense-000000.svg)](https://unlicense.org/)

<p align="center">
    <img width="256" height="256" src="./assets/peppa.png">
</p>

## Description

🤖 This bot needs for simply storing your savings in a different currencies.

Originally, you may find it at https://t.me/MySavingsV0Bot, but it not always working, cause resources are limited.&nbsp;😭

## Usage

You may run your own copy of this bot with Docker:

```bash
docker pull lordotu/my-savings-telegram-bot

docker run -dti \
  -e FIXERIO_API_KEY=<https://fixer.io API Key> \
  -e TELEGRAM_API_TOKEN=<Telegram API Token> \
  --name my-savings-telegram-bot \
  lordotu/my-savings-telegram-bot
```

But before, you should register your bot via **BotFather** https://t.me/BotFather and get API key for **Fixer** https://fixer.io/signup/free

And don't forget about **MongoDB** which is used for storing data.

## Configuring

| Option                       |                                 Default |
|------------------------------|----------------------------------------:|
| DEBUG                        |                                   false |
| DEBUG_TELEGRAM_API           |                                   false |
| EXCHANGE_RATES_GETTING_TIMER |                                    3600 |
| FIXERIO_API_KEY              |                                         |
| FIXERIO_BASE_CURRENCY        |                                     EUR |
| FIXERIO_SECURE               |                                   false |
| FIXERIO_SYMBOLS              | USD,EUR,JPY,GBP,AUD,CAD,CHF,CNY,NZD,RUB |
| MONGO_DB_URI                 |               mongodb://localhost:27017 |
| MONGO_DB_DATABASE            |                 my-savings-telegram-bot |
| RATES_PROVIDER               |                                 fixerio |
| TELEGRAM_API_TOKEN           |                                         |
| TELEGRAM_API_UPDATES_TIMEOUT |                                      60 |
