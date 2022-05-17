package warehouse

type PaymentProvider interface {
	ProcessPayment() error
}

type ChartsNotifier interface {
	Notify(title, artist string, amount int) error
}
