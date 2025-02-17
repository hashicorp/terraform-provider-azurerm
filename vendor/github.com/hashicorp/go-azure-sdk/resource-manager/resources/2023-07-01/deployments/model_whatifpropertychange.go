package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WhatIfPropertyChange struct {
	After              *interface{}            `json:"after,omitempty"`
	Before             *interface{}            `json:"before,omitempty"`
	Children           *[]WhatIfPropertyChange `json:"children,omitempty"`
	Path               string                  `json:"path"`
	PropertyChangeType PropertyChangeType      `json:"propertyChangeType"`
}
