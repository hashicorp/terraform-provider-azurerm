package customlocations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnabledResourceTypeProperties struct {
	ClusterExtensionId *string                                              `json:"clusterExtensionId,omitempty"`
	ExtensionType      *string                                              `json:"extensionType,omitempty"`
	TypesMetadata      *[]EnabledResourceTypePropertiesTypesMetadataInlined `json:"typesMetadata,omitempty"`
}
