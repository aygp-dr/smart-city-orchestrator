package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

// Severity represents the operational severity level of an urban system.
type Severity int

const (
	SeverityNormal   Severity = iota
	SeverityWarning
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "Warning"
	case SeverityCritical:
		return "Critical"
	default:
		return "Normal"
	}
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// Metric holds a single named measurement for an urban system.
type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UrbanSystem represents one monitored city subsystem.
type UrbanSystem struct {
	Name     string   `json:"name"`
	Severity Severity `json:"severity"`
	Metrics  []Metric `json:"metrics"`
}

// GenerateSystems produces mock data for all 6 urban systems.
func GenerateSystems(rng *rand.Rand) []UrbanSystem {
	return []UrbanSystem{
		genTraffic(rng),
		genEnergy(rng),
		genWater(rng),
		genWaste(rng),
		genEmergency(rng),
		genAirQuality(rng),
	}
}

func genTraffic(rng *rand.Rand) UrbanSystem {
	congestion := rng.Intn(101)
	statuses := []string{"Normal", "Degraded", "Critical"}
	idx := 0
	sev := SeverityNormal
	if congestion > 70 {
		sev = SeverityCritical
		idx = 2
	} else if congestion > 40 {
		sev = SeverityWarning
		idx = 1
	}
	return UrbanSystem{
		Name:     "Traffic",
		Severity: sev,
		Metrics: []Metric{
			{"Congestion Level", fmt.Sprintf("%d%%", congestion)},
			{"Signal Status", statuses[idx]},
		},
	}
}

func genEnergy(rng *rand.Rand) UrbanSystem {
	load := rng.Intn(101)
	renewable := 20 + rng.Intn(51)
	sev := SeverityNormal
	if load > 85 {
		sev = SeverityCritical
	} else if load > 65 {
		sev = SeverityWarning
	}
	return UrbanSystem{
		Name:     "Energy",
		Severity: sev,
		Metrics: []Metric{
			{"Grid Load", fmt.Sprintf("%d%%", load)},
			{"Renewable", fmt.Sprintf("%d%%", renewable)},
		},
	}
}

func genWater(rng *rand.Rand) UrbanSystem {
	pressure := 40 + rng.Intn(41)
	qualRoll := rng.Intn(10)
	quality := "Good"
	sev := SeverityNormal
	if qualRoll >= 9 {
		quality = "Poor"
		sev = SeverityCritical
	} else if qualRoll >= 6 {
		quality = "Fair"
		sev = SeverityWarning
	}
	if pressure < 45 {
		sev = SeverityCritical
	} else if pressure < 50 && sev < SeverityWarning {
		sev = SeverityWarning
	}
	return UrbanSystem{
		Name:     "Water",
		Severity: sev,
		Metrics: []Metric{
			{"Pressure", fmt.Sprintf("%d psi", pressure)},
			{"Quality", quality},
		},
	}
}

func genWaste(rng *rand.Rand) UrbanSystem {
	capacity := 30 + rng.Intn(61)
	status := "On Schedule"
	sev := SeverityNormal
	if capacity > 85 {
		status = "Critical"
		sev = SeverityCritical
	} else if capacity > 65 {
		status = "Delayed"
		sev = SeverityWarning
	}
	return UrbanSystem{
		Name:     "Waste",
		Severity: sev,
		Metrics: []Metric{
			{"Collection Status", status},
			{"Capacity", fmt.Sprintf("%d%%", capacity)},
		},
	}
}

func genEmergency(rng *rand.Rand) UrbanSystem {
	incidents := rng.Intn(8)
	sev := SeverityNormal
	if incidents > 4 {
		sev = SeverityCritical
	} else if incidents > 1 {
		sev = SeverityWarning
	}
	return UrbanSystem{
		Name:     "Emergency",
		Severity: sev,
		Metrics: []Metric{
			{"Active Incidents", fmt.Sprintf("%d", incidents)},
		},
	}
}

func genAirQuality(rng *rand.Rand) UrbanSystem {
	aqi := rng.Intn(201)
	pm25 := float64(rng.Intn(500)) / 10.0
	sev := SeverityNormal
	if aqi > 150 {
		sev = SeverityCritical
	} else if aqi > 100 {
		sev = SeverityWarning
	}
	return UrbanSystem{
		Name:     "Air Quality",
		Severity: sev,
		Metrics: []Metric{
			{"AQI", fmt.Sprintf("%d", aqi)},
			{"PM2.5", fmt.Sprintf("%.1f µg/m³", pm25)},
		},
	}
}
