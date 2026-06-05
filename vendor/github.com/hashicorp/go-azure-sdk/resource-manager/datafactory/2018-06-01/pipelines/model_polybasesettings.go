package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolybaseSettings struct {
	RejectSampleValue *interface{}                `json:"rejectSampleValue,omitempty"`
	RejectType        *PolybaseSettingsRejectType `json:"rejectType,omitempty"`
	RejectValue       *interface{}                `json:"rejectValue,omitempty"`
	UseTypeDefault    *interface{}                `json:"useTypeDefault,omitempty"`
}
