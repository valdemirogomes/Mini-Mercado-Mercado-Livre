package routes

import (
	"fmt"
	"go-api-meli/controller"
	"go-api-meli/database"
	"go-api-meli/repository"
	"go-api-meli/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Routers(controllers *Controllers) {

	router := mux.NewRouter()
	router.HandleFunc("/products", controllers.productController.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", controllers.productController.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", controllers.productController.GetProductById).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", controllers.productController.DeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/products/{productID}", controllers.productController.UpdateProduct).Methods(http.MethodPut)

	router.HandleFunc("/carts", controllers.cartController.AddProductToCart).Methods(http.MethodPost)
	router.HandleFunc("/carts/{cartID}", controllers.cartController.GetCartById).Methods(http.MethodGet)
	router.HandleFunc("/carts/{cartId}", controllers.cartController.Checkout).Methods(http.MethodPost)
	//router.HandleFunc("/cart/{codeTbProduct}/{codeTbCart}", controllers.cartController.InsertTbProductTbcart).Methods(http.MethodPost)

	fmt.Println("servidor rodando")
	log.Fatal(http.ListenAndServe(":5000", router))
}

type Controllers struct {
	productController controller.ProductController
	cartController    controller.CartController
}

func BuildControllers() *Controllers {
	db, _ := database.Connection()
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepository, productRepository)
	return &Controllers{
		productController: controller.NewProductController(productService),
		cartController:    controller.NewCartController(cartService, productService),
	}
}
