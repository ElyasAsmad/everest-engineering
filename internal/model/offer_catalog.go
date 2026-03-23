package model

import (
	"fmt"

	"github.com/ElyasAsmad/everestengineering2/pkg/expr"
)

type OfferCatalog struct {
	Offers []*Offer
}

func NewOfferCatalog(offers []*Offer) (*OfferCatalog, error) {
	for _, offer := range offers {
		node, err := expr.Compile(offer.Constraint)
		if err != nil {
			return nil, fmt.Errorf("offer %s: invalid constraint %q: %w", offer.Code, offer.Constraint, err)
		}
		offer.compiled = node
	}
	return &OfferCatalog{
		Offers: offers,
	}, nil
}

func (c *OfferCatalog) Apply(offerCode string, distanceKm, weightKg float64) (*Offer, bool, error) {
	ctx := expr.Context{"d": distanceKm, "w": weightKg}

	for _, offer := range c.Offers {
		if offer.Code != offerCode {
			continue
		}

		match, err := expr.Evaluate(offer.compiled, ctx)
		if err != nil {
			return nil, false, fmt.Errorf("failed to evaluate offer %s: %w", offer.Code, err)
		}
		return offer, match, nil
	}

	return nil, false, nil
}
