package shipping

import (
	"fmt"

	"github.com/ElyasAsmad/everestengineering2/internal/logger"
	"github.com/ElyasAsmad/everestengineering2/internal/model"
	f "github.com/ElyasAsmad/everestengineering2/internal/utils"
)

type Vehicle struct {
	ID            int
	AvailableTime float64
}

type Shipper struct {
	fleet    []Vehicle
	maxSpeed float64
}

func NewShipper(noOfVehicles int, maxSpeed float64) *Shipper {
	// Initialize vehicles
	fleet := make([]Vehicle, noOfVehicles)
	for i := range noOfVehicles {
		fleet[i] = Vehicle{
			ID:            i + 1,
			AvailableTime: 0.0,
		}
	}

	return &Shipper{
		fleet:    fleet,
		maxSpeed: maxSpeed,
	}
}

func (s *Shipper) FastForwardAndGetVehicle() (*Vehicle, error) {
	if len(s.fleet) == 0 {
		return nil, fmt.Errorf("no vehicles in the fleet")
	}

	earliestTime := s.fleet[0].AvailableTime
	for _, vehicle := range s.fleet {
		if vehicle.AvailableTime < earliestTime {
			earliestTime = vehicle.AvailableTime
		}
	}

	for i := range s.fleet {
		if s.fleet[i].AvailableTime == earliestTime {
			return &s.fleet[i], nil
		}
	}

	// only reachable by miracle: e.g.: cosmic ray flips the bit
	panic("unreachable code: should have returned a vehicle by now")
}

func (s *Shipper) ProcessShipment(shipment *model.PackageBundle) *[]model.DeliveryResult {
	logger := logger.NewLogger()

	// find the next vehicle (earliest available)
	vehicle, err := s.FastForwardAndGetVehicle()
	if err != nil {
		logger.Error("Error getting next vehicle:", err)
		return nil
	}
	longestDistance := 0.0

	logger.Debugf("Shipment assigned to Vehicle %d: Available Time %.2f", vehicle.ID, vehicle.AvailableTime)

	result := make([]model.DeliveryResult, len(shipment.Packages))

	for i, pkg := range shipment.Packages {
		deliveryTime := vehicle.AvailableTime + f.Truncate((pkg.DistanceKm/s.maxSpeed), 2)

		logger.Debugf("Delivering %s: (%.2f + %.2f) %.2f hrs", pkg.ID, vehicle.AvailableTime, f.Truncate((pkg.DistanceKm/s.maxSpeed), 2), deliveryTime)

		result[i] = model.DeliveryResult{
			Package:      pkg,
			DeliveryTime: deliveryTime,
		}

		if pkg.DistanceKm > longestDistance {
			longestDistance = f.Truncate(pkg.DistanceKm, 2)
		}
	}

	// use max distance for 1 complete go-return trip time calculation, truncated to 2 d.p.
	roundTripTime := f.Truncate(longestDistance/s.maxSpeed, 2) * 2

	vehicle.AvailableTime = vehicle.AvailableTime + roundTripTime

	logger.Debugf("Vehicle %d will be available after %.2f hrs", vehicle.ID, vehicle.AvailableTime)

	return &result
}
