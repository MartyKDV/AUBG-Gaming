package models

import "database/sql"

func getDiscountedPrice(p Product) float64 {
	return p.Price - (p.Discount * p.Price)
}

type Product struct {
	Id              int
	Name            string
	Price           float64
	DiscountedPrice float64
	Discount        float64
}

type Collection struct {
	Games    []Game
	Hardware []Hardware
	Services []Service
}

type Item struct {
	Product
	Game
	Hardware
	Category string
}

type Game struct {
	Genre       sql.NullString
	ReleaseDate sql.NullString
}

type Hardware struct {
	Features sql.NullString
}

type Service struct {
	ServiceType  sql.NullString
	DeliveryType sql.NullString
}
