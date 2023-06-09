package dto

type Checkout struct {
	Price float64 `json: "Price"`
}
type CheckoutValidate struct {
	IDCart          int     `json: "IDCart"`
	IDProduct       int     `json: "IDProduct"`
	Price           float64 `json: "price"`
	QuantityOfItems int     `json: "quantity_of_items"`
}
