package device

// CreateKeysAndCertificateAcceptedCh holds the bytes of a CreateKeysAndCertificateAccepted message.
type CreateKeysAndCertificateAcceptedCh []byte

// CreateKeysAndCertificateAccepted holds the data from an accepted request to create a new key and certificate.
type CreateKeysAndCertificateAccepted struct {
	CertificateID             string `json:"certificateId"`
	CertificatePem            string `json:"certificatePem"`
	PrivateKey                string `json:"privateKey"`
	CertificateOwnershipToken string `json:"certificateOwnershipToken"`
}

// AWSMQTTErrorCh holds the bytes of a CreateKeysAndCertificateRejected message.
type AWSMQTTErrorCh []byte

// AWSMQTTError holds the data from a rejected request
type AWSMQTTError struct {
	StatusCode   int    `json:"statusCode"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// RegisterThingRequest holds the values needed to make a request to register a new thing.
type RegisterThingRequest struct {
	TemplateName              string            `json:"-"`
	CertificateOwnershipToken string            `json:"certificateOwnershipToken"`
	Parameters                map[string]string `json:"parameters"`
}

// RegisterThingResponse holds the values returned from a successful registration request.
type RegisterThingResponse struct {
	ThingName           string            `json:"thingName"`
	DeviceConfiguration map[string]string `json:"deviceConfiguration"`
}

// RegisterThingAcceptedCh holds the bytes of a RegisterThingAccepted message.
type RegisterThingAcceptedCh []byte
