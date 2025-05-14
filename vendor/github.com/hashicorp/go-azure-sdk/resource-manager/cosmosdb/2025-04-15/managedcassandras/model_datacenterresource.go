package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataCenterResource struct {
	Id         *string                       `json:"id,omitempty"`
	Name       *string                       `json:"name,omitempty"`
	Properties *DataCenterResourceProperties `json:"properties,omitempty"`
	Type       *string                       `json:"type,omitempty"`
}
