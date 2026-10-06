package apitagdescription

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagDescriptionContractProperties struct {
	Description             *string `json:"description,omitempty"`
	DisplayName             *string `json:"displayName,omitempty"`
	ExternalDocsDescription *string `json:"externalDocsDescription,omitempty"`
	ExternalDocsURL         *string `json:"externalDocsUrl,omitempty"`
	TagId                   *string `json:"tagId,omitempty"`
}
