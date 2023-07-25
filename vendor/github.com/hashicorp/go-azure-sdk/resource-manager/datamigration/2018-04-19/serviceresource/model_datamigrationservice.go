package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataMigrationService struct {
	Etag       *string                         `json:"etag,omitempty"`
	Id         *string                         `json:"id,omitempty"`
	Kind       *string                         `json:"kind,omitempty"`
	Location   string                          `json:"location"`
	Name       *string                         `json:"name,omitempty"`
	Properties *DataMigrationServiceProperties `json:"properties,omitempty"`
	Sku        *ServiceSku                     `json:"sku,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
	Type       *string                         `json:"type,omitempty"`
}
