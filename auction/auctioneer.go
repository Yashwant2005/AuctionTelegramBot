package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Auctioneer struct {
	// Auction
	auction Auction
	// Chat ID
	chatID int64
	// Send channel
	send chan tgbotapi.Chattable
	// Receive channel
	bidsChannel chan Bid
	// Stop channel
	stopChannel chan string
}

func NewAuctioneer(config StartAuctionConfig, chatID int64, send chan tgbotapi.Chattable) Auctioneer {
	var auction Auction
	switch config.Type {
	case "reverse_auction":
		auction = NewReverseAuction(config.Name, config.StartPrice, config.MinStep)
	case "special_auction":
		auction = NewSpecialAuction(config.Name, config.StartPrice, config.MinStep)
	}
	stopChannel := make(chan string)
	bids := make(chan Bid)

	return Auctioneer{
		auction:     auction,
		chatID:      chatID,
		send:        send,
		bidsChannel: bids,
		stopChannel: stopChannel,
	}
}

var activeAuction Auction

func GetActiveAuction() Auction {
	return activeAuction
}

func (a *Auctioneer) Run(receive tgbotapi.UpdatesChannel) {
	activeAuction = a.auction
	a.auction.Start()

	startingMessage := tgbotapi.NewMessage(a.chatID, fmt.Sprintf(messages.START_AUCTION_MESSAGE, a.auction.Name(), a.auction.StartPrice(), a.auction.MinStep()))
	a.send <- startingMessage

	go a.auction.Auctioneer()(a)

	for {
		select {
		case <-a.stopChannel:
			activeAuction = nil
			a.auction.End()

			winner := a.auction.Winner()
			winnerPrice := a.auction.WinnerPrice()

			messageText := fmt.Sprintf(messages.END_AUCTION_MESSAGE, a.auction.Name(), winner, winnerPrice)
			message := tgbotapi.NewMessage(a.chatID, messageText)
			a.send <- message

			file := tgbotapi.FilePath("./media/leonardo-dicaprio-sold-gif.mp4")
			gif := tgbotapi.NewAnimation(a.chatID, file)
			a.send <- gif

			a.auction.WriteLog()
			return

		case update := <-receive:
			if update.Message != nil {
				bid, err := activeAuction.ParseBid(update)
				if err != nil {
					message := tgbotapi.NewMessage(update.Message.Chat.ID, messages.INVALID_BID_MESSAGE)
					message.ReplyToMessageID = update.Message.MessageID
					a.send <- message
					continue
				}
				if bid.AuctionName != a.auction.Name() {
					message := tgbotapi.NewMessage(update.Message.Chat.ID, messages.NO_AUCTION_MESSAGE)
					message.ReplyToMessageID = update.Message.MessageID
					a.send <- message
					continue
				}
				if update.Message.Chat.ID != a.chatID {
					message := tgbotapi.NewMessage(update.Message.Chat.ID, messages.INVALID_CHAT_ID_MESSAGE)
					message.ReplyToMessageID = update.Message.MessageID
					a.send <- message
					continue
				}
				a.bidsChannel <- bid
			}
		}
	}
}
