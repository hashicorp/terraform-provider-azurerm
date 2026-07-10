package synapse

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

const synapseWorkspaceBaseURIFmt = "https://%s.%s"

func NewSynapseWorkspaceBaseURI(env environments.Environment, workspaceName string) (string, error) {
	suffix, ok := env.Synapse.DomainSuffix()
	if !ok {
		return "", fmt.Errorf("determining Synapse domain suffix for environment `%s`", env.Name)
	}

	return fmt.Sprintf(synapseWorkspaceBaseURIFmt, workspaceName, *suffix), nil
}
