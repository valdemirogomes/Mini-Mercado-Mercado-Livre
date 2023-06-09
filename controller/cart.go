package controller

import (
	"encoding/json"
	"fmt"
	"go-api-meli/model"
	"go-api-meli/service"
	"go-api-meli/util"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CartController interface {
	AddProductToCart(w http.ResponseWriter, r *http.Request)
	GetCartById(w http.ResponseWriter, r *http.Request)
	Checkout(w http.ResponseWriter, r *http.Request)
}

type cartController struct {
	CartService    service.CartService
	ProductService service.ProductService
}

func NewCartController(service service.CartService, cartService service.ProductService) cartController {
	return cartController{
		CartService:    service,
		ProductService: cartService,
	}
}

func (controller cartController) AddProductToCart(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var cart model.Cart
	if err = json.Unmarshal(request, &cart); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cart = controller.CartService.AjustCart(cart)

	for _, product := range cart.Products {

		result, _ := controller.ProductService.ProductValidate(product.IDProduct)
		if product.IDProduct < 2 && result.ID != product.IDProduct {

		}
		if result.ID != product.IDProduct {
			util.JSON(w, http.StatusNotFound, "code: product_not_found.", "message: One of the cart products was not found.")
			return
		}
		if result.QuantityInStock < int64(product.Quantity) {
			util.JSON(w, http.StatusUnprocessableEntity, "code: not_enough_stock.", "message: One of the cart products does not have sufficient stock.")
			return
		}
	}

	result, err := controller.CartService.AddProductToCart(cart)
	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_error.", "message: There was an error when trying to create a shopping cart.")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("There was an error when trying to create a shopping cart.")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)

}

func (service cartController) GetCartById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, _ := strconv.ParseUint(parameters["cartID"], 10, 32)

	if err := service.CartService.CartValidate(ID); err != nil {
		util.JSON(w, http.StatusNotFound, "code: cart_not_found.", "message: Shopping cart  was not found.")
		return
	}

	cart, err := service.CartService.GetCartById(ID)
	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_error.", "message: There was an error when trying to get the shopping cart.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
	w.WriteHeader(http.StatusOK)
	return
}

func (service cartController) Checkout(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	IDCart, _ := strconv.ParseUint(paramters["cartId"], 10, 32)

	if err := service.CartService.CartValidate(IDCart); err != nil {
		util.JSON(w, http.StatusNotFound, "code: cart_not_found.", "message: Shopping cart  was not found.")
		return
	}

	items, _ := service.CartService.CheckoutValidate(IDCart)

	for _, item := range items {
		result, _ := service.ProductService.ProductValidate(uint64(item.IDProduct))
		if result.QuantityInStock < int64(item.QuantityOfItems) {
			util.JSON(w, http.StatusBadRequest, "code: not_enough_stock.", "message: A cart product does not have enough stock. This cart is invalid.")
			return
		}
	}

	response, err := service.CartService.Checkout(IDCart)

	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_error.", "message: There was an error when trying to checkout.")
		return
	}
	service.CartService.DeleteCart(IDCart)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
