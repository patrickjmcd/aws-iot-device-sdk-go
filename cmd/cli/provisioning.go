package provisioning

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// provisioningCertAndKey holds the information needed to provision a device with the certificate
// and key.
type provisioningCertAndKey struct {
	CertificatePem string `json:"certificatePem"`
	PrivateKey     string `json:"privateKey"`
	PublicKey      string `json:"publicKey"`
	CertificateARN string `json:"certificateArn"`
	CertificateID  string `json:"certificateId"`
}

func prepProvisioningCertificateAndKey(ctx context.Context, outputPath string) error {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	iotSvc := iot.NewFromConfig(cfg)

	createOutput, err := iotSvc.CreateKeysAndCertificate(ctx, &iot.CreateKeysAndCertificateInput{
		SetAsActive: true,
	})
	if err != nil {
		return err
	}

	certAndKey := provisioningCertAndKey{
		CertificatePem: *createOutput.CertificatePem,
		PrivateKey:     *createOutput.KeyPair.PrivateKey,
		PublicKey:      *createOutput.KeyPair.PublicKey,
		CertificateARN: *createOutput.CertificateArn,
		CertificateID:  *createOutput.CertificateId,
	}

	// write the CertificatePem to a file
	err = os.WriteFile(fmt.Sprintf("%s/cert.pem", outputPath), []byte(certAndKey.CertificatePem), 0644)
	if err != nil {
		return err
	}

	// write the PrivateKey to a file
	err = os.WriteFile(fmt.Sprintf("%s/private.key", outputPath), []byte(certAndKey.PrivateKey), 0644)
	if err != nil {
		return err
	}

	// write the public key to a file
	err = os.WriteFile(fmt.Sprintf("%s/public.key", outputPath), []byte(certAndKey.PublicKey), 0644)
	if err != nil {
		return err
	}

	// write the certifcate json information to a file
	outputJSON, err := json.MarshalIndent(certAndKey, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/cert.json", outputPath), outputJSON, 0644)

	return err

}
