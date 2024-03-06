package auction

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type StartAuctionConfig struct {
	Type       string
	Name       string
	StartPrice float64
	MinStep    float64
}

var startReversePattern = regexp.MustCompile(`^/start (\w+) (\w+) (\d+(\.\d+)?) (\d+(\.\d+)?)$`)
var startDutchReversePattern = regexp.MustCompile(`^/start (\w+) (\w+) (\d+(\.\d+)?) (\d+(\.\d+)?)$`)
var startSealedBidPattern = regexp.MustCompile(`^/start (\w+) (\w+)$`)

func ParseStartAuctionCommand(text string) (StartAuctionConfig, error) {
	var startPattern *regexp.Regexp
	auctionType := strings.Split(text, " ")[1]
	switch auctionType {
	case "reverse_auction":
		startPattern = startReversePattern
	case "dutch_reverse_auction":
		startPattern = startDutchReversePattern
	case "sealed_bid_auction":
		startPattern = startSealedBidPattern
	default:
		return StartAuctionConfig{}, errors.New("auction type should be reverse_auction or dutch_reverse_auction or sealed_bid_auction")
	}
	if !startPattern.MatchString(text) {
		return StartAuctionConfig{}, errors.New(fmt.Sprintf("Start command should be in format %s", startPattern.String()))
	}
	matches := startPattern.FindStringSubmatch(text)
	switch auctionType {
	case "reverse_auction", "dutch_reverse_auction":
		startPrice, _ := strconv.ParseFloat(matches[3], 64)
		minStep, _ := strconv.ParseFloat(matches[5], 64)

		return StartAuctionConfig{
			Type:       matches[1],
			Name:       matches[2],
			StartPrice: startPrice,
			MinStep:    minStep,
		}, nil
	case "sealed_bid_auction":
		return StartAuctionConfig{
			Type: matches[1],
			Name: matches[2],
		}, nil
	default:
		return StartAuctionConfig{}, errors.New("auction type should be reverse_auction or dutch_reverse_auction or sealed_bid_auction")
	}
}

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
