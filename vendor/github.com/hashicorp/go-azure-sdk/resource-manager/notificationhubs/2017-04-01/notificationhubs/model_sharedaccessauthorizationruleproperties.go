package notificationhubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedAccessAuthorizationRuleProperties struct {
	ClaimType    *string         `json:"claimType,omitempty"`
	ClaimValue   *string         `json:"claimValue,omitempty"`
	CreatedTime  *string         `json:"createdTime,omitempty"`
	KeyName      *string         `json:"keyName,omitempty"`
	ModifiedTime *string         `json:"modifiedTime,omitempty"`
	PrimaryKey   *string         `json:"primaryKey,omitempty"`
	Revision     *int64          `json:"revision,omitempty"`
	Rights       *[]AccessRights `json:"rights,omitempty"`
	SecondaryKey *string         `json:"secondaryKey,omitempty"`
}
