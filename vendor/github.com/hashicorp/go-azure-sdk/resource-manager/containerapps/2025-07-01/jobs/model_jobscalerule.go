package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobScaleRule struct {
	Auth     *[]ScaleRuleAuth `json:"auth,omitempty"`
	Identity *string          `json:"identity,omitempty"`
	Metadata *interface{}     `json:"metadata,omitempty"`
	Name     *string          `json:"name,omitempty"`
	Type     *string          `json:"type,omitempty"`
}
