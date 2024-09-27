package defenderforstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OnUploadProperties struct {
	CapGBPerMonth *int64 `json:"capGBPerMonth,omitempty"`
	IsEnabled     *bool  `json:"isEnabled,omitempty"`
}
