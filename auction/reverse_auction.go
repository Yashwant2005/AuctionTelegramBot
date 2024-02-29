package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"time"
)

type ReverseAuction struct {
	// Auction name
	name string
	// Auction start price
	startPrice float64
	// Auction current price
	currentPrice float64
	// Auction min step
	minStep float64
	// Auction status
	status string
	// Auction history
	history []Bid
	// Telegram Bot
}

func NewReverseAuction(name string, startPrice float64, minStep float64) Auction {
	bid := Bid{
		ID:     1,
		Bidder: "System",
		Amount: startPrice + minStep,
		Status: "Active",
		Time:   time.Now(),
	}
	return &ReverseAuction{
		name:         name,
		startPrice:   startPrice,
		currentPrice: startPrice,
		minStep:      minStep,
		status:       "Created",
		history:      []Bid{bid},
	}
}

func (a *ReverseAuction) Name() string {
	return a.name
}

func (a *ReverseAuction) StartPrice() float64 {
	return a.startPrice
}

func (a *ReverseAuction) StartingMessage() string {
	return fmt.Sprintf(messages.START_REVERSE_AUCTION_MESSAGE, a.Name(), a.StartPrice(), a.MinStep(), a.Name())
}

func (a *ReverseAuction) CurrentPrice() float64 {
	return a.currentPrice
}

func (a *ReverseAuction) MinStep() float64 {
	return a.minStep
}

func (a *ReverseAuction) Start() {
	a.status = "Started"
}

func (a *ReverseAuction) End() {
	a.status = "Finished"
}

func (a *ReverseAuction) Bid(bidder string, amount float64) (string, error) {
	bid := Bid{
		ID:     len(a.history) + 1,
		Bidder: bidder,
		Amount: amount,
		Time:   time.Now(),
	}

	if amount > a.currentPrice-a.minStep {
		bid.Status = "Not accepted"
		a.history = append(a.history, bid)
		return "", fmt.Errorf(messages.INVALID_BID_AMOUNT_MESSAGE, a.CurrentPrice(), a.CurrentPrice()-a.MinStep())
	}

	a.history[len(a.history)-2].Status = "Inactive"
	bid.Status = "Active"
	a.history = append(a.history, bid)
	a.currentPrice = amount
	return fmt.Sprintf(messages.ACCEPTED_BID_MESSAGE, bid.Bidder, bid.Amount), nil
}

func (a *ReverseAuction) Winner() string {
	winner := a.history[len(a.history)-1].Bidder
	return winner
}

func (a *ReverseAuction) WinnerPrice() float64 {
	winnerPrice := a.history[len(a.history)-1].Amount
	return winnerPrice
}

func (a *ReverseAuction) WriteLog(name string) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Warn("could not open log file")
	}
	defer file.Close()
	for _, bid := range a.history {
		file.WriteString(fmt.Sprintf("%s %s %f %s\n", bid.Time.Format("2006-01-02-15-04-05"), bid.Bidder, bid.Amount, bid.Status))
	}
	file.WriteString(fmt.Sprintf("Winner: %s\n", a.history[len(a.history)-1].Bidder))
	file.WriteString(fmt.Sprintf("Winner price: %f\n", a.history[len(a.history)-1].Amount))
}

var reverseAuctionBidPattern = regexp.MustCompile(`^/bid (\w+) (\d+(\.\d+)?)$`)

func (a *ReverseAuction) ParseBid(update tgbotapi.Update) (Bid, error) {
	text := update.Message.Text

	if !reverseAuctionBidPattern.MatchString(text) {
		return Bid{}, fmt.Errorf("bid command should be in format %s", reverseAuctionBidPattern.String())
	}

	matches := reverseAuctionBidPattern.FindStringSubmatch(text)
	auctionName := matches[1]
	amount, _ := strconv.ParseFloat(matches[2], 64)
	return Bid{
		AuctionName: auctionName,
		Bidder:      update.Message.From.UserName,
		Amount:      amount,
	}, nil
}

func (a *ReverseAuction) Auctioneer() func(auctioneer *Auctioneer) {
	return func(auctioneer *Auctioneer) {
		duration := 5 * time.Second

		countDown := time.NewTimer(duration)
		defer countDown.Stop()

		counter := 3
		for {
			select {
			case bid := <-auctioneer.bidsChannel:
				result, err := auctioneer.auction.Bid(bid.Bidder, bid.Amount)
				if err != nil {
					message := tgbotapi.NewMessage(auctioneer.chatID, err.Error())
					message.ReplyToMessageID = bid.Update.Message.MessageID
					auctioneer.send <- message
					continue
				}
				message := tgbotapi.NewMessage(auctioneer.chatID, result)
				message.ReplyToMessageID = bid.Update.Message.MessageID
				auctioneer.send <- message

				countDown.Reset(duration)
				counter = 3

			case <-countDown.C:
				switch counter {
				case 3:
					auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.COUNTDOWN_THREE_MESSAGE, a.CurrentPrice()))
				case 2:
					auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.COUNTDOWN_TWO_MESSAGE, a.CurrentPrice()))
				case 1:
					auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.COUNTDOWN_ONE_MESSAGE, a.CurrentPrice()))
				case 0:
					auctioneer.stopChannel <- "Auction finished"
					return
				}
				counter--
				countDown.Reset(duration)
			}
		}
	}
}

func (a *ReverseAuction) IsPrivateAllowed() bool {
	return false
}
