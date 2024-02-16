package auction

import (
	"AuctionBot/messages"
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

func Auctioneer(auction Auction, chatID int64, send chan tgbotapi.Chattable, receive tgbotapi.UpdatesChannel) {
	auction.Start()
	activeAuction = auction

	startingMessageText := fmt.Sprintf(messages.START_AUCTION_MESSAGE, auction.Name(), auction.StartPrice(), auction.MinStep())
	startingMessage := tgbotapi.NewMessage(chatID, startingMessageText)
	send <- startingMessage
	duration := 40 * time.Second

	countDown := time.NewTimer(duration)
	counter := 3

	for {
		select {
		case <-countDown.C:
			if counter == 3 {
				countDownMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.COUNTDOWN_THREE_MESSAGE, auction.CurrentPrice()))
				send <- countDownMessage
				counter--
				countDown.Reset(duration)
				continue
			} else if counter == 2 {
				countDownMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.COUNTDOWN_TWO_MESSAGE, auction.CurrentPrice()))
				send <- countDownMessage
				counter--
				countDown.Reset(duration)
				continue
			} else if counter == 1 {
				countDownMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.COUNTDOWN_ONE_MESSAGE, auction.CurrentPrice()))
				send <- countDownMessage
				counter--
				countDown.Reset(duration)
				continue
			} else {
				activeAuction = nil
				auction.End()

				winner, _ := auction.Winner()
				winnerPrice, _ := auction.WinnerPrice()

				messageText := fmt.Sprintf(messages.END_AUCTION_MESSAGE, auction.Name(), winner, winnerPrice)
				message := tgbotapi.NewMessage(chatID, messageText)
				send <- message

				file := tgbotapi.FilePath("./media/leonardo-dicaprio-sold-gif.mp4")
				gif := tgbotapi.NewAnimation(chatID, file)
				send <- gif

				auction.WriteLog()
				return
			}

		case update := <-receive:
			if update.Message != nil {
				bid, err := parseBidMessage(update)
				if err != nil {
					message := tgbotapi.NewMessage(chatID, messages.INVALID_BID_MESSAGE)
					message.ReplyToMessageID = update.Message.MessageID
					send <- message
					continue
				}
				if bid.AuctionName != auction.Name() {
					message := tgbotapi.NewMessage(chatID, messages.NO_AUCTION_MESSAGE)
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
					countDown.Reset(duration)
					counter = 3
				}
				message := tgbotapi.NewMessage(chatID, messageText)
				message.ReplyToMessageID = update.Message.MessageID
				send <- message
			}
		}
	}
}
