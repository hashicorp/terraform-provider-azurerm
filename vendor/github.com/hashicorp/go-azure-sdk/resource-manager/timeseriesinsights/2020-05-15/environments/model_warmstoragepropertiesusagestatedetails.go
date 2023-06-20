package environments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WarmStoragePropertiesUsageStateDetails struct {
	CurrentCount *int64 `json:"currentCount,omitempty"`
	MaxCount     *int64 `json:"maxCount,omitempty"`
}
