package repository

import (
	"database/sql"
	"go-api-meli/model"
)

type products struct {
	db *sql.DB
}

type ProductRepository interface {
	CreateProduct(product model.Product) (*model.Product, error)
	GetAll() ([]model.Product, error)
	GetById(ID uint64) (model.Product, error)
	DeleteProduct(ID uint64) error
	UpdateProduct(ID uint64, product model.Product) (model.Product, error)
	ProductValidate(productID uint64) (model.Product, error)
}

func NewProductRepository(db *sql.DB) *products {
	return &products{db}
}
func (p products) CreateProduct(product model.Product) (*model.Product, error) {

	statement, err := p.db.Prepare(
		"insert into tb_product (title, price, quantity_in_stock) values (?,?,?)",
	)
	if err != nil {
		return nil, err
	}

	defer statement.Close()

	result, err := statement.Exec(product.Title, product.Price, product.QuantityInStock)
	if err != nil {
		return nil, err
	}
	IDProduct, _ := result.LastInsertId()

	response := model.Product{
		ID:              uint64(IDProduct),
		Title:           product.Title,
		Price:           product.Price,
		QuantityInStock: product.QuantityInStock,
	}
	return &response, nil

}
func (p products) GetAll() ([]model.Product, error) {
	rows, err := p.db.Query("select * from tb_product")
	if err != nil {
		return nil, err

	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product

		if err = rows.Scan(&product.ID, &product.Title, &product.Price, &product.QuantityInStock); err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, nil

}
func (p products) GetById(ID uint64) (model.Product, error) {

	row, err := p.db.Query("select idtb_product, title, price, quantity_in_stock from tb_product where idtb_product = ?", ID)
	if err != nil {
		return model.Product{}, err
	}
	defer row.Close()

	var prd model.Product

	if row.Next() {
		if err = row.Scan(
			&prd.ID,
			&prd.Title,
			&prd.Price,
			&prd.QuantityInStock,
		); err != nil {
			return model.Product{}, err

		}
	}
	return prd, nil
}

func (p products) UpdateProduct(ID uint64, product model.Product) (model.Product, error) {
	statement, err := p.db.Prepare("update tb_product set title = ?, price = ?, quantity_in_stock = ? where idtb_product = ? ")
	if err != nil {
		return model.Product{}, err
	}
	defer statement.Close()

	_, err = statement.Exec(&product.Title, &product.Price, &product.QuantityInStock, ID)
	if err != nil {
		return model.Product{}, err

	}

	response := model.Product{
		ID:              ID,
		Title:           product.Title,
		Price:           product.Price,
		QuantityInStock: product.QuantityInStock,
	}
	return response, nil

}
func (p products) DeleteProduct(ID uint64) error {
	statement, err := p.db.Prepare("delete from tb_product where idtb_product = ? ")
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
func (p products) ProductValidate(productID uint64) (model.Product, error) {
	row, err := p.db.Query("select idtb_product, title, price, quantity_in_stock from tb_product where idtb_product = ?", productID)
	if err != nil {
		return model.Product{}, err
	}
	defer row.Close()

	var prd model.Product

	if row.Next() {
		if err = row.Scan(
			&prd.ID,
			&prd.Title,
			&prd.Price,
			&prd.QuantityInStock,
		); err != nil {
			return model.Product{}, err
		}
	}
	return prd, nil
}
