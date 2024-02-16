package auction

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type StartAuctionConfig struct {
	Name       string
	StartPrice float64
	MinStep    float64
}

var StartPattern = regexp.MustCompile(`^/start (\w+) (\d+(\.\d+)?) (\d+(\.\d+)?)$`)

func ParseStartAuctionCommand(text string) (StartAuctionConfig, error) {
	if !StartPattern.MatchString(text) {
		return StartAuctionConfig{}, errors.New(fmt.Sprintf("Start command should be in format %s", StartPattern.String()))
	}
	matches := StartPattern.FindStringSubmatch(text)
	startPrice, err := strconv.ParseFloat(matches[2], 32)
	if err != nil {
		return StartAuctionConfig{}, err
	}
	minStep, err := strconv.ParseFloat(matches[4], 32)
	if err != nil {
		return StartAuctionConfig{}, err
	}
	return StartAuctionConfig{
		Name:       matches[1],
		StartPrice: startPrice,
		MinStep:    minStep,
	}, nil
}
