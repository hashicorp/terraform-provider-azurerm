package oraclesubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudAccountDetails struct {
	CloudAccountHomeRegion *string `json:"cloudAccountHomeRegion,omitempty"`
	CloudAccountName       *string `json:"cloudAccountName,omitempty"`
}
