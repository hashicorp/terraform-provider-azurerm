package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RefreshConfiguration struct {
	DateFormat      *string                  `json:"dateFormat,omitempty"`
	PathPattern     *string                  `json:"pathPattern,omitempty"`
	RefreshInterval *string                  `json:"refreshInterval,omitempty"`
	RefreshType     *UpdatableUdfRefreshType `json:"refreshType,omitempty"`
	TimeFormat      *string                  `json:"timeFormat,omitempty"`
}
