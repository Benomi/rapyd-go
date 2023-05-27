package resources

type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Regex       string `json:"regex"`
	Description string `json:"description"`
	Required    bool   `json:"is_required"`
	Updatable   bool   `json:"is_updatable"`
}

type PaymentMethod struct {
	Fields map[string]interface{} `json:"fields"`
	Type   string                 `json:"type"`
}

type PaymentMethodRequiredFields struct {
	Type                 string  `json:"type"`
	Fields               []Field `json:"fields"`
	MethodOptions        []Field `json:"payment_method_options"`
	PaymentOptions       []Field `json:"payment_options"`
	MinExpirationSeconds int64   `json:"minimum_expiration_seconds"`
	MaxExpirationSeconds int64   `json:"maximum_expiration_seconds"`
}

type PaymentMethodRequiredFieldsResponse struct {
	Data PaymentMethodRequiredFields `json:"data"`
}

type CountryPaymentMethod struct {
	Type                   string   `json:"type"`
	Name                   string   `json:"name"`
	Category               string   `json:"category"`
	Image                  string   `json:"image"`
	Country                string   `json:"country"`
	PaymentFlowType        string   `json:"payment_flow_type"`
	Currencies             []string `json:"currencies"`
	Status                 int      `json:"status"`
	IsCancelable           bool     `json:"is_cancelable"`
	PaymentOptions         []Field  `json:"payment_options"`
	IsExpirable            bool     `json:"is_expirable"`
	IsOnline               bool     `json:"is_online"`
	IsRefundable           bool     `json:"is_refundable"`
	IsVirtual              bool     `json:"is_virtual"`
	MultipleOverageAllowed bool     `json:"multiple_overage_allowed"`
	IsTokenizable          bool     `json:"is_tokenizable"`
	MinExpirationSeconds   int64    `json:"minimum_expiration_seconds"`
	MaxExpirationSeconds   int64    `json:"maximum_expiration_seconds"`
	VirtualType            string   `json:"virtual_payment_method_type"`
}

type CountryPaymentMethodsResponse struct {
	Data []CountryPaymentMethod `json:"data"`
}
