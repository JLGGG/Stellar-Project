## Stellar Test Network Telegram Bot

[![MIT licensed](https://img.shields.io/github/license/JLGGG/Stellar-Project?label=License)](https://github.com/JLGGG/Stellar-Project/blob/main/LICENSE)

This is a stellar bot running on telegram. This bot provides a function that makes it easy to use the stellar test network. Currently, there are a total of 7 commands supported by bot.
To use this bot, you need to create a personal bot using Telegram's [BotFather](https://core.telegram.org/bots).   
You can create a token with BotFather's /newbot command. Keep your token secure and store it safely, it can be used by anyone to control your bot.
You must use this generated token.   

## How to use the program
- Package management with Go modules. Refer to [go.mod](https://github.com/JLGGG/Stellar-Project/blob/main/go.mod) for the package information used.
- After downloading the program, go to the folder where main.go is located.
- `go run main.go token`: In the Token field, enter your personal bot token created using BotFather.

## Supported commands
- `/`: You can see a list of commands supported by Stellar bot.
- `/hello`: Output "hello world"
- `/make_account`: Create an account that can be used on the stellar test network.
- `/show_account`: This command shows the list of accounts you currently have.
- `/send_money`: This command allows the user to transfer the money to the receiving account of the stellar test network.
- `/show_favorite`: This is a command that shows the list of receiving accounts that the user likes to use.
- `/save_favorite`: This command helps to save a new receiving account.
- `/delete_favorite`: This is a command to help you delete an unwanted account from the list of receiving accounts.

## Folder information
- Practice: Pratice using stellar test network API.
- stellar: Files used in main.go

## Reference
- [`Stellar`](https://developers.stellar.org/docs)
- [`Telegram`](https://core.telegram.org/bots)
- [`Telebot`](https://github.com/tucnak/telebot)
