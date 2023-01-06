package provider

import (
	"fmt"
	"github.com/doruk581/cdc-wholesale-workshop-go/model"
	"github.com/doruk581/cdc-wholesale-workshop-go/provider/repository"
	"log"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
)

// The Provider verification
func TestPactProvider(t *testing.T) {
	go startInstrumentedProvider()

	pact := createPact()

	// Verify the Provider - Tag-based Published Pacts for any known consumers
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL:    fmt.Sprintf("http://127.0.0.1:%d", port),
		Tags:               []string{"dev"},
		FailIfNoPactsFound: false,
		// Use this if you want to test without the Pact Broker
		// PactURLs:                   []string{filepath.FromSlash(fmt.Sprintf("%s/golistingservice-goproductservice.json", os.Getenv("PACT_DIR")))},
		BrokerURL:                  fmt.Sprintf("%s://%s", os.Getenv("PACT_BROKER_PROTO"), os.Getenv("PACT_BROKER_URL")),
		PublishVerificationResults: true,
		ProviderVersion:            "1.0.0",
		StateHandlers:              stateHandlers,
		RequestFilter:              fixBearerToken,
	})

	if err != nil {
		t.Fatal(err)
	}

}

// Simulates the neeed to set a time-bound authorization token,
// such as an OAuth bearer token
func fixBearerToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only set the correct bearer token, if one was provided in the first place
		if r.Header.Get("Authorization") != "" {
			r.Header.Set("Authorization", getAuthToken())
		}
		next.ServeHTTP(w, r)
	})
}

var stateHandlers = types.StateHandlers{
	"Product trendyolmilla exists": func() error {
		productRepository = trendyolMillaExist
		return nil
	},
	"Product trendyolmilla does not exist": func() error {
		productRepository = trendyolMillaDoesNotExist
		return nil
	},
}

// Starts the provider API with hooks for provider states.
// This essentially mirrors the main.go file, with extra routes added.
func startInstrumentedProvider() {
	mux := GetHTTPHandler()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("API starting: port %d (%s)", port, ln.Addr())
	log.Printf("API terminating: %v", http.Serve(ln, mux))

}

// Configuration / Test Data
var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/../../pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)
var port, _ = utils.GetFreePort()

// Provider States data sets
var trendyolMillaExist = &repository.ProductRepository{
	Products: map[string]*model.Product{
		"TrendyolMilla": &model.Product{
			Name:     "TrendyolMilla",
			Category: "Short",
			Color:    "green",
			ID:       10,
		},
	},
}

var trendyolMillaDoesNotExist = &repository.ProductRepository{}

var trendyolMillaNotAuthorized = &repository.ProductRepository{
	Products: map[string]*model.Product{
		"TrendyolMilla": &model.Product{
			Name:     "TrendyolMilla",
			Category: "Short",
			Color:    "green",
			ID:       10,
		},
	},
}

// Setup the Pact client.
func createPact() dsl.Pact {
	return dsl.Pact{
		Provider:                 "GoProductService",
		LogDir:                   logDir,
		PactDir:                  pactDir,
		DisableToolValidityCheck: true,
		LogLevel:                 "INFO",
	}
}
