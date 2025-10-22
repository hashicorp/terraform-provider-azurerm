// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func expandPaloAltoLegacyToUserAssignedIdentity(input []identity.ModelUserAssigned) (*identity.LegacySystemAndUserAssignedMap, error) {
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

func flattenPaloAltoUserAssignedToLegacyIdentity(input *identity.LegacySystemAndUserAssignedMap) ([]identity.ModelUserAssigned, error) {
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
