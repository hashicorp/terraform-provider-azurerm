package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RollbackStatusInfo struct {
	FailedRolledbackInstanceCount       *int64    `json:"failedRolledbackInstanceCount,omitempty"`
	RollbackError                       *ApiError `json:"rollbackError,omitempty"`
	SuccessfullyRolledbackInstanceCount *int64    `json:"successfullyRolledbackInstanceCount,omitempty"`
}
