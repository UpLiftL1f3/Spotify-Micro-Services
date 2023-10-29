package card

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.createPaymentIntent(currency, amount)
}

func (c *Card) createPaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// params.AddMetadata("key", "value") //! If you want to add information to the transaction

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}

	return pi, "", err
}

func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card is Expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Your card is Expired"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Your card was declined"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "Your card is Expired"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "Your card was declined"
	case stripe.ErrorCodeInsufficientFunds:
		msg = "Your card is Expired"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your card was declined"
	default:
		msg = "Your card was declined"
	}

	return msg
}
