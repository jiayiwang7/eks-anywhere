package snow

import (
	"fmt"

	"github.com/aws/eks-anywhere/pkg/aws"
	"github.com/aws/eks-anywhere/pkg/cluster"
)

func (p *SnowProvider) setupBootstrapCreds(clusterSpec *cluster.Spec) error {
	creds, err := aws.EncodeFileFromEnv(eksaSnowCredentialsFileKey)
	if err != nil {
		return fmt.Errorf("failed to set up snow credentials: %v", err)
	}

	certs, err := aws.EncodeFileFromEnv(eksaSnowCABundlesFileKey)
	if err != nil {
		return fmt.Errorf("failed to set up snow certificates: %v", err)
	}

	// if eks-a credentials secret not specified in the cluster spec file, create from ENVs
	if clusterSpec.SnowCredentialsSecret == nil {
		clusterSpec.SnowCredentialsSecret = EksaCredentialsSecret(clusterSpec, creds, certs)
	}

	return nil
}
