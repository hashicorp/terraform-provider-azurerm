package quota

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaUpdateParameters struct {
	Location *string                `json:"location,omitempty"`
	Value    *[]QuotaBaseProperties `json:"value,omitempty"`
}
