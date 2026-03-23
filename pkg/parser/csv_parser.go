package parser

import (
	"fmt"
	"os"

	"github.com/ElyasAsmad/everestengineering2/internal/model"
	"github.com/gocarina/gocsv"
)

type ParsedOffer struct {
	Offers []*model.OfferCSV
}

func ParseOffersCSV(fileName string) (*ParsedOffer, error) {
	clientsFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed to open CSV file: %w", err)
	}
	defer clientsFile.Close()

	offers := []*model.OfferCSV{}

	if err := gocsv.UnmarshalFile(clientsFile, &offers); err != nil {
		return nil, fmt.Errorf("Failed to parse CSV file: %w", err)
	}

	return &ParsedOffer{Offers: offers}, nil
}

func (o *ParsedOffer) ConvertToOffer() []*model.Offer {
	var convertedOffers []*model.Offer
	for _, offer := range o.Offers {
		// Perform conversion logic here
		convertedOffers = append(convertedOffers, &model.Offer{
			Code:       offer.Code,
			Discount:   offer.Discount,
			Constraint: offer.Distance + " && " + offer.Weight,
		})
	}
	return convertedOffers
}
