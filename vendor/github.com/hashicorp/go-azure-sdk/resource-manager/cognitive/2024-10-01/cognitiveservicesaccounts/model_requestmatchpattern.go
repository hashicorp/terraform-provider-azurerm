package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequestMatchPattern struct {
	Method *string `json:"method,omitempty"`
	Path   *string `json:"path,omitempty"`
}
