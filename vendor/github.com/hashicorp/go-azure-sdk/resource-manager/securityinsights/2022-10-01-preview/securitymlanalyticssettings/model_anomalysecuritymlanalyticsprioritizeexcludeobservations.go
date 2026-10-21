package securitymlanalyticssettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnomalySecurityMLAnalyticsPrioritizeExcludeObservations struct {
	DataType       *interface{} `json:"dataType,omitempty"`
	Description    *string      `json:"description,omitempty"`
	Exclude        *string      `json:"exclude,omitempty"`
	Name           *string      `json:"name,omitempty"`
	Prioritize     *string      `json:"prioritize,omitempty"`
	Rerun          *interface{} `json:"rerun,omitempty"`
	SequenceNumber *interface{} `json:"sequenceNumber,omitempty"`
}
