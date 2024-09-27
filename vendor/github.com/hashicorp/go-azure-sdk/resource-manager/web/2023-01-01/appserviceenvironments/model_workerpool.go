package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkerPool struct {
	ComputeMode   *ComputeModeOptions `json:"computeMode,omitempty"`
	InstanceNames *[]string           `json:"instanceNames,omitempty"`
	WorkerCount   *int64              `json:"workerCount,omitempty"`
	WorkerSize    *string             `json:"workerSize,omitempty"`
	WorkerSizeId  *int64              `json:"workerSizeId,omitempty"`
}
