package experiments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Selector struct {
	Filter *Filter      `json:"filter,omitempty"`
	Id     string       `json:"id"`
	Type   SelectorType `json:"type"`
}
