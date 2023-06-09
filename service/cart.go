package service

import (
	"go-api-meli/dto"
	"go-api-meli/model"
	"go-api-meli/repository"
)

type CartService interface {
	AddProductToCart(cart model.Cart) (model.Cart, error)
	GetCartById(ID uint64) ([]model.CartResponse, error)
	Checkout(IDCart uint64) (dto.Checkout, error)
	CartValidate(cartID uint64) error
	CheckoutValidate(cartID uint64) ([]dto.CheckoutValidate, error)
	AjustCart(cart model.Cart) model.Cart
	DeleteCart(ID uint64) error
}
type cartService struct {
	CartRepository    repository.CartRepository
	ProductRepository repository.ProductRepository
}

func NewCartService(cr repository.CartRepository, pr repository.ProductRepository) cartService {
	return cartService{
		CartRepository:    cr,
		ProductRepository: pr,
	}
}
func (cs cartService) AddProductToCart(cart model.Cart) (model.Cart, error) {
	return cs.CartRepository.AddProductToCart(cart)
}

func (cs cartService) GetCartById(ID uint64) ([]model.CartResponse, error) {
	return cs.CartRepository.GetCartById(ID)
}
func (cs cartService) Checkout(IDCart uint64) (dto.Checkout, error) {
	return cs.CartRepository.Checkout(IDCart)
}

func (cs cartService) CartValidate(cartID uint64) error {
	return cs.CartRepository.CartValidate(cartID)
}
func (cs cartService) CheckoutValidate(cartID uint64) ([]dto.CheckoutValidate, error) {
	return cs.CartRepository.CheckoutValidate(cartID)
}
func (cs cartService) AjustCart(cart model.Cart) model.Cart {
	return cs.CartRepository.AjustCart(cart)
}
func (cs cartService) DeleteCart(IDCart uint64) error {
	return cs.CartRepository.DeleteCart(IDCart)
}
