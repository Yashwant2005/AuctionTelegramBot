package auction

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"time"
)

var bidPattern = regexp.MustCompile(`^/bid (\w+) (\d+(\.\d+)?)$`)

func parseBidMessage(update tgbotapi.Update) (Bid, error) {
	text := update.Message.Text

	if !bidPattern.MatchString(text) {
		return Bid{}, fmt.Errorf("bid command should be in format %s", bidPattern.String())
	}

	matches := bidPattern.FindStringSubmatch(text)
	auctionName := matches[1]
	amount, _ := strconv.ParseFloat(matches[2], 32)
	return Bid{
		AuctionName: auctionName,
		Bidder:      update.Message.From.UserName,
		Amount:      amount,
	}, nil
}

var activeAuction Auction

func GetActiveAuction() Auction {
	return activeAuction
}

func Auctioneer(auction Auction, chatID int64, send chan tgbotapi.MessageConfig, receive tgbotapi.UpdatesChannel) {
	startingMessageText := auction.Start()
	activeAuction = auction

	startingMessage := tgbotapi.NewMessage(chatID, startingMessageText)
	send <- startingMessage
	duration := 40 * time.Second

	notice := time.NewTimer(duration - 1*time.Second)
	timer := time.NewTimer(duration)

	for {
		select {
		case <-notice.C:
			noticeMessageText := "Auction is going to end in 10 seconds. Bid sooner if you want to win!"
			noticeMessage := tgbotapi.NewMessage(chatID, noticeMessageText)
			send <- noticeMessage

		case <-timer.C:
			endingMessageText := auction.End()
			activeAuction = nil

			endingMessage := tgbotapi.NewMessage(chatID, endingMessageText)
			send <- endingMessage

			winner, _ := auction.Winner()
			winnerPrice, _ := auction.WinnerPrice()
			messageText := fmt.Sprintf("Winner of the auction is %s with price %f", winner, winnerPrice)
			message := tgbotapi.NewMessage(chatID, messageText)
			send <- message

			auction.WriteLog()
			return

		case update := <-receive:
			if update.Message != nil {
				bid, err := parseBidMessage(update)
				if err != nil {
					message := tgbotapi.NewMessage(chatID, err.Error())
					message.ReplyToMessageID = update.Message.MessageID
					send <- message
					continue
				}
				if bid.AuctionName != auction.Name() {
					message := tgbotapi.NewMessage(chatID, "This bid is for another auction")
					message.ReplyToMessageID = update.Message.MessageID
					send <- message
					continue
				}
				result, err := auction.Bid(bid.Bidder, bid.Amount)
				var messageText string
				if err != nil {
					messageText = err.Error()
				} else {
					messageText = result
					notice.Reset(duration - 5*time.Second)
					timer.Reset(duration)
				}
				message := tgbotapi.NewMessage(chatID, messageText)
				message.ReplyToMessageID = update.Message.MessageID
				send <- message
			}
		}
	}
}
