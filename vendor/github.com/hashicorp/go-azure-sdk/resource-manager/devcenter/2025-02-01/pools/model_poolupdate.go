package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolUpdate struct {
	Location   *string               `json:"location,omitempty"`
	Properties *PoolUpdateProperties `json:"properties,omitempty"`
	Tags       *map[string]string    `json:"tags,omitempty"`
}
