package shipping

import (
	"testing"

	"github.com/ElyasAsmad/everestengineering2/internal/model"
)

// Tests for GenerateCombinations function

func TestGenerateCombinations_SinglePackageFits(t *testing.T) {
	packages := []model.Package{
		{ID: "PKG1", WeightKg: 5},
	}
	result := GenerateCombinations(packages, 10)
	if len(result) != 1 {
		t.Errorf("Expected 1 combination, got %d", len(result))
	}
}

func TestGenerateCombinations_SinglePackageExceedsLimit(t *testing.T) {
	packages := []model.Package{
		{ID: "PKG1", WeightKg: 15},
	}
	result := GenerateCombinations(packages, 10)
	if len(result) != 0 {
		t.Errorf("Expected 0 combinations, got %d", len(result))
	}
}

func TestGenerateCombinations_ExactWeightLimit(t *testing.T) {
	packages := []model.Package{
		{ID: "PKG1", WeightKg: 5},
		{ID: "PKG2", WeightKg: 5},
	}
	result := GenerateCombinations(packages, 10)
	// should produce: [PKG1], [PKG2], [PKG1, PKG2]
	found := false
	for _, combo := range result {
		if combo.TotalWeight == 10 {
			found = true
		}
	}
	if !found {
		t.Error("Expected a combination with total weight exactly 10")
	}
}

func TestGenerateCombinations_EmptyPackages(t *testing.T) {
	result := GenerateCombinations([]model.Package{}, 10)
	if len(result) != 0 {
		t.Errorf("Expected 0 combinations, got %d", len(result))
	}
}

func TestGenerateCombinations_TotalWeightIsCorrect(t *testing.T) {
	packages := []model.Package{
		{ID: "PKG1", WeightKg: 3},
		{ID: "PKG2", WeightKg: 4},
	}
	result := GenerateCombinations(packages, 10)
	for _, combo := range result {
		actual := 0.0
		for _, pkg := range combo.Packages {
			actual += pkg.WeightKg
		}
		if actual != combo.TotalWeight {
			t.Errorf("Expected total weight %.2f, got %.2f", combo.TotalWeight, actual)
		}
	}
}

// Tests for GetOptimalShipment function

func TestGetOptimalShipment_ReturnsHeaviest(t *testing.T) {
	combinations := []model.PackageBundle{
		{Packages: []model.Package{{ID: "PKG1"}}, TotalWeight: 5},
		{Packages: []model.Package{{ID: "PKG2"}}, TotalWeight: 10},
		{Packages: []model.Package{{ID: "PKG3"}}, TotalWeight: 3},
	}
	result, err := GetOptimalShipment(combinations)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.TotalWeight != 10 {
		t.Errorf("Expected heaviest bundle (10), got %.2f", result.TotalWeight)
	}
}

func TestGetOptimalShipment_EmptyCombinations(t *testing.T) {
	_, err := GetOptimalShipment([]model.PackageBundle{})
	if err == nil {
		t.Error("Expected error for empty combinations, got nil")
	}
}

func TestGetOptimalShipment_SingleCombination(t *testing.T) {
	combinations := []model.PackageBundle{
		{Packages: []model.Package{{ID: "PKG1"}}, TotalWeight: 7},
	}
	result, err := GetOptimalShipment(combinations)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result.TotalWeight != 7 {
		t.Errorf("Expected total weight 7, got %.2f", result.TotalWeight)
	}
}

// Tests for FilterOutPackages function

func TestFilterOutPackages_RemovesShipped(t *testing.T) {
	all := []model.Package{
		{ID: "PKG1"},
		{ID: "PKG2"},
		{ID: "PKG3"},
	}
	shipped := []model.Package{
		{ID: "PKG2"},
	}
	result := FilterOutPackages(all, shipped)
	if len(result) != 2 {
		t.Errorf("Expected 2 packages after filtering, got %d", len(result))
	}
	for _, pkg := range result {
		if pkg.ID == "PKG2" {
			t.Error("Expected PKG2 to be filtered out, but it was found in the result")
		}
	}
}

func TestFilterOutPackages_NoShipped(t *testing.T) {
	all := []model.Package{
		{ID: "PKG1"},
		{ID: "PKG2"},
	}
	result := FilterOutPackages(all, []model.Package{})
	if len(result) != 2 {
		t.Errorf("Expected 2 packages when no shipped, got %d", len(result))
	}
}

func TestFilterOutPackages_AllShipped(t *testing.T) {
	all := []model.Package{
		{ID: "PKG1"},
		{ID: "PKG2"},
	}
	shipped := []model.Package{
		{ID: "PKG1"},
		{ID: "PKG2"},
	}
	result := FilterOutPackages(all, shipped)
	if len(result) != 0 {
		t.Errorf("Expected 0 packages when all shipped, got %d", len(result))
	}
}

// Tests for FastForwardAndGetVehicle function

func TestFastForwardAndGetVehicle_ReturnsEarliestVehicle(t *testing.T) {
	s := &Shipper{
		fleet: []Vehicle{
			{ID: 1, AvailableTime: 5.0},
			{ID: 2, AvailableTime: 2.0},
			{ID: 3, AvailableTime: 8.0},
		},
	}
	v, err := s.FastForwardAndGetVehicle()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if v.ID != 2 {
		t.Errorf("Expected vehicle with ID 2 (earliest), got ID %d", v.ID)
	}
}

func TestFastForwardAndGetVehicle_AllSameTime(t *testing.T) {
	s := &Shipper{
		fleet: []Vehicle{
			{ID: 1, AvailableTime: 3.0},
			{ID: 2, AvailableTime: 3.0},
		},
	}
	v, err := s.FastForwardAndGetVehicle()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// should return the first one in case of tie
	if v.ID != 1 {
		t.Errorf("Expected vehicle with ID 1 (first in tie), got ID %d", v.ID)
	}
}

// Tests for ProcessShipment function

func TestProcessShipment_DeliveryTimeCalculation(t *testing.T) {
	// 1 vehicle, 70 km/h max speed
	s := NewShipper(1, 70)
	shipment := &model.PackageBundle{
		Packages: []model.Package{
			{ID: "PKG1", WeightKg: 5, DistanceKm: 70}, // should take 1 hour at max speed
		},
	}
	results := s.ProcessShipment(shipment)
	if results == nil {
		t.Error("Expected results, got nil")
	}
	if len(*results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(*results))
	}
	if (*results)[0].DeliveryTime != 1.0 {
		t.Errorf("Expected delivery time 1.0 hours, got %.2f", (*results)[0].DeliveryTime)
	}
}

func TestProcessShipment_VehicleAvailableTimeUpdated(t *testing.T) {
	s := NewShipper(1, 70)
	shipment := &model.PackageBundle{
		Packages: []model.Package{
			{ID: "PKG1", WeightKg: 5, DistanceKm: 70}, // 1 hour (2 hours round trip)
		},
	}
	s.ProcessShipment(shipment)
	if s.fleet[0].AvailableTime != 2.0 {
		t.Errorf("Expected vehicle available time to be 2.0 hours after round trip, got %.2f", s.fleet[0].AvailableTime)
	}
}

func TestProcessShipment_MultiplePackagesUseLongestDistanceForRoundTrip(t *testing.T) {
	s := NewShipper(1, 100)
	shipment := &model.PackageBundle{
		Packages: []model.Package{
			{ID: "PKG1", WeightKg: 5, DistanceKm: 100}, // 1 hour
			{ID: "PKG2", WeightKg: 5, DistanceKm: 200}, // 2 hours
		},
		TotalWeight: 10,
	}
	s.ProcessShipment(shipment)
	if s.fleet[0].AvailableTime != 4.0 {
		t.Errorf("Expected vehicle available time to be 4.0 hours after round trip, got %.2f", s.fleet[0].AvailableTime)
	}
}
