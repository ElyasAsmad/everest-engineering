package shipping

import (
	"fmt"

	"github.com/ElyasAsmad/everestengineering2/internal/model"
)

func GenerateCombinations(packages []model.Package, maxWeight float64) []model.PackageBundle {
	validCombinations := []model.PackageBundle{{}}

	for _, pkg := range packages {
		currentCount := len(validCombinations)

		for i := range currentCount {
			existingCombination := validCombinations[i]

			currentWeight := 0.0
			for _, p := range existingCombination.Packages {
				currentWeight += p.WeightKg
			}

			if currentWeight+pkg.WeightKg <= maxWeight {
				newCombination := append([]model.Package{}, existingCombination.Packages...)
				newCombination = append(newCombination, pkg)
				validCombinations = append(validCombinations, model.PackageBundle{
					Packages:    newCombination,
					TotalWeight: currentWeight + pkg.WeightKg,
				})
			}
		}
	}

	var finalResult []model.PackageBundle
	for _, combination := range validCombinations {
		if len(combination.Packages) > 0 {
			finalResult = append(finalResult, combination)
		}
	}

	return finalResult
}

func GetOptimalShipment(packageCombinations []model.PackageBundle) (*model.PackageBundle, error) {

	if len(packageCombinations) == 0 {
		return nil, fmt.Errorf("no valid combinations generated")
	}

	optimal := (packageCombinations)[0]
	for _, combo := range packageCombinations {
		if combo.TotalWeight > optimal.TotalWeight {
			optimal = combo
		}
	}

	return &optimal, nil
}
