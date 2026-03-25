package app

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/ElyasAsmad/everestengineering2/internal/logger"
	"github.com/ElyasAsmad/everestengineering2/internal/model"
	"github.com/ElyasAsmad/everestengineering2/internal/parser"
	"github.com/ElyasAsmad/everestengineering2/internal/shipping"
	"github.com/shopspring/decimal"
)

func Run(in io.Reader, inputFile string) (string, error) {
	logger := logger.NewLogger()

	offers, err := parser.ParseOffersCSV(inputFile)
	if err != nil {
		return "", fmt.Errorf("Failed to parse CSV file: %v", err)
	}

	var parsedOffers strings.Builder
	for _, offer := range offers {
		fmt.Fprintf(&parsedOffers, "\tOffer Code: %s, Discount: %.2f%%, Constraint: %s\n",
			offer.Code, offer.Discount, offer.Constraint)
	}
	logger.Debugf("Parsed Offers:\n%s", parsedOffers.String())

	catalog, err := model.NewOfferCatalog(offers)
	if err != nil {
		return "", fmt.Errorf("Failed to create offer catalog: %v", err)
	}

	// 1st scan: base cost, number of packages
	var baseCost decimal.Decimal
	var noOfPackages int

	// 2nd scan: no of vehicles, max speed, max load
	var noOfVehicles int
	var maxSpeed float64
	var maxLoad float64

	scanner := bufio.NewScanner(in)

	if scanner.Scan() {
		line := scanner.Text()

		// read base cost, no of packages
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid input format. Expected: <baseCost> <noOfPackages>")
		}

		parsedBaseCost, err := decimal.NewFromString(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid input format. Expected: <baseCost> <noOfPackages>")
		}

		parsedNoOfPackages, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid input format. Expected: <baseCost> <noOfPackages>")
		}

		baseCost = parsedBaseCost
		noOfPackages = parsedNoOfPackages
	} else {
		return "", fmt.Errorf("failed to read base cost and number of packages")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading base cost and number of packages: %v", err)
	}

	logger.Debugf("Base Cost: %s", baseCost.String())
	logger.Debugf("Number of Packages: %d", noOfPackages)

	packages := make([]model.Package, noOfPackages)

	for i := 0; i < noOfPackages; i++ {

		if scanner.Scan() {
			line := scanner.Text()

			var id string
			var weight float64
			var distance float64
			var offerCode string

			// package details (id, weight, distance, offer code)
			n, err := fmt.Sscanf(line, "%s %f %f %s", &id, &weight, &distance, &offerCode)

			if err != nil || n != 4 {
				return "", fmt.Errorf("invalid package details format. Expected: <id> <weight> <distance> <offerCode>")
			}

			logger.Debugf("Package ID: %s, Weight: %f, Distance: %f, Offer Code: %s", id, weight, distance, offerCode)

			packages[i] = model.Package{
				ID:         id,
				WeightKg:   weight,
				DistanceKm: distance,
				OfferCode:  offerCode,
			}
		} else {
			return "", fmt.Errorf("failed to read package details")
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading package details: %v", err)
		}
	}

	if scanner.Scan() {
		line := scanner.Text()

		n, err := fmt.Sscanf(line, "%d %f %f", &noOfVehicles, &maxSpeed, &maxLoad)

		if err != nil || n != 3 {
			return "", fmt.Errorf("invalid vehicle details format. Expected: <noOfVehicles> <maxSpeed> <maxLoad>")
		}

		logger.Debugf("Number of Vehicles: %d, Max Speed: %f, Max Load: %f", noOfVehicles, maxSpeed, maxLoad)
	} else {
		return "", fmt.Errorf("failed to read vehicle details")
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading vehicle details: %v", err)
	}

	// fleet manager
	shipper := shipping.NewShipper(noOfVehicles, maxSpeed)
	deliveryResult := make([]model.DeliveryResult, 0)

	// start loop
	for len(packages) > 0 {
		logger.Debugf("---- %d packages remaining ---", len(packages))

		// generate combinations
		combinations := shipping.GenerateCombinations(packages, maxLoad)

		// get optimal shipment from generated combinations
		op, err := shipping.GetOptimalShipment(combinations)

		if err != nil {
			return "", fmt.Errorf("error getting optimal shipment: %v", err)
		}

		logger.Debugf("Optimal Shipment: %v, Total Weight: %f", op.Packages, op.TotalWeight)

		result := shipper.ProcessShipment(op)

		deliveryResult = append(deliveryResult, *result...)

		// remove shipped packages from original list
		packages = shipping.FilterOutPackages(packages, op.Packages)

		remainingIDs := make([]string, len(packages))
		for i, pkg := range packages {
			remainingIDs[i] = pkg.ID
		}
		logger.Debugf("Remaining Package IDs: %v", remainingIDs)
	}

	var output strings.Builder

	// sort delivery results by package ID
	slices.SortFunc(deliveryResult, func(a, b model.DeliveryResult) int {
		return cmp.Compare(a.Package.ID, b.Package.ID)
	})

	for _, res := range deliveryResult {
		pkg := res.Package

		offer, qualifies, err := catalog.Apply(pkg.OfferCode, pkg.DistanceKm, pkg.WeightKg)
		if err != nil {
			return "", fmt.Errorf("error applying offer %s to package %s: %v", pkg.OfferCode, pkg.ID, err)
		}

		// by default, no discount
		discount := decimal.Zero
		deliveryCost := baseCost.
			Add(decimal.NewFromFloat(pkg.WeightKg).Mul(decimal.NewFromInt(10))).
			Add(decimal.NewFromFloat(pkg.DistanceKm).Mul(decimal.NewFromInt(5)))

		// if qualifies for discount, then calculate discount and final delivery cost
		if qualifies {
			discount = decimal.NewFromFloat(offer.Discount).
				Div(decimal.NewFromInt(100)).
				Mul(deliveryCost)
			deliveryCost = deliveryCost.Sub(discount)
		}

		fmt.Fprintf(
			&output,
			"%s %s %s %.2f\n",
			pkg.ID,
			discount.Round(0).StringFixed(0),
			deliveryCost.Round(0).StringFixed(0),
			res.DeliveryTime,
		)
	}

	return output.String(), nil
}
