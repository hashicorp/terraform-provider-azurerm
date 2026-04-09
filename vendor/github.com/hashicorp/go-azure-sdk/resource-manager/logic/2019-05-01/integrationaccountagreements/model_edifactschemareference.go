package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactSchemaReference struct {
	AssociationAssignedCode    *string `json:"associationAssignedCode,omitempty"`
	MessageId                  string  `json:"messageId"`
	MessageRelease             string  `json:"messageRelease"`
	MessageVersion             string  `json:"messageVersion"`
	SchemaName                 string  `json:"schemaName"`
	SenderApplicationId        *string `json:"senderApplicationId,omitempty"`
	SenderApplicationQualifier *string `json:"senderApplicationQualifier,omitempty"`
}
