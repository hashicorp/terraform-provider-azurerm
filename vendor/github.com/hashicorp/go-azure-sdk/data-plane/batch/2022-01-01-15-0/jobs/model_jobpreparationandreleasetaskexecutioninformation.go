package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobPreparationAndReleaseTaskExecutionInformation struct {
	JobPreparationTaskExecutionInfo *JobPreparationTaskExecutionInformation `json:"jobPreparationTaskExecutionInfo,omitempty"`
	JobReleaseTaskExecutionInfo     *JobReleaseTaskExecutionInformation     `json:"jobReleaseTaskExecutionInfo,omitempty"`
	NodeId                          *string                                 `json:"nodeId,omitempty"`
	NodeURL                         *string                                 `json:"nodeUrl,omitempty"`
	PoolId                          *string                                 `json:"poolId,omitempty"`
}
