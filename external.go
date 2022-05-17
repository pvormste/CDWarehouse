package warehouse

type PaymentProvider interface {
	ProcessPayment() error
}

type ChartsProvider interface {
	Notify(title, artist string, amount int) error
}
