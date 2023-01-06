package client

import (
	"encoding/json"
	"fmt"
	"github.com/doruk581/cdc-wholesale-workshop-go/model"
	"testing"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func TestClientUnit_GetProduct(t *testing.T) {
	productID := 10

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), fmt.Sprintf("/product/%d", productID))
		product, _ := json.Marshal(model.Product{
			Name:     "TrendyolMilla",
			Category: "Clothing",
			ID:       productID,
			Color:    "green",
		})
		rw.Write([]byte(product))
	}))
	defer server.Close()

	u, _ := url.Parse(server.URL)
	client := &Client{
		BaseURL: u,
	}
	product, err := client.GetProduct(productID)
	assert.NoError(t, err)

	assert.Equal(t, product.ID, productID)
}
