package resources

type CreateCheckout struct {
	Amount                      float64  `json:"amount"`
	Country                     string   `json:"country"`
	Currency                    string   `json:"currency"`
	Customer                    string   `json:"customer"`
	CancelCheckoutURL           string   `json:"cancel_checkout_url"`
	MerchantReference 				 string  `json:"merchant_reference_id"`
	CompleteCheckoutURL         string   `json:"complete_checkout_url"`
	ErrorCheckoutURL            string   `json:"error_checkout_url"`
	PaymentMethodTypeCategories []string `json:"payment_method_type_categories"`
	Expiration                  *int64   `json:"expiration"`
	RequestedCurrency           *string  `json:"requested_currency,omitempty"`
}

type CheckoutResponse struct {
	Data Data `json:"data"`
}
