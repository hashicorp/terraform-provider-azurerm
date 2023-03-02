package watcher

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatcherUpdateParameters struct {
	Name       *string                  `json:"name,omitempty"`
	Properties *WatcherUpdateProperties `json:"properties,omitempty"`
}
