package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StopRequest struct {
	SoftStopTimeoutSeconds *int64 `json:"softStopTimeoutSeconds,omitempty"`
}
