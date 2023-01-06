package model

import "errors"

type Product struct {
	Name     string `json:"firstName" pact:"example=Tshirt"`
	Category string `json:"lastName" pact:"example=Clothing"`
	Color    string `json:"username" pact:"example=green"`
	ID       int    `json:"id" pact:"example=10"`
}

var (
	ErrNotFound = errors.New("not found")

	ErrUnauthorized = errors.New("unauthorized")

	ErrEmpty = errors.New("empty string")
)

type AddProductRequest struct {
	Name     string `json:"name" pact:"example=Short"`
	Category string `json:"category" pact:"example=Clothing"`
	Color    string `json:"color" pact:"example=blue"`
}

type AddProductResponse struct {
	Name     string `json:"name" pact:"example=Tshirt"`
	Category string `json:"category" pact:"example=Clothing"`
	Color    string `json:"color" pact:"example=green"`
	ID       int    `json:"id" pact:"example=10"`
}
