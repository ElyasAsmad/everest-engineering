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

func ParseOffersCSV(fileName string) ([]*model.Offer, error) {
	clientsFile, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer clientsFile.Close()

	var raw []*model.OfferCSV
	if err := gocsv.UnmarshalFile(clientsFile, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse CSV file: %w", err)
	}

	offers := make([]*model.Offer, len(raw))
	for i, o := range raw {
		offers[i] = &model.Offer{
			Code:       o.Code,
			Discount:   o.Discount,
			Constraint: o.Distance + " && " + o.Weight,
		}
	}

	return offers, nil
}
