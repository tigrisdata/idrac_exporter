package collector

import (
	"encoding/json"
	"testing"
)

func TestIntrusion2value(t *testing.T) {
	cases := map[string]int{
		"":                  -1,
		"Normal":            0,
		"HardwareIntrusion": 1,
		"TamperingDetected": 1,
	}

	for in, want := range cases {
		if got := intrusion2value(in); got != want {
			t.Errorf("%q: got %d, want %d", in, got, want)
		}
	}
}

func TestDecodeChassisPhysicalSecurity(t *testing.T) {
	data := `{"Id": "1", "PhysicalSecurity": {"IntrusionSensorNumber": 170, "IntrusionSensor": "HardwareIntrusion", "IntrusionSensorReArm": "Manual"}}`

	var c ChassisResponse
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		t.Fatal(err)
	}
	if c.PhysicalSecurity == nil {
		t.Fatal("PhysicalSecurity not decoded")
	}
	if got := c.PhysicalSecurity.IntrusionSensor; got != "HardwareIntrusion" {
		t.Errorf("got %q, want %q", got, "HardwareIntrusion")
	}
}
