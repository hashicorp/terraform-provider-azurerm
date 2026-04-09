package connections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiConnectionTestLink struct {
	Method     *string `json:"method,omitempty"`
	RequestUri *string `json:"requestUri,omitempty"`
}
