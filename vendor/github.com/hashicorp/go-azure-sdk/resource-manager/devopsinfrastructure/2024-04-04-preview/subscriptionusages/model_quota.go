package subscriptionusages

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Quota struct {
	CurrentValue int64      `json:"currentValue"`
	Id           string     `json:"id"`
	Limit        int64      `json:"limit"`
	Name         *QuotaName `json:"name,omitempty"`
	Unit         string     `json:"unit"`
}
