// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// NOTE: @tombuildsstuff - this entire file is not a recommended pattern and is a workaround for a resource which is
// managing two types of Access Policy (an Object ID-only and with an Application ID) instead of being split into
// two separate resources.
//
// If you find yourself needing to emulate this pattern, consider splitting/specialising the resource as required.

var _ resourceids.Id = AccessPolicyId{}

type AccessPolicyId struct {
	applicationId *AccessPolicyApplicationId
	objectId      *AccessPolicyObjectId
}

func NewAccessPolicyId(keyVaultId commonids.KeyVaultId, objectId, applicationId string) AccessPolicyId {
	out := AccessPolicyId{}
	if applicationId != "" {
		id := NewAccessPolicyApplicationID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, objectId, applicationId)
		out.applicationId = &id
	} else {
		id := NewAccessPolicyObjectID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroupName, keyVaultId.VaultName, objectId)
		out.objectId = &id
	}
	return out
}

func AccessPolicyID(input string) (*AccessPolicyId, error) {
	accessPolicyObjectId, _ := AccessPolicyObjectID(input)
	if accessPolicyObjectId != nil {
		return &AccessPolicyId{
			objectId: accessPolicyObjectId,
		}, nil
	}
	accessPolicyApplicationId, _ := AccessPolicyApplicationID(input)
	if accessPolicyApplicationId != nil {
		return &AccessPolicyId{
			applicationId: accessPolicyApplicationId,
		}, nil
	}

	return nil, fmt.Errorf("%q didn't parse as either an Object ID ('{keyVaultId}/objectId/{objId}') or an Application ID ('{keyVaultId}/objectId/{objId}/applicationId/{appId}')", input)
}

func (a AccessPolicyId) ID() string {
	if a.applicationId != nil {
		return a.applicationId.ID()
	}

	// whilst this is a pointer, as it has to be either/or it's fine
	return a.objectId.ID()
}

func (a AccessPolicyId) String() string {
	if a.applicationId != nil {
		return a.applicationId.String()
	}

	// whilst this is a pointer, as it has to be either/or it's fine
	return a.objectId.String()
}

func (a AccessPolicyId) ApplicationId() string {
	if a.applicationId != nil {
		return a.applicationId.ApplicationIdName
	}

	return ""
}

func (a AccessPolicyId) KeyVaultId() commonids.KeyVaultId {
	if a.applicationId != nil {
		return commonids.NewKeyVaultID(a.applicationId.SubscriptionId, a.applicationId.ResourceGroup, a.applicationId.VaultName)
	}

	// whilst this is a pointer, as it has to be either/or it's fine
	return commonids.NewKeyVaultID(a.objectId.SubscriptionId, a.objectId.ResourceGroup, a.objectId.VaultName)
}

func (a AccessPolicyId) ObjectID() string {
	if a.applicationId != nil {
		return a.applicationId.ObjectIdName
	}

	// whilst this is a pointer, as it has to be either/or it's fine
	return a.objectId.ObjectIdName
}
