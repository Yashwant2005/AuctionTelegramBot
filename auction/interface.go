package auction

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Auction interface {
	// Name of auction
	Name() string
	// Start auction
	Start()
	// End auction
	End()
	// StartingMessage Starting message
	StartingMessage() string
	// Bid in auction
	Bid(bidder string, amount float64) (string, error)
	// Winner of auction
	Winner() string
	// WinnerPrice Winner price
	WinnerPrice() float64
	// WriteLog Write log
	WriteLog()
	// Auctioneer that will be used to notify about auction events
	Auctioneer() func(auctioneer *Auctioneer)
	// ParseBid Parse Bid from string
	ParseBid(bid tgbotapi.Update) (Bid, error)
	// IsPrivateAllowed Is private allowed
	IsPrivateAllowed() bool
}
