package warehouse

type PaymentProvider interface {
	ProcessPayment() error
}
