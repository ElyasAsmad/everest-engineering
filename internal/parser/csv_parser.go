package parser

import (
	"fmt"
	"os"
	"strings"

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

	if err := validateOfferCSV(raw); err != nil {
		return nil, fmt.Errorf("invalid offer data: %w", err)
	}

	offers := make([]*model.Offer, len(raw))

	for i, o := range raw {
		offers[i] = &model.Offer{
			Code:       o.Code,
			Discount:   o.Discount,
			Constraint: buildConstraint(o.Distance, o.Weight),
		}
	}

	return offers, nil
}

func validateOfferCSV(offers []*model.OfferCSV) error {
	for _, offer := range offers {
		if offer.Code == "" {
			return fmt.Errorf("offer code is required")
		}
		if offer.Discount < 0 || offer.Discount > 100 {
			return fmt.Errorf("discount must be between 0 and 100")
		}
	}

	return nil
}

func buildConstraint(distance, weight string) string {
	distance = strings.TrimSpace(distance)
	weight = strings.TrimSpace(weight)

	switch {
	case distance != "" && weight != "":
		return fmt.Sprintf("%s && %s", distance, weight)
	case distance != "":
		return distance
	case weight != "":
		return weight
	default:
		return ""
	}
}
