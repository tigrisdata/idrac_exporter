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

func TestFillIdentity(t *testing.T) {
	chassis := ChassisResponse{Manufacturer: "HGST", Model: "H4102-J", SerialNumber: "SN-CHASSIS", SKU: "SKU-CHASSIS"}

	empty := SystemResponse{}
	empty.FillIdentity(&chassis)
	if empty.Manufacturer != "HGST" || empty.Model != "H4102-J" || empty.SerialNumber != "SN-CHASSIS" || empty.SKU != "SKU-CHASSIS" {
		t.Errorf("empty system: got %q %q %q %q", empty.Manufacturer, empty.Model, empty.SerialNumber, empty.SKU)
	}

	full := SystemResponse{Manufacturer: "Supermicro", Model: "AS -1115HS-TNR", SerialNumber: " ", SKU: "SKU-SYSTEM"}
	full.FillIdentity(&chassis)
	if full.Manufacturer != "Supermicro" || full.Model != "AS -1115HS-TNR" || full.SKU != "SKU-SYSTEM" {
		t.Errorf("populated system fields were overwritten: %q %q %q", full.Manufacturer, full.Model, full.SKU)
	}
	if full.SerialNumber != "SN-CHASSIS" {
		t.Errorf("blank serial: got %q, want %q", full.SerialNumber, "SN-CHASSIS")
	}
}
