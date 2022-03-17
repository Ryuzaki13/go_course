package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func prepareProduct() []string {
	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	errorList := make([]string, 0)
	var e error

	query["ProductSearch"], e = Link.Prepare(`SELECT id, "name", price FROM product WHERE "name" ILIKE '%' || $1 || '%' ORDER BY price`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	query["ProductDelete"], e = Link.Prepare(`DELETE FROM product WHERE id = $1`)
	if e != nil {
		errorList = append(errorList, e.Error())
	}

	return errorList
}

func SearchProduct(name string) []Product {
	q, ok := query["ProductSearch"]
	if !ok {
		fmt.Println("ОШИБКА")
		return nil
	}

	rows, e := q.Query(name)
	if e != nil {
		fmt.Println("ОШИБКА", e)
		return nil
	}

	defer rows.Close()

	p := Product{}
	products := make([]Product, 0)

	for rows.Next() {
		e = rows.Scan(&p.ID, &p.Name, &p.Price)
		if e != nil {
			fmt.Println(e)
			return nil
		}

		products = append(products, Product{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	return products
}

func DeleteProduct(id int) error {
	q, ok := query["ProductDelete"]
	if !ok {
		fmt.Println("ОШИБКА")
		return errors.New("ошощывавыа ыва")
	}

	_, e := q.Exec(id)
	if e != nil {
		fmt.Println("ОШИБКА", e)
		return e
	}

	return nil
}
