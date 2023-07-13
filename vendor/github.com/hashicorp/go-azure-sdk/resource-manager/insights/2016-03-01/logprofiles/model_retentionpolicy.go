package logprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetentionPolicy struct {
	Days    int64 `json:"days"`
	Enabled bool  `json:"enabled"`
}
