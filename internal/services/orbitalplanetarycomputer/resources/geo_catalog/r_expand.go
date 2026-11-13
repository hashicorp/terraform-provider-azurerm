// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_catalog

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbitalplanetarycomputer/2026-04-15/geocatalogs"
)

func (r Resource) expandIdentity(input []identity.ModelSystemAssignedUserAssigned) *geocatalogs.ManagedServiceIdentityUpdate {
	managedServiceIdentityUpdate := &geocatalogs.ManagedServiceIdentityUpdate{
		Type:                   pointer.To(geocatalogs.ManagedServiceIdentityTypeNone),
		UserAssignedIdentities: &map[string]geocatalogs.UserAssignedIdentity{},
	}

	if len(input) == 0 {
		return managedServiceIdentityUpdate
	}

	managedServiceIdentityUpdate.Type = pointer.ToEnum[geocatalogs.ManagedServiceIdentityType](string(input[0].Type))

	for _, identityId := range input[0].IdentityIds {
		(*managedServiceIdentityUpdate.UserAssignedIdentities)[identityId] = geocatalogs.UserAssignedIdentity{}
	}

	return managedServiceIdentityUpdate
}
