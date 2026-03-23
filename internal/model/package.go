package model

type Package struct {
	ID         string
	WeightKg   float64
	DistanceKm float64
	OfferCode  string
}

type PackageBundle struct {
	Packages    []Package
	TotalWeight float64
}
