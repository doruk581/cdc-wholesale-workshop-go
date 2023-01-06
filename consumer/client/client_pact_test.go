package client

import (
	"fmt"
	"github.com/doruk581/cdc-wholesale-workshop-go/model"
	"os"
	"testing"

	"net/url"

	"github.com/pact-foundation/pact-go/dsl"
)

var commonHeaders = dsl.MapMatcher{
	"Content-Type":         term("application/json; charset=utf-8", `application\/json`),
	"X-Api-Correlation-Id": dsl.Like("100"),
}

var headersWithToken = dsl.MapMatcher{
	"Authorization": dsl.Like("Bearer 2019-01-01"),
}

var u *url.URL
var client *Client

func TestMain(m *testing.M) {
	var exitCode int

	if os.Getenv("PACT_TEST") != "" {

		fmt.Println("PACT_TEST not null")
		// Setup Pact and related test stuff
		setup()

		// Run all the tests
		exitCode = m.Run()

		// Shutdown the Mock Service and Write pact files to disk
		if err := pact.WritePact(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pact.Teardown()
	} else {
		exitCode = m.Run()
	}
	os.Exit(exitCode)
}

func TestClientPact_GetProduct(t *testing.T) {
	t.Run("the product exists", func(t *testing.T) {
		id := 10

		pact.
			AddInteraction().
			Given("Product TrendyolMilla exists").
			UponReceiving("A request to fetch product 'TrendyolMilla'").
			WithRequest(request{
				Method:  "GET",
				Path:    term("/product/10", "/product/[0-9]+"),
				Headers: headersWithToken,
			}).
			WillRespondWith(dsl.Response{
				Status:  200,
				Body:    dsl.Match(model.Product{}),
				Headers: commonHeaders,
			})

		err := pact.Verify(func() error {
			product, err := client.WithToken("2019-01-01").GetProduct(id)

			// Assert basic fact
			if product.ID != id {
				return fmt.Errorf("wanted product with ID %d but got %d", id, product.ID)
			}

			return err
		})

		if err != nil {
			t.Fatalf("Error on Verify: %v", err)
		}
	})
}

// Common test data
var pact dsl.Pact

// Aliases
var term = dsl.Term

type request = dsl.Request

func setup() {
	pact = createPact()

	// Proactively start service to get access to the port
	pact.Setup(true)

	u, _ = url.Parse(fmt.Sprintf("http://localhost:%d", pact.Server.Port))

	client = &Client{
		BaseURL: u,
	}

}

func createPact() dsl.Pact {
	return dsl.Pact{
		Consumer:                 os.Getenv("CONSUMER_NAME"),
		Provider:                 os.Getenv("PROVIDER_NAME"),
		LogDir:                   os.Getenv("LOG_DIR"),
		PactDir:                  os.Getenv("PACT_DIR"),
		LogLevel:                 "INFO",
		PactFileWriteMode:        "overwrite",
		DisableToolValidityCheck: true,
	}
}
