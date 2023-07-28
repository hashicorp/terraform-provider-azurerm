package quota

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateWorkspaceQuotas struct {
	Id     *string    `json:"id,omitempty"`
	Limit  *int64     `json:"limit,omitempty"`
	Status *Status    `json:"status,omitempty"`
	Type   *string    `json:"type,omitempty"`
	Unit   *QuotaUnit `json:"unit,omitempty"`
}
