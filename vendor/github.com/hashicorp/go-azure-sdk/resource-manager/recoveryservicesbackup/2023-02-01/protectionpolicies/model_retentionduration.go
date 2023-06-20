package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionDuration struct {
	Count        *int64                 `json:"count,omitempty"`
	DurationType *RetentionDurationType `json:"durationType,omitempty"`
}
