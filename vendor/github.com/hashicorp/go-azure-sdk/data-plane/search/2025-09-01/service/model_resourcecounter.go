package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceCounter struct {
	Quota *int64 `json:"quota,omitempty"`
	Usage int64  `json:"usage"`
}
