package keys

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyReleasePolicy struct {
	ContentType *string `json:"contentType,omitempty"`
	Data        *string `json:"data,omitempty"`
	Immutable   *bool   `json:"immutable,omitempty"`
}
