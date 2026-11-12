// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_hybrid_connection

type RelayHybridConnectionResourceModel struct {
	Name                        string `tfschema:"name"`
	ResourceGroupName           string `tfschema:"resource_group_name"`
	RelayNamespaceName          string `tfschema:"relay_namespace_name"`
	RequiresClientAuthorization bool   `tfschema:"requires_client_authorization"`
	UserMetadata                string `tfschema:"user_metadata"`
}

func (r Resource) ModelObject() any {
	return &RelayHybridConnectionResourceModel{}
}
