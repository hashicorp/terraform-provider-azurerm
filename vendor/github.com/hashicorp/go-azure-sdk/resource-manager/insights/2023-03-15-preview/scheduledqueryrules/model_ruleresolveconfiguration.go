package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuleResolveConfiguration struct {
	AutoResolved  *bool   `json:"autoResolved,omitempty"`
	TimeToResolve *string `json:"timeToResolve,omitempty"`
}
