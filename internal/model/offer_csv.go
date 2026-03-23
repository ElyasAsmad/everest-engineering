package model

type OfferCSV struct {
	Code     string  `csv:"code"`
	Discount float64 `csv:"discount"`
	Distance string  `csv:"distance"`
	Weight   string  `csv:"weight"`
}
