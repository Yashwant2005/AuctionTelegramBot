package auction

type Auction interface {
	// Name of auction
	Name() string
	// Start auction
	Start()
	// End auction
	End()
	// StartPrice of auction
	StartPrice() float64
	// CurrentPrice of auction
	CurrentPrice() float64
	// MinStep of auction
	MinStep() float64
	// Bid in auction
	Bid(bidder string, amount float64) (string, error)
	// Winner of auction
	Winner() string
	// WinnerPrice Winner price
	WinnerPrice() float64
	// WriteLog Write log
	WriteLog()
}
