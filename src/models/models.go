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
	Service
	Category string
}

type Game struct {
	Genre       sql.NullString
	ReleaseDate sql.NullString
}

type Hardware struct {
	Features     sql.NullString
	HardwareType sql.NullString
}

type Service struct {
	ServiceType sql.NullString
}

type Cart struct {
	CartItems []CartItem
}

type CartItem struct {
	ItemID   int
	Quantity int
}

type CartItemDetails struct {
	ID       int
	Name     string
	Quantity int
}
