package connectionanalyzers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionAnalyzerProperties struct {
	Destination                  ConnectionAnalyzerEndpoint    `json:"destination"`
	DiagnosticOperations         []DiagnosticOperation         `json:"diagnosticOperations"`
	DiagnosticOperationsSettings *DiagnosticOperationsSettings `json:"diagnosticOperationsSettings,omitempty"`
	ExpiryInDays                 *int64                        `json:"expiryInDays,omitempty"`
	OutputSettings               *OutputSettings               `json:"outputSettings,omitempty"`
	ProtocolSettings             *ProtocolSettings             `json:"protocolSettings,omitempty"`
	ProvisioningState            *ProvisioningState            `json:"provisioningState,omitempty"`
	RequestTime                  *string                       `json:"requestTime,omitempty"`
	Source                       ConnectionAnalyzerEndpoint    `json:"source"`
	Status                       *ConnectionAnalyzerStatus     `json:"status,omitempty"`
}

func (o *ConnectionAnalyzerProperties) GetRequestTimeAsTime() (*time.Time, error) {
	if o.RequestTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RequestTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectionAnalyzerProperties) SetRequestTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RequestTime = &formatted
}
