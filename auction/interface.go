package auction

type Auction interface {
	// Name of auction
	Name() string
	// Start auction
	Start() string
	// End auction
	End() string
	// Bid in auction
	Bid(bidder string, amount float64) (string, error)
	// Winner of auction
	Winner() (string, error)
	// WinnerPrice Winner price
	WinnerPrice() (float64, error)
	// WriteLog Write log
	WriteLog()
}
