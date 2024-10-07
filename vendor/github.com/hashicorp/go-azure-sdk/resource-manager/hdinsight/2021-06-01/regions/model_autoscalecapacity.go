package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleCapacity struct {
	MaxInstanceCount *int64 `json:"maxInstanceCount,omitempty"`
	MinInstanceCount *int64 `json:"minInstanceCount,omitempty"`
}
