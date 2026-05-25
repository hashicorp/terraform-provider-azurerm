// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

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
