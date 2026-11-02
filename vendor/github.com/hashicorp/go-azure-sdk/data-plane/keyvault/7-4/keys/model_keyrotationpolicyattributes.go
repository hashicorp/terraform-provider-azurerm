package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyRotationPolicyAttributes struct {
	Created    *int64  `json:"created,omitempty"`
	ExpiryTime *string `json:"expiryTime,omitempty"`
	Updated    *int64  `json:"updated,omitempty"`
}
