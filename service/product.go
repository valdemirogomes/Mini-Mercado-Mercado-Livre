package service

import (
	"go-api-meli/model"
	"go-api-meli/repository"
)

type ProductService interface {
	CreateProduct(product model.Product) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetById(ID uint64) (model.Product, error)
	DeleteProduct(ID uint64) error
	UpdateProduct(ID uint64, product model.Product) (model.Product, error)
	ProductValidate(productID uint64) (model.Product, error)
}
type productService struct {
	Repository repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) productService {
	return productService{
		Repository: repo,
	}
}
func (pr productService) CreateProduct(product model.Product) (*model.Product, error) {
	return pr.Repository.CreateProduct(product)
}
func (pr productService) GetAll() ([]model.Product, error) {
	return pr.Repository.GetAll()
}
func (pr productService) DeleteProduct(ID uint64) error {
	return pr.Repository.DeleteProduct(ID)
}
func (pr productService) GetById(ID uint64) (model.Product, error) {
	return pr.Repository.GetById(ID)
}
func (pr productService) UpdateProduct(ID uint64, product model.Product) (model.Product, error) {
	return pr.Repository.UpdateProduct(ID, product)
}
func (pr productService) ProductValidate(productID uint64) (model.Product, error) {
	return pr.Repository.ProductValidate(productID)
}
