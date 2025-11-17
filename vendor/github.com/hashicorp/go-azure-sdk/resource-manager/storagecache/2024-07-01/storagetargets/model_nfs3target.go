package storagetargets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Nfs3Target struct {
	Target            *string `json:"target,omitempty"`
	UsageModel        *string `json:"usageModel,omitempty"`
	VerificationTimer *int64  `json:"verificationTimer,omitempty"`
	WriteBackTimer    *int64  `json:"writeBackTimer,omitempty"`
}
