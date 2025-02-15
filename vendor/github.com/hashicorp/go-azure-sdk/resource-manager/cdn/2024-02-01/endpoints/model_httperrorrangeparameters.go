package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPErrorRangeParameters struct {
	Begin *int64 `json:"begin,omitempty"`
	End   *int64 `json:"end,omitempty"`
}
