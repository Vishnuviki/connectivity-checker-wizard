package fsm

import (
	"testing"
)

func TestInputData(t *testing.T) {
	cases := []struct {
		namespace          string
		port               string
		destinationAddress string

		expectIsIP    bool
		expectIsFQDN  bool
		expectIsValid bool
	}{
		{"ns1", "80", "1.2.3.4", true, false, true},
		{"ns2", "443", "svc1.potato.com", false, true, true},
		{"ns3", "443", ".com", false, false, false},
		{"ns4", "443", "?.", false, false, false},
	}

	for _, c := range cases {
		inputData := NewInputData(c.namespace, c.port, c.destinationAddress)

		if inputData.IsDestinationAddressIP() != c.expectIsIP {
			t.Errorf("expected inputData.isDestinationAddressIP() to be %t, got %t", c.expectIsIP, inputData.IsDestinationAddressIP())
		}
		if inputData.IsDestinationAddressFQDN() != c.expectIsFQDN {
			t.Errorf("expected inputData.isDestinationAddressFQDN() to be %t, got %t", c.expectIsFQDN, inputData.IsDestinationAddressFQDN())
		}
		if inputData.IsDestinationAddressValid() != c.expectIsValid {
			t.Errorf("expected inputData.isDestinationAddressValid() to be %t, got %t", c.expectIsValid, inputData.IsDestinationAddressValid())
		}
	}
}
