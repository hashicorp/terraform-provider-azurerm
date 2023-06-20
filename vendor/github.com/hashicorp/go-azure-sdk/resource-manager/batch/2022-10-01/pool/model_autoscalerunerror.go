package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoScaleRunError struct {
	Code    string               `json:"code"`
	Details *[]AutoScaleRunError `json:"details,omitempty"`
	Message string               `json:"message"`
}
