package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Frontend struct {
	Replicas int64  `json:"replicas"`
	Workers  *int64 `json:"workers,omitempty"`
}
