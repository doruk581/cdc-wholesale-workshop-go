package provider

import (
	"encoding/json"
	"github.com/doruk581/cdc-wholesale-workshop-go/model"
	"github.com/doruk581/cdc-wholesale-workshop-go/provider/repository"
	"net/http"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var productRepository = &repository.ProductRepository{
	Products: map[string]*model.Product{
		"TrendyolMilla": &model.Product{
			Name:     "TrendyolMilla",
			Category: "Short",
			Color:    "green",
			ID:       10,
		},
	},
}

func getAuthToken() string {
	return fmt.Sprintf("Bearer %s", time.Now().Format("2006-01-02T15:04"))
}

func WithCorrelationID(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New()
		w.Header().Set("X-Api-Correlation-Id", uuid.String())
		h.ServeHTTP(w, r)
	}
}

// IsAuthenticated checks for a correct bearer token
func IsAuthenticated(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == getAuthToken() {
			h.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Get name from path
	a := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(a[len(a)-1])

	product, err := productRepository.ByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(product)
		w.Write(resBody)
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resBody, _ := json.Marshal(productRepository.GetProducts())
	w.Write(resBody)
}

func commonMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return WithCorrelationID(IsAuthenticated(f))
}

func GetHTTPHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/product/", commonMiddleware(GetProduct))
	mux.HandleFunc("/products/", commonMiddleware(GetProducts))

	return mux
}
