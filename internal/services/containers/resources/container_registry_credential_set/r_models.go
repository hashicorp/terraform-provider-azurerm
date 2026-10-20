// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_credential_set

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ContainerRegistryCredentialSetModel struct {
	Name                     string                         `tfschema:"name"`
	ContainerRegistryId      string                         `tfschema:"container_registry_id"`
	LoginServer              string                         `tfschema:"login_server"`
	AuthenticationCredential []AuthenticationCredential     `tfschema:"authentication_credentials"`
	Identity                 []identity.ModelSystemAssigned `tfschema:"identity"`
}

func (Resource) ModelObject() any {
	return &ContainerRegistryCredentialSetModel{}
}
