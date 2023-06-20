package dataexport

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Destination struct {
	MetaData   *DestinationMetaData `json:"metaData,omitempty"`
	ResourceId string               `json:"resourceId"`
	Type       *Type                `json:"type,omitempty"`
}
