package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerificationIPFlowResult struct {
	Access   *Access `json:"access,omitempty"`
	RuleName *string `json:"ruleName,omitempty"`
}
