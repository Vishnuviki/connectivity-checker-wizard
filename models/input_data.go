package models

import (
	gonet "github.com/THREATINT/go-net"
)

type InputData struct {
	SourceNamespace    string `form:"sourceNamespace"`
	DestinationPort    string `form:"destinationPort"`
	DestinationAddress string `form:"destinationAddress"`
}

func NewInputData(sourceNamespace, destinationPort, destinationAddress string) *InputData {
	return &InputData{
		SourceNamespace:    sourceNamespace,
		DestinationPort:    destinationPort,
		DestinationAddress: destinationAddress,
	}
}

func (i *InputData) IsDestinationAddressFQDN() bool {
	return gonet.IsFQDN(i.DestinationAddress)
}

func (i *InputData) IsDestinationAddressIP() bool {
	return gonet.IsIPAddr(i.DestinationAddress)
}

func (i *InputData) IsDestinationAddressValid() bool {
	return i.IsDestinationAddressFQDN() || i.IsDestinationAddressIP()
}