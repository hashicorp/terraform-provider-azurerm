package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type X12SchemaReference struct {
	MessageId           string  `json:"messageId"`
	SchemaName          string  `json:"schemaName"`
	SchemaVersion       string  `json:"schemaVersion"`
	SenderApplicationId *string `json:"senderApplicationId,omitempty"`
}
