package fsm

import (
	gonet "github.com/THREATINT/go-net"
)

type InputData struct {
	SourceNamespace    string `form:"sourceNamespace"`
	DestinationPort    string `form:"destinationPort"`
	DestinationAddress string `form:"destinationAddress"`
}

func newInputData(sourceNamespace, destinationPort, destinationAddress string) *InputData {
	return &InputData{
		SourceNamespace:    sourceNamespace,
		DestinationPort:    destinationPort,
		DestinationAddress: destinationAddress,
	}
}

func (i *InputData) isDestinationAddressFQDN() bool {
	return gonet.IsFQDN(i.DestinationAddress)
}

func (i *InputData) isDestinationAddressIP() bool {
	return gonet.IsIPAddr(i.DestinationAddress)
}

func (i *InputData) isDestinationAddressValid() bool {
	return i.isDestinationAddressFQDN() || i.isDestinationAddressIP()
}
