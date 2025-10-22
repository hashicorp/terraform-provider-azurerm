package recordsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NaptrRecord struct {
	Flags       *string `json:"flags,omitempty"`
	Order       *int64  `json:"order,omitempty"`
	Preference  *int64  `json:"preference,omitempty"`
	Regexp      *string `json:"regexp,omitempty"`
	Replacement *string `json:"replacement,omitempty"`
	Services    *string `json:"services,omitempty"`
}
