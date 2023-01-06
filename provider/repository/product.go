package repository

import "github.com/doruk581/cdc-wholesale-workshop-go/model"

type ProductRepository struct {
	Products map[string]*model.Product
}

func (p *ProductRepository) GetProducts() []model.Product {
	var response []model.Product

	for _, product := range p.Products {
		response = append(response, *product)
	}

	return response
}

func (p *ProductRepository) ByName(name string) (*model.Product, error) {
	if product, ok := p.Products[name]; ok {
		return product, nil
	}
	return nil, model.ErrNotFound
}

// ByID finds a product by their ID
func (p *ProductRepository) ByID(ID int) (*model.Product, error) {
	for _, product := range p.Products {
		if product.ID == ID {
			return product, nil
		}
	}
	return nil, model.ErrNotFound
}
