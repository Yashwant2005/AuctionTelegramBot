package auction

import (
	"fmt"
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
}

type FirstPriceAuction struct {
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

func (a *FirstPriceAuction) Name() string {
	return a.name
}

func (a *FirstPriceAuction) Start() string {
	a.status = "Started"
	text := fmt.Sprintf("Auction %s has started. Starting price is %f and minimum step for new bid is %f. If you want to bid, you have to send a message with format \"/bid %s your_bid\".", a.Name(), a.startPrice, a.minStep, a.Name())
	return text
}

func (a *FirstPriceAuction) End() string {
	a.status = "Finished"
	text := fmt.Sprintf("Auction %s has ended.", a.Name())
	return text
}

func (a *FirstPriceAuction) Bid(bidder string, amount float64) (string, error) {
	if a.status == "Created" {
		return "", fmt.Errorf("auction %s is not started yet", a.Name())
	}
	if a.status == "Finished" {
		return "", fmt.Errorf("auction %s is already finished", a.Name())
	}
	if amount > a.currentPrice-a.minStep {
		return "", fmt.Errorf("bid amount should be less than or equal to %f", a.currentPrice-a.minStep)
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
	return fmt.Sprintf("bid %d from %s with amount %f has been accepted", bid.ID, bid.Bidder, bid.Amount), nil
}

func (a *FirstPriceAuction) Winner() (string, error) {
	if a.status != "Finished" {
		return "", fmt.Errorf("auction %s is not finished yet", a.Name())
	}
	winner := a.history[len(a.history)-1].Bidder
	return winner, nil
}

func (a *FirstPriceAuction) WinnerPrice() (float64, error) {
	if a.status != "Finished" {
		return 0, fmt.Errorf("auction %s is not finished yet", a.Name())
	}
	winnerPrice := a.history[len(a.history)-1].Amount
	return winnerPrice, nil
}

func (a *FirstPriceAuction) WriteLog() {
	name := fmt.Sprintf("./log/%s-%s.log", a.Name(), time.Now().Format("2006-01-02-15-04-05"))
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

func NewFirstPriceAuction(name string, startPrice float64, minStep float64) Auction {
	auction := &FirstPriceAuction{
		name:         name,
		startPrice:   startPrice,
		currentPrice: startPrice,
		minStep:      minStep,
		status:       "Created",
	}
	bid := Bid{
		ID:     1,
		Bidder: "System",
		Amount: startPrice,
		Status: "Active",
	}
	auction.history = append(auction.history, bid)
	return auction
}
