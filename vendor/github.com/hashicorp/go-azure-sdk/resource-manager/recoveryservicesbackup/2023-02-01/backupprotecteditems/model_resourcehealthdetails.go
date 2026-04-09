package backupprotecteditems

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceHealthDetails struct {
	Code            *int64    `json:"code,omitempty"`
	Message         *string   `json:"message,omitempty"`
	Recommendations *[]string `json:"recommendations,omitempty"`
	Title           *string   `json:"title,omitempty"`
}
