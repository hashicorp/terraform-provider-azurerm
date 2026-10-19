// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_key

type KeyResourceModel struct {
	ConfigurationStoreId string         `tfschema:"configuration_store_id"`
	Key                  string         `tfschema:"key"`
	ContentType          string         `tfschema:"content_type"`
	Etag                 string         `tfschema:"etag"`
	Label                string         `tfschema:"label"`
	Value                string         `tfschema:"value"`
	Locked               bool           `tfschema:"locked"`
	Tags                 map[string]any `tfschema:"tags"`
	Type                 string         `tfschema:"type"`
	VaultKeyReference    string         `tfschema:"vault_key_reference"`
}

func (k Resource) ModelObject() any {
	return &KeyResourceModel{}
}
