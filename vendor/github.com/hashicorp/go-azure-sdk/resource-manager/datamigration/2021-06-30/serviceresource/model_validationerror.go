package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidationError struct {
	Severity *Severity `json:"severity,omitempty"`
	Text     *string   `json:"text,omitempty"`
}
