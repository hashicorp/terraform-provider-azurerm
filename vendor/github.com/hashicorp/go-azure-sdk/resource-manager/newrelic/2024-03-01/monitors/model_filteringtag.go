package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilteringTag struct {
	Action *TagAction `json:"action,omitempty"`
	Name   *string    `json:"name,omitempty"`
	Value  *string    `json:"value,omitempty"`
}
