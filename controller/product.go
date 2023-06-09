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

type ProductController interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProductById(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
}

type productController struct {
	ProductService service.ProductService
}

func NewProductController(service service.ProductService) productController {
	return productController{
		ProductService: service,
	}
}
func (service productController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var product model.Product

	if err = json.Unmarshal(request, &product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := service.ProductService.CreateProduct(product)
	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_erro", "message: There was an error when trying to insert the product.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (service productController) GetProducts(w http.ResponseWriter, r *http.Request) {

	product, err := service.ProductService.GetAll()
	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_erro", "message: Error when searching for products.")
		return
	}
	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)
}

func (service productController) GetProductById(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	IDProduct, err := strconv.ParseUint(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, _ := service.ProductService.ProductValidate(IDProduct)
	if result.ID != IDProduct {
		util.JSON(w, http.StatusNotFound, "code: product_not_found", "message: Product not found.")
		return
	}

	product, err := service.ProductService.GetById(IDProduct)
	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_erro", "message: There was an error when trying to get the product.")
		return
	}
	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)
}

func (service productController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ProductID, err := strconv.ParseInt(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var product model.Product

	if err = json.Unmarshal(request, &product); err != nil {
		w.Write([]byte(fmt.Sprintf("Error converting product object to struct.")))
		return
	}

	result, _ := service.ProductService.ProductValidate(uint64(ProductID))
	if result.ID != uint64(ProductID) {
		util.JSON(w, http.StatusNotFound, "code: product_not_found", "message: Product not found.")
		return
	}

	response, err := service.ProductService.UpdateProduct(uint64(ProductID), product)
	if err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_error.", "message: There was an error when trying to update the product.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (service productController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	productID, err := strconv.ParseUint(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, _ := service.ProductService.ProductValidate(productID)
	if result.ID != productID {
		util.JSON(w, http.StatusNotFound, "code: product_not_found", "message: Product not found.")
		return
	}

	if err = service.ProductService.DeleteProduct(uint64(productID)); err != nil {
		util.JSON(w, http.StatusInternalServerError, "code: internal_server_error", "message: There was an error when trying to delete the product.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
