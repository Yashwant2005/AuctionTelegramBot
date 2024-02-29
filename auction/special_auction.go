package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
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
		Time:   time.Now(),
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

func (a *SpecialAuction) StartingMessage() string {
	return fmt.Sprintf(messages.START_SPECIAL_AUCTION_MESSAGE, a.Name(), a.StartPrice(), a.MinStep(), a.Name())
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

func (a *SpecialAuction) Bid(bid Bid) (string, error) {
	a.history = append(a.history, bid)
	a.currentPrice = bid.Amount

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

func (a *SpecialAuction) WriteLog(name string) {
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
		Amount:      a.CurrentPrice(),
		Time:        time.Now(),
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
				bid := Bid{
					Bidder: "System",
					Amount: a.CurrentPrice() + a.MinStep(),
					Time:   time.Now(),
				}
				a.Bid(bid)
				auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.SPECIAL_AUCTION_PRICE_RAISED_MESSAGE, a.CurrentPrice()))
			case bid := <-auctioneer.bidsChannel:
				a.Bid(bid)
				message := tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.SPECIAL_AUCTION_BID_ACCEPTED_MESSAGE, bid.Bidder))
				message.ReplyToMessageID = bid.Update.Message.MessageID
				auctioneer.send <- message
				auctioneer.stopChannel <- "Auction finished"
				return
			}
		}
	}
}

func (a *SpecialAuction) IsPrivateAllowed() bool {
	return false
}
