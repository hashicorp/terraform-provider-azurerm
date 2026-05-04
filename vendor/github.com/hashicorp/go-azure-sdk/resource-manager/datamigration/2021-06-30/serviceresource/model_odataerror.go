package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ODataError struct {
	Code    *string       `json:"code,omitempty"`
	Details *[]ODataError `json:"details,omitempty"`
	Message *string       `json:"message,omitempty"`
}
