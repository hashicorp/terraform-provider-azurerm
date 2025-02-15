package containerapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomScaleRule struct {
	Auth     *[]ScaleRuleAuth   `json:"auth,omitempty"`
	Metadata *map[string]string `json:"metadata,omitempty"`
	Type     *string            `json:"type,omitempty"`
}
