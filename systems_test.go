package main

import (
	"encoding/json"
	"math/rand"
	"testing"
)

func TestGenerateSystemsCount(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	systems := GenerateSystems(rng)
	if len(systems) != 6 {
		t.Errorf("expected 6 systems, got %d", len(systems))
	}
}

func TestGenerateSystemsNames(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	systems := GenerateSystems(rng)
	expected := []string{"Traffic", "Energy", "Water", "Waste", "Emergency", "Air Quality"}
	for i, name := range expected {
		if systems[i].Name != name {
			t.Errorf("system %d: expected %q, got %q", i, name, systems[i].Name)
		}
	}
}

func TestSeverityString(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityNormal, "Normal"},
		{SeverityWarning, "Warning"},
		{SeverityCritical, "Critical"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

func TestSeverityJSON(t *testing.T) {
	for _, tt := range []struct {
		sev  Severity
		want string
	}{
		{SeverityNormal, `"Normal"`},
		{SeverityWarning, `"Warning"`},
		{SeverityCritical, `"Critical"`},
	} {
		data, err := json.Marshal(tt.sev)
		if err != nil {
			t.Fatalf("marshal severity %d: %v", tt.sev, err)
		}
		if string(data) != tt.want {
			t.Errorf("severity %d JSON = %s, want %s", tt.sev, data, tt.want)
		}
	}
}

func TestSystemMetricCounts(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	systems := GenerateSystems(rng)
	expected := map[string]int{
		"Traffic": 2, "Energy": 2, "Water": 2,
		"Waste": 2, "Emergency": 1, "Air Quality": 2,
	}
	for _, sys := range systems {
		want, ok := expected[sys.Name]
		if !ok {
			t.Errorf("unexpected system: %s", sys.Name)
			continue
		}
		if len(sys.Metrics) != want {
			t.Errorf("%s: expected %d metrics, got %d", sys.Name, want, len(sys.Metrics))
		}
	}
}

func TestGenerateSystemsDeterministic(t *testing.T) {
	rng1 := rand.New(rand.NewSource(123))
	rng2 := rand.New(rand.NewSource(123))
	sys1 := GenerateSystems(rng1)
	sys2 := GenerateSystems(rng2)

	for i := range sys1 {
		if sys1[i].Name != sys2[i].Name {
			t.Errorf("system %d names differ", i)
		}
		if sys1[i].Severity != sys2[i].Severity {
			t.Errorf("system %d severities differ", i)
		}
		for j := range sys1[i].Metrics {
			if sys1[i].Metrics[j].Value != sys2[i].Metrics[j].Value {
				t.Errorf("system %d metric %d values differ: %s vs %s",
					i, j, sys1[i].Metrics[j].Value, sys2[i].Metrics[j].Value)
			}
		}
	}
}

func TestTrafficSeverityRanges(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	counts := [3]int{}
	for i := 0; i < 1000; i++ {
		sys := genTraffic(rng)
		if sys.Name != "Traffic" {
			t.Fatalf("expected Traffic, got %s", sys.Name)
		}
		counts[sys.Severity]++
	}
	for i, name := range []string{"Normal", "Warning", "Critical"} {
		if counts[i] == 0 {
			t.Errorf("traffic: %s severity never appeared in 1000 runs", name)
		}
	}
}

func TestEnergySeverityRanges(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	counts := [3]int{}
	for i := 0; i < 1000; i++ {
		sys := genEnergy(rng)
		counts[sys.Severity]++
	}
	for i, name := range []string{"Normal", "Warning", "Critical"} {
		if counts[i] == 0 {
			t.Errorf("energy: %s severity never appeared in 1000 runs", name)
		}
	}
}

func TestWaterSeverityAndMetrics(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		sys := genWater(rng)
		if len(sys.Metrics) != 2 {
			t.Fatalf("water: expected 2 metrics, got %d", len(sys.Metrics))
		}
		if sys.Metrics[0].Name != "Pressure" {
			t.Errorf("water metric 0: expected Pressure, got %s", sys.Metrics[0].Name)
		}
		if sys.Metrics[1].Name != "Quality" {
			t.Errorf("water metric 1: expected Quality, got %s", sys.Metrics[1].Name)
		}
	}
}

func TestWasteSeverityCorrelation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 500; i++ {
		sys := genWaste(rng)
		status := sys.Metrics[0].Value
		switch sys.Severity {
		case SeverityCritical:
			if status != "Critical" {
				t.Errorf("waste critical severity but status=%s", status)
			}
		case SeverityWarning:
			if status != "Delayed" {
				t.Errorf("waste warning severity but status=%s", status)
			}
		case SeverityNormal:
			if status != "On Schedule" {
				t.Errorf("waste normal severity but status=%s", status)
			}
		}
	}
}

func TestEmergencyMetrics(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		sys := genEmergency(rng)
		if len(sys.Metrics) != 1 {
			t.Fatalf("emergency: expected 1 metric, got %d", len(sys.Metrics))
		}
		if sys.Metrics[0].Name != "Active Incidents" {
			t.Errorf("emergency metric: expected Active Incidents, got %s", sys.Metrics[0].Name)
		}
	}
}

func TestAirQualityMetrics(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		sys := genAirQuality(rng)
		if len(sys.Metrics) != 2 {
			t.Fatalf("air quality: expected 2 metrics, got %d", len(sys.Metrics))
		}
		if sys.Metrics[0].Name != "AQI" {
			t.Errorf("air quality metric 0: expected AQI, got %s", sys.Metrics[0].Name)
		}
		if sys.Metrics[1].Name != "PM2.5" {
			t.Errorf("air quality metric 1: expected PM2.5, got %s", sys.Metrics[1].Name)
		}
	}
}

func TestSystemsJSONRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	systems := GenerateSystems(rng)
	data, err := json.Marshal(systems)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("empty JSON output")
	}
	// Verify it's valid JSON by unmarshaling into generic structure
	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result) != 6 {
		t.Errorf("expected 6 systems in JSON, got %d", len(result))
	}
}
