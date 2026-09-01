package securitymlanalyticssettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnomalySecurityMLAnalyticsMultiSelectObservations struct {
	Description        *string      `json:"description,omitempty"`
	Name               *string      `json:"name,omitempty"`
	Rerun              *interface{} `json:"rerun,omitempty"`
	SequenceNumber     *interface{} `json:"sequenceNumber,omitempty"`
	SupportedValues    *[]string    `json:"supportedValues,omitempty"`
	SupportedValuesKql *interface{} `json:"supportedValuesKql,omitempty"`
	Values             *[]string    `json:"values,omitempty"`
	ValuesKql          *interface{} `json:"valuesKql,omitempty"`
}
