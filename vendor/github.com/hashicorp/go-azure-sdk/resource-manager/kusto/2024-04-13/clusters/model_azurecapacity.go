package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureCapacity struct {
	Default   int64          `json:"default"`
	Maximum   int64          `json:"maximum"`
	Minimum   int64          `json:"minimum"`
	ScaleType AzureScaleType `json:"scaleType"`
}
