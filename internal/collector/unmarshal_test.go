package collector

import (
	"encoding/json"
	"testing"
)

func TestXfloat64(t *testing.T) {
	cases := map[string]float64{
		`1600`:    1600,
		`"1600"`:  1600,
		`" 12.5"`: 12.5,
		`null`:    0,
		`""`:      0,
		`"n/a"`:   0,
	}

	for in, want := range cases {
		var got xfloat64
		if err := json.Unmarshal([]byte(in), &got); err != nil {
			t.Errorf("%s: unexpected error: %v", in, err)
			continue
		}
		if float64(got) != want {
			t.Errorf("%s: got %v, want %v", in, got, want)
		}
	}
}

// Western Digital Ultrastar Data102 IOMs report PowerCapacityWatts as a string
// and fan and temperature MemberIds as numbers.
func TestDecodeWesternDigitalData102(t *testing.T) {
	power := `{
		"PowerControl": [],
		"PowerSupplies": [
			{"MemberId": "A", "Name": "PSU A", "PowerCapacityWatts": "1600", "LastPowerOutputWatts": null, "LineInputVoltage": null, "Status": {"Health": "OK", "State": "Enabled"}},
			{"MemberId": "B", "Name": "PSU B", "PowerCapacityWatts": "1600", "LastPowerOutputWatts": null, "LineInputVoltage": null, "Status": {"Health": "OK", "State": "Enabled"}}
		]
	}`

	var p PowerResponse
	if err := json.Unmarshal([]byte(power), &p); err != nil {
		t.Fatalf("power: %v", err)
	}
	if len(p.PowerSupplies) != 2 {
		t.Fatalf("power: got %d supplies, want 2", len(p.PowerSupplies))
	}
	if got := float64(p.PowerSupplies[0].PowerCapacityWatts); got != 1600 {
		t.Errorf("power: capacity got %v, want 1600", got)
	}

	thermal := `{
		"Fans": [
			{"MemberId": 1, "Name": null, "Reading": 53, "ReadingUnits": "Percent", "Status": {"Health": "OK", "State": "Enabled"}}
		],
		"Temperatures": [
			{"MemberId": 1, "Name": "TEMP SLOT 000", "ReadingCelsius": 31, "Status": {"Health": "OK", "State": "Enabled"}}
		]
	}`

	var th ThermalResponse
	if err := json.Unmarshal([]byte(thermal), &th); err != nil {
		t.Fatalf("thermal: %v", err)
	}
	if got := th.Fans[0].GetId(0); got != "1" {
		t.Errorf("thermal: fan id got %q, want %q", got, "1")
	}
	if got := th.Temperatures[0].GetId(0); got != "1" {
		t.Errorf("thermal: temperature id got %q, want %q", got, "1")
	}
}
