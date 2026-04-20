package views

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KpiProperties struct {
	Enabled *bool        `json:"enabled,omitempty"`
	Id      *string      `json:"id,omitempty"`
	Type    *KpiTypeType `json:"type,omitempty"`
}
