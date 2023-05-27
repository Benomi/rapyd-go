package rapyd

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/benomi/rapyd-go/resources"
	"github.com/pkg/errors"
)

const (
	createCheckoutPath   = "/v1/checkout"
	retrieveCheckoutPath = "/v1/checkout/"

	createCustomerPath   = "/v1/customers"
	updateCustomerPath   = "/v1/customers/"
	retrieveCustomerPath = "/v1/customers/"

	getPaymentFieldsPath  = "/v1/payment_methods/required_fields/"
	getPaymentMethodsPath = "/v1/payment_methods/country?country="
)

type Client interface {
	Resolve(path string) string
	GetSigned(path string) ([]byte, error)
	PostSigned(data interface{}, path string) ([]byte, error)
	DeleteSigned(path string) ([]byte, error)

	CreateCheckout(data resources.CreateCheckout) (*resources.CheckoutResponse, error)
	RetrieveCheckout(checkoutID string) (*resources.CheckoutResponse, error)

	CreateCustomer(data resources.Customer) (*resources.CustomerResponse, error)
	RetrieveCustomer(customerID string) (*resources.RetrieveCustomerResponse, error)
	UpdateCustomer(customerID string, data resources.Customer) (*resources.CustomerResponse, error)

	GetPaymentMethodFields(method string) (*resources.PaymentMethodRequiredFieldsResponse, error)
	GetCountryPaymentMethods(country string) (*resources.CountryPaymentMethodsResponse, error)

	ValidateWebhook(r *http.Request) bool
}

type client struct {
	Signer
	*http.Client
	url *url.URL
}

func NewClient(signer Signer, url *url.URL, cli *http.Client) Client {
	return &client{
		Signer: signer,
		Client: cli,
		url:    url,
	}
}

func (c *client) Resolve(path string) string {
	endpoint, err := url.Parse(path)
	if err != nil {
		panic(errors.New("error parsing path"))
	}
	return c.url.ResolveReference(endpoint).String()
}

func (c *client) GetSigned(path string) ([]byte, error) {
	request, err := http.NewRequest("GET", c.Resolve(path), nil)

	err = c.signRequest(request, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error signing request")
	}

	r, err := c.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "error sending request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		errorResponse, _ := ioutil.ReadAll(r.Body)
		return nil, errors.Errorf("error: got status code %d, response %s", r.StatusCode, string(errorResponse))
	}

	return ioutil.ReadAll(r.Body)
}

func (c *client) PostSigned(data interface{}, path string) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling data")
	}

	request, err := http.NewRequest("POST", c.Resolve(path), bytes.NewBuffer(body))

	err = c.signRequest(request, body)
	if err != nil {
		return nil, errors.Wrap(err, "error signing request")
	}

	r, err := c.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "error sending request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		errorResponse, _ := ioutil.ReadAll(r.Body)
		return nil, errors.Errorf("error: got status code %d, response %s", r.StatusCode, string(errorResponse))
	}

	return ioutil.ReadAll(r.Body)
}

func (c *client) DeleteSigned(path string) ([]byte, error) {
	request, err := http.NewRequest("DELETE", c.Resolve(path), nil)

	err = c.signRequest(request, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error signing request")
	}

	r, err := c.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "error sending request")
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		errorResponse, _ := ioutil.ReadAll(r.Body)
		return nil, errors.Errorf("error: got status code %d, response %s", r.StatusCode, string(errorResponse))
	}

	return ioutil.ReadAll(r.Body)
}

func (c *client) CreateCheckout(data resources.CreateCheckout) (*resources.CheckoutResponse, error) {
	response, err := c.PostSigned(data, createCheckoutPath)
	if err != nil {
		return nil, errors.Wrap(err, "error sending create checkout request")
	}

	var body resources.CheckoutResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return &body, nil
}

func (c *client) RetrieveCheckout(checkoutID string) (*resources.CheckoutResponse, error) {
	response, err := c.GetSigned(retrieveCheckoutPath + checkoutID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting checkout")
	}

	var body resources.CheckoutResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling checkout response")
	}

	return &body, nil
}

func (c *client) CreateCustomer(data resources.Customer) (*resources.CustomerResponse, error) {
	response, err := c.PostSigned(data, createCustomerPath)
	if err != nil {
		return nil, errors.Wrap(err, "error sending create customer request")
	}

	var body resources.CustomerResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return &body, nil
}

func (c *client) RetrieveCustomer(customerID string) (*resources.RetrieveCustomerResponse, error) {
	response, err := c.GetSigned(retrieveCustomerPath + customerID)
	if err != nil {
		return nil, errors.Wrap(err, "error sending retrieve customer request")
	}

	var body resources.RetrieveCustomerResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return &body, nil
}

func (c *client) UpdateCustomer(customerID string, data resources.Customer) (*resources.CustomerResponse, error) {
	response, err := c.PostSigned(data, updateCustomerPath+customerID)
	if err != nil {
		return nil, errors.Wrap(err, "error sending update customer request")
	}

	var body resources.CustomerResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return &body, nil
}

func (c *client) GetPaymentMethodFields(method string) (*resources.PaymentMethodRequiredFieldsResponse, error) {
	response, err := c.GetSigned(getPaymentFieldsPath + method)
	if err != nil {
		return nil, errors.Wrap(err, "error getting payment method fields")
	}

	var body resources.PaymentMethodRequiredFieldsResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return &body, nil
}

func (c *client) GetCountryPaymentMethods(country string) (*resources.CountryPaymentMethodsResponse, error) {
	response, err := c.GetSigned(getPaymentMethodsPath + country)
	if err != nil {
		return nil, errors.Wrap(err, "error getting country payment methods")
	}

	var body resources.CountryPaymentMethodsResponse

	err = json.Unmarshal(response, &body)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling response")
	}

	return &body, nil
}

func (c *client) ValidateWebhook(r *http.Request) bool {
	if webhookBytes, err := ioutil.ReadAll(r.Body); err == nil {
		if err := r.Body.Close(); err != nil {
			return false
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(webhookBytes))

		data := SignatureData{
			Path:      fmt.Sprintf("https://%s%s", r.Host, r.RequestURI),
			Salt:      r.Header.Get(SaltHeader),
			Timestamp: r.Header.Get(TimestampHeader),
			Body:      string(webhookBytes),
		}

		generatedSignature := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(c.signData(data))))
		return generatedSignature == r.Header.Get(SignatureHeader)
	}

	return false
}
