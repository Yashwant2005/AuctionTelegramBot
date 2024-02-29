package auction

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReverseAuction(t *testing.T) {
	auction := NewReverseAuction("test", 100, 10).(Auction)

	type BidTest struct {
		bidder        string
		amount        float64
		expectedError bool
	}

	var tests = []BidTest{
		{
			bidder:        "bidder1",
			amount:        90,
			expectedError: false,
		},
		{
			bidder:        "bidder2",
			amount:        80,
			expectedError: false,
		},
		{
			bidder:        "bidder2",
			amount:        100,
			expectedError: true,
		},
		{
			bidder:        "bidder2",
			amount:        75,
			expectedError: true,
		},
		{
			bidder:        "bidder1",
			amount:        10,
			expectedError: false,
		},
	}

	for _, test := range tests {
		bid := Bid{
			Bidder: test.bidder,
			Amount: test.amount,
		}
		_, err := auction.Bid(bid)
		if test.expectedError {
			require.NotNil(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	require.Equal(t, 10.0, auction.WinnerPrice())
	require.Equal(t, "bidder1", auction.Winner())
}
