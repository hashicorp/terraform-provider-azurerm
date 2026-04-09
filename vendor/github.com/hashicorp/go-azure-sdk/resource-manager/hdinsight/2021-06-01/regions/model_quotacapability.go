package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaCapability struct {
	CoresUsed       *int64                     `json:"coresUsed,omitempty"`
	MaxCoresAllowed *int64                     `json:"maxCoresAllowed,omitempty"`
	RegionalQuotas  *[]RegionalQuotaCapability `json:"regionalQuotas,omitempty"`
}
