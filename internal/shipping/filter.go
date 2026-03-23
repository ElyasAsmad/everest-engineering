package shipping

import "github.com/ElyasAsmad/everestengineering2/internal/model"

func FilterOutPackages(packages []model.Package, shipped []model.Package) []model.Package {
	shippedIDs := make(map[string]struct{}, len(shipped))
	for _, p := range shipped {
		shippedIDs[p.ID] = struct{}{}
	}

	result := make([]model.Package, 0, len(packages))
	for _, p := range packages {
		if _, ok := shippedIDs[p.ID]; !ok {
			result = append(result, p)
		}
	}
	return result
}
