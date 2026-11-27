// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service_datasource_blob

import "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/schema"

type SearchServiceDatasourceBlobModel struct {
	Name                      string                                      `tfschema:"name"`
	SearchServiceEndpoint     string                                      `tfschema:"search_service_endpoint"`
	ContainerName             string                                      `tfschema:"container_name"`
	ConnectionStringWOVersion int64                                       `tfschema:"connection_string_wo_version"`
	Description               string                                      `tfschema:"description"`
	ContainerQuery            string                                      `tfschema:"container_query"`
	SoftDeleteColumnName      string                                      `tfschema:"soft_delete_column_name"`
	SoftDeleteMarkerValue     string                                      `tfschema:"soft_delete_marker_value"`
	EncryptionKey             []schema.SearchDatasourceEncryptionKeyModel `tfschema:"encryption_key"`
}

func (r Resource) ModelObject() any {
	return &SearchServiceDatasourceBlobModel{}
}
