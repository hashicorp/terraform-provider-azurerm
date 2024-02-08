package images

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceRange struct {
	Max *int64 `json:"max,omitempty"`
	Min *int64 `json:"min,omitempty"`
}
