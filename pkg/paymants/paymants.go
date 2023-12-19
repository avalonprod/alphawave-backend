package paymants

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/setupintent"
	"github.com/stripe/stripe-go/sub"
)

type PaymentProvider struct {
	StripeAPIKey string
}

func NewPaymentProvider(stripeAPIKey string) *PaymentProvider {
	stripe.Key = stripeAPIKey
	return &PaymentProvider{
		StripeAPIKey: stripeAPIKey,
	}
}

func (p *PaymentProvider) CreateCustomer(name, email, descr string) (*string, error) {
	params := &stripe.CustomerParams{
		Name:        &name,
		Email:       &email,
		Description: &descr,
	}

	c, err := customer.New(params)

	if err != nil {
		return nil, err
	}

	return &c.ID, nil
}

func (p *PaymentProvider) NewCard(customerID string) (secret *string, err error) {

	params := &stripe.SetupIntentParams{
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Customer: &customerID,
	}
	res, err := setupintent.New(params)
	if err != nil {
		return nil, err
	}
	return &res.ClientSecret, nil
}

// func (p *PaymentProvider) GetListPaymentMethod(customerID string) *[]string {
// 	params := &stripe.PaymentMethodListParams{
// 		Customer: stripe.String(customerID),
// 		Type:     stripe.String("card"),
// 	}
// 	i := paymentmethod.List(params)

// 	for i.Next() {
// 		pm := i.PaymentMethod()
// 	}
// }

func (p *PaymentProvider) CreateSubscription(customerID, priceID string) (*string, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(priceID),
			},
		},
	}

	s, err := sub.New(params)
	if err != nil {
		return nil, err
	}
	return &s.ID, nil
}

// func (p *PaymentProvider) UpdateSubscription(subID, priceID string) (*string, error) {
// 	subItemParams := &stripe.SubscriptionItemListParams{
// 		Subscription: &subID,
// 	}
// 	i := subitem.List(subItemParams)
// 	var si *stripe.SubscriptionItem

// 	for i.Next() {
// 		si = i.SubscriptionItem()
// 		break
// 	}
// 	if si == nil {
// 		return nil, errors.New("doesn't exist items for subscription")
// 	}

// 	subParams := &stripe.SubscriptionParams{
// 		CancelAtPeriodEnd: stripe.Bool(false),
// 		ProrationBehavior: stripe.String(string(stripe.SubscriptionProrationBehaviorCreateProrations)),
// 		Items: []*stripe.SubscriptionItemsParams{
// 			{
// 				ID:   &si.ID,
// 				Plan: &priceID,
// 			},
// 		},
// 		DefaultPaymentMethod: &paymentMethodId,
// 	}
// }

func (p *PaymentProvider) CancelSubscription(subscriptionID string) error {
	cancel := true

	subscriptionParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: &cancel,
	}

	_, err := sub.Update(subscriptionID, subscriptionParams)
	if err != nil {
		return err
	}
	return nil
}
