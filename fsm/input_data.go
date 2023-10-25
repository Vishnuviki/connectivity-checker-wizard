package fsm

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
	// TODO: implement
	return false
}

func (i *InputData) isDestinationAddressIP() bool {
	// TODO: implement
	return false
}

func (i *InputData) isDestinationAddressValid() bool {
	return i.isDestinationAddressFQDN() || i.isDestinationAddressIP()
}
