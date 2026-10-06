package batchaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchAccountKeys struct {
	AccountName *string `json:"accountName,omitempty"`
	Primary     *string `json:"primary,omitempty"`
	Secondary   *string `json:"secondary,omitempty"`
}
