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

type DutchReverseAuction struct {
	name         string
	startPrice   float64
	currentPrice float64
	minStep      float64
	status       string
	history      []Bid
}

func NewDutchReverseAuction(name string, startPrice float64, minStep float64) Auction {
	bid := Bid{
		ID:     1,
		Bidder: "System",
		Amount: startPrice,
		Time:   time.Now(),
	}
	return &DutchReverseAuction{
		name:         name,
		startPrice:   startPrice,
		currentPrice: startPrice,
		minStep:      minStep,
		status:       "Created",
		history:      []Bid{bid},
	}
}

func (a *DutchReverseAuction) Name() string {
	return a.name
}

func (a *DutchReverseAuction) StartPrice() float64 {
	return a.startPrice
}

func (a *DutchReverseAuction) StartingMessage() string {
	return fmt.Sprintf(messages.START_DUTCH_REVERSE_AUCTION_MESSAGE, a.Name(), a.StartPrice(), a.MinStep(), a.Name())
}

func (a *DutchReverseAuction) CurrentPrice() float64 {
	return a.currentPrice
}

func (a *DutchReverseAuction) MinStep() float64 {
	return a.minStep
}

func (a *DutchReverseAuction) Start() {
	a.status = "Started"
}

func (a *DutchReverseAuction) End() {
	a.status = "Finished"
}

func (a *DutchReverseAuction) Bid(bid Bid) (string, error) {
	a.history = append(a.history, bid)
	a.currentPrice = bid.Amount

	return fmt.Sprintf(messages.ACCEPTED_BID_MESSAGE, bid.Bidder, bid.Amount), nil
}

func (a *DutchReverseAuction) Winner() string {
	winner := a.history[len(a.history)-1].Bidder
	return winner
}

func (a *DutchReverseAuction) WinnerPrice() float64 {
	winnerPrice := a.history[len(a.history)-1].Amount
	return winnerPrice
}

func (a *DutchReverseAuction) WriteLog(name string) {
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

var dutchReverseAuctionBidPattern = regexp.MustCompile(`^/bid (\w+)$`)

func (a *DutchReverseAuction) ParseBid(update tgbotapi.Update) (Bid, error) {
	text := update.Message.Text

	if !dutchReverseAuctionBidPattern.MatchString(text) {
		return Bid{}, fmt.Errorf("bid command should be in format %s", dutchReverseAuctionBidPattern.String())
	}

	matches := dutchReverseAuctionBidPattern.FindStringSubmatch(text)
	auctionName := matches[1]
	return Bid{
		AuctionName: auctionName,
		Bidder:      update.Message.From.UserName,
		Amount:      a.CurrentPrice(),
		Time:        time.Now(),
	}, nil
}

func (a *DutchReverseAuction) Auctioneer() func(auctioneer *Auctioneer) {
	return func(auctioneer *Auctioneer) {
		duration := 10 * time.Second

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
				auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.DUTCH_REVERSE_AUCTION_PRICE_RAISED_MESSAGE, a.CurrentPrice()))
			case bid := <-auctioneer.bidsChannel:
				a.Bid(bid)
				message := tgbotapi.NewMessage(auctioneer.chatID, fmt.Sprintf(messages.DUTCH_REVERSE_AUCTION_BID_ACCEPTED_MESSAGE, bid.Bidder))
				message.ReplyToMessageID = bid.Update.Message.MessageID
				auctioneer.send <- message
				auctioneer.stopChannel <- "Auction finished"
				return
			}
		}
	}
}

func (a *DutchReverseAuction) IsPrivateAllowed() bool {
	return false
}
