package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"time"
)

type SpecialAuction struct {
	name         string
	startPrice   float64
	currentPrice float64
	minStep      float64
	status       string
	history      []Bid
}

func NewSpecialAuction(name string, startPrice float64, minStep float64) Auction {
	bid := Bid{
		ID:     1,
		Bidder: "System",
		Amount: startPrice,
		Status: "Active",
	}
	return &SpecialAuction{
		name:         name,
		startPrice:   startPrice,
		currentPrice: startPrice,
		minStep:      minStep,
		status:       "Created",
		history:      []Bid{bid},
	}
}

func (a *SpecialAuction) Name() string {
	return a.name
}

func (a *SpecialAuction) StartPrice() float64 {
	return a.startPrice
}

func (a *SpecialAuction) CurrentPrice() float64 {
	return a.currentPrice
}

func (a *SpecialAuction) MinStep() float64 {
	return a.minStep
}

func (a *SpecialAuction) Start() {
	a.status = "Started"
}

func (a *SpecialAuction) End() {
	a.status = "Finished"
}

func (a *SpecialAuction) Bid(bidder string, amount float64) (string, error) {
	if a.status != "Started" {
		return "", fmt.Errorf("auction is not started")
	}

	bid := Bid{
		Bidder: bidder,
		Amount: amount,
		Status: "Active",
	}

	a.history[len(a.history)-1].Status = "Inactive"
	a.history = append(a.history, bid)
	a.currentPrice = amount

	return fmt.Sprintf(messages.ACCEPTED_BID_MESSAGE, bid.Bidder, bid.Amount), nil
}

func (a *SpecialAuction) Winner() string {
	winner := a.history[len(a.history)-1].Bidder
	return winner
}

func (a *SpecialAuction) WinnerPrice() float64 {
	winnerPrice := a.history[len(a.history)-1].Amount
	return winnerPrice
}

func (a *SpecialAuction) WriteLog() {
	// Write log
}

var specialAuctionBidPattern = regexp.MustCompile(`^/bid (\w+)$`)

func (a *SpecialAuction) ParseBid(update tgbotapi.Update) (Bid, error) {
	text := update.Message.Text

	if !specialAuctionBidPattern.MatchString(text) {
		return Bid{}, fmt.Errorf("bid command should be in format %s", specialAuctionBidPattern.String())
	}

	matches := specialAuctionBidPattern.FindStringSubmatch(text)
	auctionName := matches[1]
	return Bid{
		AuctionName: auctionName,
		Bidder:      update.Message.From.UserName,
	}, nil
}

func (a *SpecialAuction) Auctioneer() func(auctioneer *Auctioneer) {
	return func(auctioneer *Auctioneer) {
		duration := 5 * time.Second

		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				a.Bid("System", a.CurrentPrice()+a.MinStep())
				auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.SPECIAL_AUCTION_PRICE_RAISED_MESSAGE, a.CurrentPrice()))
			case bid := <-auctioneer.bidsChannel:
				a.Bid(bid.Bidder, a.CurrentPrice())
				auctioneer.stopChannel <- "Auction finished"
				return
			}
		}
	}
}
