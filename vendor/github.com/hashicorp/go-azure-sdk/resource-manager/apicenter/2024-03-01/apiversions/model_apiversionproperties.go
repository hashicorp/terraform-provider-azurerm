package apiversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiVersionProperties struct {
	LifecycleStage LifecycleStage `json:"lifecycleStage"`
	Title          string         `json:"title"`
}
