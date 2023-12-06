package diagnostic

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataMasking struct {
	Headers     *[]DataMaskingEntity `json:"headers,omitempty"`
	QueryParams *[]DataMaskingEntity `json:"queryParams,omitempty"`
}
