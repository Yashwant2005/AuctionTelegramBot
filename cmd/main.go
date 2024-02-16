package main

import (
	"AuctionBot/auction"
	"AuctionBot/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

func IsAdmin(username string, adminUserLists []string) bool {
	for _, admin := range adminUserLists {
		if admin == username {
			return true
		}
	}
	return false
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Error reading config file")
	}
}

func main() {
	token := viper.GetString("bot_token")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	auctionBids := make(chan tgbotapi.Update)
	auctioneerMessages := make(chan tgbotapi.Chattable)

	go func(messages chan tgbotapi.Chattable) {
		for {
			select {
			case message := <-messages:
				bot.Send(message)
			}
		}
	}(auctioneerMessages)

	var adminUserLists = viper.GetStringSlice("admin_usernames")

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			if update.Message.IsCommand() {
				chatID := update.Message.Chat.ID
				switch update.Message.Command() {
				case "start":
					sender := update.Message.From.UserName
					if !IsAdmin(sender, adminUserLists) {
						message := tgbotapi.NewMessage(chatID, messages.NOT_ADMIN_MESSAGE)
						message.ReplyToMessageID = update.Message.MessageID
						bot.Send(message)
						continue
					}
					if auction.GetActiveAuction() != nil {
						message := tgbotapi.NewMessage(chatID, messages.ACTIVE_AUCTION_EXISTS_MESSAGE)
						message.ReplyToMessageID = update.Message.MessageID
						bot.Send(message)
						continue
					}
					startConfig, err := auction.ParseStartAuctionCommand(update.Message.Text)
					if err != nil {
						message := tgbotapi.NewMessage(chatID, messages.INVALID_START_MESSAGE)
						message.ReplyToMessageID = update.Message.MessageID
						bot.Send(message)
						continue
					}

					a := auction.NewFirstPriceAuction(startConfig.Name, startConfig.StartPrice, startConfig.MinStep)
					go auction.Auctioneer(a, chatID, auctioneerMessages, auctionBids)

				case "bid":
					if auction.GetActiveAuction() != nil {
						auctionBids <- update
					} else {
						message := tgbotapi.NewMessage(chatID, messages.NO_ACTIVE_AUCTION_MESSAGE)
						message.ReplyToMessageID = update.Message.MessageID
						bot.Send(message)
					}

				case "help":
					message := tgbotapi.NewMessage(chatID, messages.HELP_MESSAGE)
					message.ReplyToMessageID = update.Message.MessageID
					bot.Send(message)
				}
			}
		}
	}
}
