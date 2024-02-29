package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type SealedBidAuction struct {
	name    string
	status  string
	history []Bid
	mutex   *sync.Mutex
}

func NewSealedBidAuction(name string) Auction {
	return &SealedBidAuction{
		name:   name,
		status: "Created",
		mutex:  &sync.Mutex{},
	}
}

func (a *SealedBidAuction) Name() string {
	return a.name
}

func (a *SealedBidAuction) StartingMessage() string {
	return fmt.Sprintf(messages.START_SEALED_BID_AUCTION_MESSAGE, a.Name(), a.Name())
}

func (a *SealedBidAuction) Start() {
	a.status = "Started"
}

func (a *SealedBidAuction) End() {
	a.status = "Finished"

	winner := a.Winner()
	for i, b := range a.history {
		if b.Bidder == winner {
			a.history[i].Status = "Winner"
		}
	}
}

func (a *SealedBidAuction) Bid(bid Bid) (string, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, b := range a.history {
		if b.Bidder == bid.Bidder {
			return "", fmt.Errorf(messages.BID_ALREADY_PLACED_MESSAGE)
		}
	}

	if len(a.history) > 0 {
		a.history[len(a.history)-1].Status = "Not winner"
	}
	a.history = append(a.history, bid)

	return messages.SEALED_BID_ACCEPTED_MESSAGE, nil
}

func (a *SealedBidAuction) Winner() string {
	if len(a.history) == 0 {
		return "System"
	}
	maxBid := a.history[0]
	for _, b := range a.history {
		if b.Amount < maxBid.Amount {
			maxBid = b
		}
	}
	return maxBid.Bidder
}

func (a *SealedBidAuction) WinnerPrice() float64 {
	if len(a.history) == 0 {
		return 0
	}
	maxBid := a.history[0]
	for _, b := range a.history {
		if b.Amount < maxBid.Amount {
			maxBid = b
		}
	}
	return maxBid.Amount
}

func (a *SealedBidAuction) WriteLog(name string) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Warn("could not open log file")
	}
	defer file.Close()
	for _, bid := range a.history {
		file.WriteString(fmt.Sprintf("%s %s %f %s\n", bid.Time.Format("2006-01-02-15-04-05"), bid.Bidder, bid.Amount, bid.Status))
	}
	file.WriteString(fmt.Sprintf("Winner: %s\n", a.Winner()))
	file.WriteString(fmt.Sprintf("Winner price: %f\n", a.WinnerPrice()))
}

var sealedBidAuctionBidPattern = regexp.MustCompile(`^/bid (\w+) (\d+(\.\d+)?)$`)

func (a *SealedBidAuction) ParseBid(update tgbotapi.Update) (Bid, error) {
	text := update.Message.Text

	if !sealedBidAuctionBidPattern.MatchString(text) {
		return Bid{}, fmt.Errorf("bid command should be in format %s", sealedBidAuctionBidPattern.String())
	}

	matches := sealedBidAuctionBidPattern.FindStringSubmatch(text)
	auctionName := matches[1]
	amount, _ := strconv.ParseFloat(matches[2], 64)
	return Bid{
		AuctionName: auctionName,
		Bidder:      update.Message.From.UserName,
		Amount:      amount,
		Time:        time.Now(),
	}, nil
}

func (a *SealedBidAuction) Auctioneer() func(auctioneer *Auctioneer) {
	return func(auctioneer *Auctioneer) {
		duration := 5 * time.Second

		countDown := time.NewTimer(duration)
		defer countDown.Stop()

		counter := 1
		for {
			select {
			case bid := <-auctioneer.bidsChannel:
				result, err := auctioneer.auction.Bid(bid)
				if err != nil {
					message := tgbotapi.NewMessage(bid.Update.Message.Chat.ID, err.Error())
					message.ReplyToMessageID = bid.Update.Message.MessageID
					auctioneer.send <- message
					continue
				}
				message := tgbotapi.NewMessage(bid.Update.Message.Chat.ID, result)
				message.ReplyToMessageID = bid.Update.Message.MessageID
				auctioneer.send <- message

			case <-countDown.C:
				switch counter {
				case 1:
					auctioneer.send <- tgbotapi.NewMessage(auctioneer.chatID, messages.SEALED_HALF_TIME_MESSAGE)
				case 0:
					auctioneer.stopChannel <- "Time is Up"
					return
				}
				counter--
				countDown.Reset(duration)
			}
		}
	}
}

func (a *SealedBidAuction) IsPrivateAllowed() bool {
	return true
}
