package auction

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFirstPrice(t *testing.T) {
	auction := NewFirstPriceAuction("test", 100, 10).(Auction)

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
		_, err := auction.Bid(test.bidder, test.amount)
		if test.expectedError {
			require.NotNil(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	require.Equal(t, 10.0, auction.CurrentPrice())
	require.Equal(t, "bidder1", auction.Winner())
}
