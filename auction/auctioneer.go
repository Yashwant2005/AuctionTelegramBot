package auction

import (
	"AuctionBot/messages"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"time"
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
}

func NewAuctioneer(config StartAuctionConfig, chatID int64, send chan tgbotapi.Chattable) Auctioneer {
	var auction Auction
	if config.Type == "reverse_auction" {
		auction = NewReverseAuction(config.Name, config.StartPrice, config.MinStep)
	} else {
		return Auctioneer{}
	}
	activeAuction = auction
	auction.Start()
	bids := make(chan Bid)

	return Auctioneer{
		auction:     auction,
		chatID:      chatID,
		send:        send,
		bidsChannel: bids,
	}
}

var bidPattern = regexp.MustCompile(`^/bid (\w+) (\d+(\.\d+)?)$`)

func parseBidMessage(update tgbotapi.Update) (Bid, error) {
	text := update.Message.Text

	if !bidPattern.MatchString(text) {
		return Bid{}, fmt.Errorf("bid command should be in format %s", bidPattern.String())
	}

	matches := bidPattern.FindStringSubmatch(text)
	auctionName := matches[1]
	amount, _ := strconv.ParseFloat(matches[2], 64)
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

func (a *Auctioneer) ReverseAuctionStoppingRule(stopChannel chan string) {
	duration := 15 * time.Second

	countDown := time.NewTimer(duration)
	counter := 3
	for {
		select {
		case bid := <-a.bidsChannel:
			result, err := a.auction.Bid(bid.Bidder, bid.Amount)
			if err != nil {
				message := tgbotapi.NewMessage(a.chatID, err.Error())
				a.send <- message
				continue
			}
			message := tgbotapi.NewMessage(a.chatID, result)
			a.send <- message

			countDown.Reset(duration)
			counter = 3

		case <-countDown.C:
			switch counter {
			case 3:
				a.send <- tgbotapi.NewMessage(a.chatID, fmt.Sprintf(messages.COUNTDOWN_THREE_MESSAGE, a.auction.CurrentPrice()))
			case 2:
				a.send <- tgbotapi.NewMessage(a.chatID, fmt.Sprintf(messages.COUNTDOWN_TWO_MESSAGE, a.auction.CurrentPrice()))
			case 1:
				a.send <- tgbotapi.NewMessage(a.chatID, fmt.Sprintf(messages.COUNTDOWN_ONE_MESSAGE, a.auction.CurrentPrice()))
			case 0:
				stopChannel <- fmt.Sprint("Auction ended")
				return
			}
			counter--
			countDown.Reset(duration)
		}
	}
}

func (a *Auctioneer) Run(receive tgbotapi.UpdatesChannel) {
	activeAuction = a.auction
	a.auction.Start()

	startingMessage := tgbotapi.NewMessage(a.chatID, fmt.Sprintf(messages.START_AUCTION_MESSAGE, a.auction.Name(), a.auction.StartPrice(), a.auction.MinStep()))
	a.send <- startingMessage

	stoppingRuleChannel := make(chan string)

	go a.ReverseAuctionStoppingRule(stoppingRuleChannel)

	for {
		select {
		case <-stoppingRuleChannel:
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
				bid, err := parseBidMessage(update)
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
