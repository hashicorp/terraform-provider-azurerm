package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleRuleAuth struct {
	SecretRef        *string `json:"secretRef,omitempty"`
	TriggerParameter *string `json:"triggerParameter,omitempty"`
}
