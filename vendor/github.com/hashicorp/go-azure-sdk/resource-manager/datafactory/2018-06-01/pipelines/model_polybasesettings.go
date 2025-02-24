package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolybaseSettings struct {
	RejectSampleValue *int64                      `json:"rejectSampleValue,omitempty"`
	RejectType        *PolybaseSettingsRejectType `json:"rejectType,omitempty"`
	RejectValue       *int64                      `json:"rejectValue,omitempty"`
	UseTypeDefault    *bool                       `json:"useTypeDefault,omitempty"`
}
