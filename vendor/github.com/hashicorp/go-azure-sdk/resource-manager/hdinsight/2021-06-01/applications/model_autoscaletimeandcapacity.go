package applications

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleTimeAndCapacity struct {
	MaxInstanceCount *int64  `json:"maxInstanceCount,omitempty"`
	MinInstanceCount *int64  `json:"minInstanceCount,omitempty"`
	Time             *string `json:"time,omitempty"`
}
