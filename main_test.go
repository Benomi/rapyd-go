package rapyd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	endpoint = "https://sandboxapi.rapyd.net"
)

var accessKey = os.Getenv("RAPYD_ACCESS_KEY")
var secretKey = os.Getenv("RAPYD_SECRET_KEY")

func TestClient_GetCountryPaymentMethods(t *testing.T) {
	addr, err := url.Parse(endpoint)
	assert.NoError(t, err)

	rapyd := NewClient(NewRapydSigner([]byte(accessKey), []byte(secretKey)), addr, http.DefaultClient)

	resp, err := rapyd.GetCountryPaymentMethods("US")
	fmt.Println(resp)

	assert.NoError(t, err)
}

func TestClient_GetPaymentMethodFields(t *testing.T) {
	addr, err := url.Parse(endpoint)
	assert.NoError(t, err)

	rapyd := NewClient(NewRapydSigner([]byte(accessKey), []byte(secretKey)), addr, http.DefaultClient)

	resp, err := rapyd.GetPaymentMethodFields("us_debit_visa_card")
	fmt.Println(resp)

	assert.NoError(t, err)
}