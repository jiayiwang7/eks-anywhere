package snow

import (
	"fmt"

	"github.com/aws/eks-anywhere/pkg/aws"
)

type BootstrapCreds struct {
	credsB64 string
	certsB64 string
}

// Set is only used for unit test purpose.
func (b *BootstrapCreds) Set(credsB64, certsB64 string) {
	b.credsB64 = credsB64
	b.certsB64 = certsB64
}

func (p *SnowProvider) setupBootstrapCreds() error {
	creds, err := aws.EncodeFileFromEnv(eksaSnowCredentialsFileKey)
	if err != nil {
		return fmt.Errorf("failed to set up snow credentials: %v", err)
	}
	p.bootstrapCreds.credsB64 = creds

	certs, err := aws.EncodeFileFromEnv(eksaSnowCABundlesFileKey)
	if err != nil {
		return fmt.Errorf("failed to set up snow certificates: %v", err)
	}
	p.bootstrapCreds.certsB64 = certs

	return nil
}
