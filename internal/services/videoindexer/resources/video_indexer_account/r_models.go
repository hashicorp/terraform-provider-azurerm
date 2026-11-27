// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package video_indexer_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type AccountModel struct {
	Name                string                                     `tfschema:"name"`
	ResourceGroup       string                                     `tfschema:"resource_group_name"`
	Location            string                                     `tfschema:"location"`
	Storage             []StorageModel                             `tfschema:"storage"`
	Identity            []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccess string                                     `tfschema:"public_network_access"`
	Tags                map[string]string                          `tfschema:"tags"`
}

type StorageModel struct {
	StorageAccountId       string `tfschema:"storage_account_id"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
}

func (r Resource) ModelObject() any {
	return &AccountModel{}
}
