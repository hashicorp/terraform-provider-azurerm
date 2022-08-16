package workbooksapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookPropertiesUpdateParameters struct {
	Category       *string   `json:"category,omitempty"`
	Description    *string   `json:"description,omitempty"`
	DisplayName    *string   `json:"displayName,omitempty"`
	Revision       *string   `json:"revision,omitempty"`
	SerializedData *string   `json:"serializedData,omitempty"`
	Tags           *[]string `json:"tags,omitempty"`
}
