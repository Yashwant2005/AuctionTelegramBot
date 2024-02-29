package auction

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"time"
)

type StartAuctionConfig struct {
	Type       string
	Name       string
	StartPrice float64
	MinStep    float64
}

var StartPattern = regexp.MustCompile(`^/start (\w+) (\w+) (\d+(\.\d+)?) (\d+(\.\d+)?)$`)

func ParseStartAuctionCommand(text string) (StartAuctionConfig, error) {
	if !StartPattern.MatchString(text) {
		return StartAuctionConfig{}, errors.New(fmt.Sprintf("Start command should be in format %s", StartPattern.String()))
	}
	matches := StartPattern.FindStringSubmatch(text)
	if matches[1] != "reverse_auction" && matches[1] != "special_auction" {
		return StartAuctionConfig{}, errors.New("Auction type should be reverse_auction or special_auction")

	}
	startPrice, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return StartAuctionConfig{}, err
	}
	minStep, err := strconv.ParseFloat(matches[5], 64)
	if err != nil {
		return StartAuctionConfig{}, err
	}
	return StartAuctionConfig{
		Type:       matches[1],
		Name:       matches[2],
		StartPrice: startPrice,
		MinStep:    minStep,
	}, nil
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
