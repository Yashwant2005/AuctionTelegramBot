package auction

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

type BidResponse struct {
	Message string
	Success bool
	Bid     Bid
}
