# Telegram Auction Bot

This is an auction bot for Telegram. It is written in Golang and uses the [telegram-bot-api](https://core.telegram.org/bots/api).
This bot is mainly developed for course Algorithmic Game Theory at the Sharif University of Technology by [Dr. Masoud Seddighin](https://sites.google.com/view/masoudseddighin).

# Setup

## Pre-requisites

1. Make sure you have Go installed on your system. You can download it from [here](https://golang.org/dl/).
2. Make sure you have a Telegram bot token. You can get one by talking to [BotFather](https://t.me/botfather).

## Steps

1. Create a `config.yaml` file in the `config` directory. You can use the `config.example.yaml` file as a template. Fill in the `bot_token` field with your bot token. Also fill the `admin_usernames` field with the usernames of the users you want to be able to start an auction. 
2. Create a `log` directory in the root of the project.
3. Resolve the dependencies by running `go mod tidy`.
4. Run the bot using `go run cmd/main.go`.

# Room for Improvement
[ ]: Add ascending auction

# Contact
If you have any questions or suggestions, feel free to contact me at [absoltani02@gmail.com](mailto:absoltani02@gmail.com).