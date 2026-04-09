package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MxRecord struct {
	Exchange   *string `json:"exchange,omitempty"`
	Preference *int64  `json:"preference,omitempty"`
}
