package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProfile struct {
	MsiClientId   string `json:"msiClientId"`
	MsiObjectId   string `json:"msiObjectId"`
	MsiResourceId string `json:"msiResourceId"`
}
