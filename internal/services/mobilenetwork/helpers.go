// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22634
// some resources defined both systemAssigned and userAssigned Identity type in Swagger but only support userAssigned Identity,
// so add a workaround to convert type here.
func expandMobileNetworkLegacyToUserAssignedIdentity(input []identity.ModelUserAssigned) (*identity.LegacySystemAndUserAssignedMap, error) {
	if len(input) == 0 {
		return nil, nil
	}

	identityValue, err := identity.ExpandUserAssignedMapFromModel(input)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	output := identity.LegacySystemAndUserAssignedMap{
		Type:        identityValue.Type,
		IdentityIds: identityValue.IdentityIds,
	}

	return &output, nil
}

func flattenMobileNetworkUserAssignedToNetworkLegacyIdentity(input *identity.LegacySystemAndUserAssignedMap) ([]identity.ModelUserAssigned, error) {
	if input == nil {
		return nil, nil
	}

	tmp := identity.UserAssignedMap{
		Type:        input.Type,
		IdentityIds: input.IdentityIds,
	}

	output, err := identity.FlattenUserAssignedMapToModel(&tmp)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	return *output, nil
}
