package securitymlanalyticssettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnomalySecurityMLAnalyticsThresholdObservations struct {
	Description    *string      `json:"description,omitempty"`
	Maximum        *string      `json:"maximum,omitempty"`
	Minimum        *string      `json:"minimum,omitempty"`
	Name           *string      `json:"name,omitempty"`
	Rerun          *interface{} `json:"rerun,omitempty"`
	SequenceNumber *interface{} `json:"sequenceNumber,omitempty"`
	Value          *string      `json:"value,omitempty"`
}
