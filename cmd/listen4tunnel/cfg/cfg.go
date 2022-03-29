package cfg

var (
	// Endpoint represents the AWS IoT endpoint
	Endpoint string
	// ThingName is the AWS IoT Thing name for the current device
	ThingName string
	// PrivateKeyPath is the path to the device private key
	PrivateKeyPath = "/certs/device.private.key"
	// CertificatePath is the path to the device certificate
	CertificatePath = "/certs/device.certificate.pem"
	// RootCAPath is the path to the root CA certificate
	RootCAPath = "/certs/ca.pem"
)
