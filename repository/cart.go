package repository

import (
	"database/sql"
	"errors"
	_ "errors"
	"go-api-meli/dto"
	"go-api-meli/model"
	"go-api-meli/util"
	"time"
)

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(cart model.Cart) (model.Cart, error)
	GetCartById(ID uint64) ([]model.CartResponse, error)
	Checkout(IDCart uint64) (dto.Checkout, error)
	UpdateInventoryColumn(quantity int64, idCart uint64) error
	SubtractOfItems(IDCart uint64) (float64, int, error)
	CartValidate(cartID uint64) error
	UpdateColumnPriceFinal(ID uint64, ValueFinal float64) (float64, error)
	CheckoutValidate(cartID uint64) ([]dto.CheckoutValidate, error)
	AjustCart(cart model.Cart) model.Cart
	DeleteCart(ID uint64) error
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

func (c cart) AddProductToCart(cart model.Cart) (model.Cart, error) {
	var IDCart, _ = c.CreateCart()
	detail := []model.Details{}

	for _, products := range cart.Products {

		if products.Quantity > 0 {
			statement, err := c.db.Prepare("insert into tb_cart_tb_product (codetb_product, codetb_cart, quantity_of_items) values (?,?,?)")

			if err != nil {
				return model.Cart{}, err
			}

			defer statement.Close()

			statement.Exec(products.IDProduct, IDCart, products.Quantity)
			detail = append(detail, model.Details{IDProduct: products.IDProduct, Quantity: products.Quantity})
		}
	}
	cartResponse := model.Cart{
		IDCart:   IDCart,
		Products: detail,
	}
	return cartResponse, nil
}

func (c cart) GetCartById(ID uint64) ([]model.CartResponse, error) {

	row, err := c.db.Query(`
	select 
    c.idtb_cart, 
    p.idtb_product,
	cp.quantity_of_items
    from tb_cart c
    join tb_cart_tb_product cp
    on c.idtb_cart = cp.codetb_cart
    join tb_product p
    on p.idtb_product=cp.codetb_product
    where  c.idtb_cart = ?`, ID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var cart []model.CartResponse

	for row.Next() {
		var c model.CartResponse
		if err = row.Scan(
			&c.IdCart,
			&c.QuantityOfItems,
			&c.Title,
		); err != nil {
			return nil, err
		}
		cart = append(cart, c)
	}
	return cart, nil
}

func (c cart) Checkout(IDCart uint64) (dto.Checkout, error) {

	result, _, err := c.SubtractOfItems(IDCart)

	if err != nil {
		return dto.Checkout{}, err
	}

	return dto.Checkout{Price: result}, nil
}

func (c cart) SubtractOfItems(IDCart uint64) (float64, int, error) {

	rows, err := c.db.Query(`select
	p.idtb_product,
	p.quantity_in_stock, 
	cp.quantity_of_items,
	p.price
	from
	tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	where c.idtb_cart = ?`, IDCart)
	if err != nil {
		return 0, 0, err
	}

	defer rows.Close()

	var stockFinal, priceFinal, result float64
	var product model.Purchase
	for rows.Next() {
		if err = rows.Scan(&product.ID, &product.QuantityStock, &product.QuantityItems, &product.Price); err != nil {
			return 0, 0, err
		}
		priceFinal = product.QuantityItems * product.Price
		result += priceFinal
		stockFinal = product.QuantityStock - product.QuantityItems
		c.UpdateInventoryColumn(int64(stockFinal), uint64(product.ID))
	}
	price, _ := c.UpdateColumnPriceFinal(IDCart, result)
	return price, int(product.QuantityItems), nil
}

func (c cart) UpdateInventoryColumn(quantity int64, idCart uint64) error {

	statement, err := c.db.Prepare(
		`update tb_cart c
			join tb_cart_tb_product cp
			on cp.codetb_cart = c.idtb_cart
			join tb_product p
			on p.idtb_product = cp.codetb_product
			set p.quantity_in_stock = ?
			where p.idtb_product = ? `)
	if err != nil {
		return err
	}
	defer statement.Close()
	statement.Exec(quantity, idCart)
	return nil
}
func (c cart) UpdateColumnPriceFinal(ID uint64, ValueFinal float64) (float64, error) {
	statement, err := c.db.Prepare("insert into tb_price_final (idtb_cart, price_final) values (?,?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	_, err = statement.Exec(ID, ValueFinal)
	if err != nil {
		return 0, err
	}
	return float64(ValueFinal), nil
}

func (c cart) CreateCart() (int64, error) {
	date := time.Now()
	statement, err := c.db.Prepare(
		"insert into tb_cart (date) values (?)",
	)
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(date)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return ID, err
}

func (c cart) CartValidate(cartID uint64) error {
	row, err := c.db.Query(`select 
	count(*)
    from tb_cart c
    join tb_cart_tb_product cp
    on c.idtb_cart = cp.codetb_cart
    join tb_product p
    on p.idtb_product=cp.codetb_product
    where  c.idtb_cart = ?`, cartID)
	if err != nil {
		return err
	}
	defer row.Close()

	var count int
	if row.Next() {
		if err = row.Scan(
			&count,
		); err != nil {
			return err
		}
	}
	if count == 0 {
		return errors.New("carrinho nao encontrado")
	}
	return nil
}

func (c cart) CheckoutValidate(cartID uint64) ([]dto.CheckoutValidate, error) {
	rows, err := c.db.Query(`select 
	c.idtb_cart, 
	p.idtb_product,
	cp.quantity_of_items
    from tb_cart c
    join tb_cart_tb_product cp
    on c.idtb_cart = cp.codetb_cart
    join tb_product p
    on p.idtb_product=cp.codetb_product
    where  c.idtb_cart = ?`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []dto.CheckoutValidate

	for rows.Next() {
		var cart dto.CheckoutValidate

		if err = rows.Scan(
			&cart.IDCart,
			&cart.IDProduct,
			&cart.QuantityOfItems,
		); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func (c cart) AjustCart(cart model.Cart) model.Cart {
	var newProduct []model.Details
	indices := []int{}
	for i, products := range cart.Products {
		if !util.Contains(indices, i) {

			for j, prod := range cart.Products {
				if (i != j) && (prod.IDProduct == products.IDProduct) {
					products.Quantity += prod.Quantity
					indices = append(indices, j)
				}
			}
			newProduct = append(newProduct, products)
		}
	}
	newCart := model.Cart{IDCart: cart.IDCart, Products: newProduct}
	return newCart
}

func (c cart) DeleteCart(ID uint64) error {
	statement, err := c.db.Prepare("delete from tb_cart_tb_product where codetb_cart = ? ")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(ID)
	if err != nil {
		return err
	}
	return nil
}
