package constants

const (
	// rule names
	DISPATCH_IP_RULE    = "dispatchIPRule"
	DNS_LOOK_UP_RULE    = "dnsLookUPRule"
	VALIDATION_RULE     = "validationRule"
	NETWORK_POLICY_RULE = "networkPolicyRule"

	// response messages
	PAGE_NOT_FOUND_MESSAGE              = "Page Not Found"
	INVALID_DESTINATION_ADDRESS_MESSAGE = "Validation failed: Destination address not valid, must be IP or FQDN"
	INVALID_PORT_NUMBER_MESSAGE         = "Validation failed: Invalid Destination Port"
)
