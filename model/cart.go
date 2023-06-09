package model

type Details struct {
	IDProduct  uint64  `json: "IDProduct"`
	Quantity   int8    `json: "quantity"`
	PriceFinal float64 `json: "PriceFinal"`
}

type Cart struct {
	IDCart   int64     `json:id`
	Products []Details `json: "products"`
}

type Purchase struct {
	ID            int64   `json: "id"`
	QuantityStock int64   `json: "QuantityInStock"`
	QuantityItems int64   `json: "QuantityInItems"`
	Price         float64 `json: "Price"`
}

type CartResponse struct {
	IdCart          uint64  `json: "idtb_cart"`
	Title           string  `json: "title"`
	QuantityOfItems uint64  `json: "QuantityOfItems"`
	PriceFinal      float64 `json: "PriceFinal"`
}
