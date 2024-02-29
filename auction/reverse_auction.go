package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Bid struct {
	// Bid ID
	ID int
	// AuctionName bidder
	AuctionName string
	// Bidder
	Bidder string
	// Bid amount
	Amount float64
	// Bid status
	Status string
	// Bid time
	Time time.Time
	// Bid telegram message
	Update tgbotapi.Update
}

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

func (a *ReverseAuction) Name() string {
	return a.name
}

func (a *ReverseAuction) StartPrice() float64 {
	return a.startPrice
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
	if amount > a.currentPrice-a.minStep {
		return "", fmt.Errorf(messages.INVALID_BID_AMOUNT_MESSAGE, a.CurrentPrice(), a.CurrentPrice()-a.MinStep())
	}
	bid := Bid{
		ID:     len(a.history) + 1,
		Bidder: bidder,
		Amount: amount,
		Status: "Active",
		Time:   time.Now(),
	}
	a.history[len(a.history)-1].Status = "Inactive"
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

func (a *ReverseAuction) WriteLog() {
	name := fmt.Sprintf("./log/%s-%s.log", time.Now().Format("2006-01-02-15-04-05"), a.Name())
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

func NewReverseAuction(name string, startPrice float64, minStep float64) Auction {
	auction := &ReverseAuction{
		name:         name,
		startPrice:   startPrice,
		currentPrice: startPrice,
		minStep:      minStep,
		status:       "Created",
	}
	bid := Bid{
		ID:     1,
		Bidder: "System",
		Amount: startPrice + minStep,
		Status: "Active",
	}
	auction.history = append(auction.history, bid)
	return auction
}
