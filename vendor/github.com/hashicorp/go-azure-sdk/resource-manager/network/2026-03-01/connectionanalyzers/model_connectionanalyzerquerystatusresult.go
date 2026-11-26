package connectionanalyzers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionAnalyzerQueryStatusResult struct {
	ConnectionAnalyzerStatus   *ConnectionAnalyzerStatus    `json:"connectionAnalyzerStatus,omitempty"`
	DiagnosticOperationResults *[]DiagnosticOperationResult `json:"diagnosticOperationResults,omitempty"`
	Error                      *string                      `json:"error,omitempty"`
	ExpiryInUtc                *string                      `json:"expiryInUtc,omitempty"`
	Id                         *string                      `json:"id,omitempty"`
	OutputStoragePath          *string                      `json:"outputStoragePath,omitempty"`
}

func (o *ConnectionAnalyzerQueryStatusResult) GetExpiryInUtcAsTime() (*time.Time, error) {
	if o.ExpiryInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectionAnalyzerQueryStatusResult) SetExpiryInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryInUtc = &formatted
}
