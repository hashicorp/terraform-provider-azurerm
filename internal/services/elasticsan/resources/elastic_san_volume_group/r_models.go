// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_san_volume_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

func (r Resource) ModelObject() any {
	return &ElasticSANVolumeGroupResourceModel{}
}

type ElasticSANVolumeGroupResourceModel struct {
	SanId          string                                          `tfschema:"elastic_san_id"`
	EncryptionType string                                          `tfschema:"encryption_type"`
	Encryption     []ElasticSANVolumeGroupResourceEncryptionModel  `tfschema:"encryption"`
	Identity       []identity.ModelSystemAssignedUserAssigned      `tfschema:"identity"`
	Name           string                                          `tfschema:"name"`
	NetworkRule    []ElasticSANVolumeGroupResourceNetworkRuleModel `tfschema:"network_rule"`
	ProtocolType   string                                          `tfschema:"protocol_type"`
}

type ElasticSANVolumeGroupResourceEncryptionModel struct {
	CurrentVersionedKeyExpirationTimestamp string `tfschema:"current_versioned_key_expiration_timestamp"`
	CurrentVersionedKeyId                  string `tfschema:"current_versioned_key_id"`
	UserAssignedIdentityId                 string `tfschema:"user_assigned_identity_id"`
	KeyVaultKeyId                          string `tfschema:"key_vault_key_id"`
	LastKeyRotationTimestamp               string `tfschema:"last_key_rotation_timestamp"`
}

type ElasticSANVolumeGroupResourceNetworkRuleModel struct {
	Action   string `tfschema:"action"`
	SubnetId string `tfschema:"subnet_id"`
}
