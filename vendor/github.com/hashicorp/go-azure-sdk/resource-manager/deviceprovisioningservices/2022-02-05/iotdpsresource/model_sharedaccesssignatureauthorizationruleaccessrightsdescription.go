package iotdpsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedAccessSignatureAuthorizationRuleAccessRightsDescription struct {
	KeyName      string                  `json:"keyName"`
	PrimaryKey   *string                 `json:"primaryKey,omitempty"`
	Rights       AccessRightsDescription `json:"rights"`
	SecondaryKey *string                 `json:"secondaryKey,omitempty"`
}
