package models

type FormData struct {
	SourceNamespace    string `form:"sourceNamespace"`
	DestinationPort    string `form:"destinationPort"`
	DestinationAddress string `form:"destinationAddress"`
}
