package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequestsBasedTrigger struct {
	Count        *int64  `json:"count,omitempty"`
	TimeInterval *string `json:"timeInterval,omitempty"`
}
