package updates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateStateProperties struct {
	NotifyMessage      *string  `json:"notifyMessage,omitempty"`
	ProgressPercentage *float64 `json:"progressPercentage,omitempty"`
}
