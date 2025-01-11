package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateSyncCompleteCommandOutput struct {
	Errors *[]ReportableException `json:"errors,omitempty"`
	Id     *string                `json:"id,omitempty"`
}
