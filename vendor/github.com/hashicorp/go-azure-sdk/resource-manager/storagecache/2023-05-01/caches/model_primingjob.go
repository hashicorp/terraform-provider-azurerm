package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrimingJob struct {
	PrimingJobDetails         *string          `json:"primingJobDetails,omitempty"`
	PrimingJobId              *string          `json:"primingJobId,omitempty"`
	PrimingJobName            string           `json:"primingJobName"`
	PrimingJobPercentComplete *float64         `json:"primingJobPercentComplete,omitempty"`
	PrimingJobState           *PrimingJobState `json:"primingJobState,omitempty"`
	PrimingJobStatus          *string          `json:"primingJobStatus,omitempty"`
	PrimingManifestURL        string           `json:"primingManifestUrl"`
}
