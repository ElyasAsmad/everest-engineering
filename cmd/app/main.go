package main

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/ElyasAsmad/everestengineering2/internal/logger"
	"github.com/ElyasAsmad/everestengineering2/internal/model"
	"github.com/ElyasAsmad/everestengineering2/internal/parser"
	"github.com/ElyasAsmad/everestengineering2/internal/shipping"
)

func main() {
	logger := logger.NewLogger()

	// TODO: maybe make the file name passable via arg
	offers, err := parser.ParseOffersCSV("offers.csv")
	if err != nil {
		logger.Error("Failed to parse CSV file: %v", err)
		os.Exit(1)
	}

	logger.Debugf("Parsed Offers: %v", offers)

	catalog, err := model.NewOfferCatalog(offers)
	if err != nil {
		logger.Error("Failed to create offer catalog: %v", err)
		os.Exit(1)
	}

	// 1st scan: base cost, number of packages
	var baseCost float64
	var noOfPackages int

	// 2nd scan: no of vehicles, max speed, max load
	var noOfVehicles int
	var maxSpeed float64
	var maxLoad float64

	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		line := scanner.Text()

		// read base cost, no of packages
		n, err := fmt.Sscanf(line, "%f %d", &baseCost, &noOfPackages)

		if err != nil || n != 2 {
			logger.Error("Invalid input format. Expected: <baseCost> <noOfPackages>")
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error reading base cost and number of packages:", err)
		os.Exit(1)
	}

	logger.Debugf("Base Cost: %f", baseCost)
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
				logger.Error("Invalid package details format. Expected: <id> <weight> <distance> <offerCode>")
				os.Exit(1)
			}

			logger.Debugf("Package ID: %s, Weight: %f, Distance: %f, Offer Code: %s", id, weight, distance, offerCode)

			packages[i] = model.Package{
				ID:         id,
				WeightKg:   weight,
				DistanceKm: distance,
				OfferCode:  offerCode,
			}
		}

		if err := scanner.Err(); err != nil {
			logger.Error("Error reading package details:", err)
			os.Exit(1)
		}
	}

	if scanner.Scan() {
		line := scanner.Text()

		n, err := fmt.Sscanf(line, "%d %f %f", &noOfVehicles, &maxSpeed, &maxLoad)

		if err != nil || n != 3 {
			logger.Error("Invalid vehicle details format. Expected: <noOfVehicles> <maxSpeed> <maxLoad>")
			os.Exit(1)
		}

		logger.Debugf("Number of Vehicles: %d, Max Speed: %f, Max Load: %f", noOfVehicles, maxSpeed, maxLoad)
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error reading vehicle details:", err)
		os.Exit(1)
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
			logger.Error("Error getting optimal shipment:", err)
			os.Exit(1)
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

	// sort delivery results by package ID
	slices.SortFunc(deliveryResult, func(a, b model.DeliveryResult) int {
		return cmp.Compare(a.Package.ID, b.Package.ID)
	})

	for _, res := range deliveryResult {
		pkg := res.Package

		offer, qualifies, err := catalog.Apply(pkg.OfferCode, pkg.DistanceKm, pkg.WeightKg)
		if err != nil {
			logger.Fatalf("Error applying offer %s to package %s: %v", pkg.OfferCode, pkg.ID, err)
		}

		discount := 0.0
		deliveryCost := baseCost + (pkg.WeightKg * 10) + (pkg.DistanceKm * 5)

		if qualifies {
			discount = (offer.Discount / 100.0) * deliveryCost
			deliveryCost = deliveryCost - discount
		}

		logger.Printf("%s %.0f %.0f %.2f", pkg.ID, discount, deliveryCost, res.DeliveryTime)
	}

}
