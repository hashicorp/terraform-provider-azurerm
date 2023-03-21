package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OptimizedAutoscale struct {
	IsEnabled bool  `json:"isEnabled"`
	Maximum   int64 `json:"maximum"`
	Minimum   int64 `json:"minimum"`
	Version   int64 `json:"version"`
}
