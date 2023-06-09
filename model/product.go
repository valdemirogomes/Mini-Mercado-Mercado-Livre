package model

type Product struct {
	ID              uint64  `json: "id"`
	Title           string  `json: "title"`
	Price           float64 `json: "price"`
	QuantityInStock int64   `json: "QuantityInStock"`
}
