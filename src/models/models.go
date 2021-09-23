package models

type Product struct {
	Id              int
	Name            string
	Price           float64
	DiscountedPrice float64
	Discount        float64
}

type Game struct {
	Genre       string
	ReleaseDate string
	Product
}

type Hardware struct {
	Features []string
	Product
}

type Service struct {
	ServiceType  string
	DeliveryType string
	Product
}

func getDiscountedPrice(p Product) float64 {
	return p.Discount * p.Price
}
