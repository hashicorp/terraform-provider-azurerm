package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleCapacity struct {
	Default string `json:"default"`
	Maximum string `json:"maximum"`
	Minimum string `json:"minimum"`
}
