package models

type CiliumResponse struct {
	IsHostname  bool `json:"isHostname"`
	IsIPAddress bool `json:"isIPAddress"`
}
