package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineGroupPropertiesUpdate struct {
	ExecutionPlacement *ExecutionPlacement `json:"executionPlacement,omitempty"`
	Exporters          *[]Exporter         `json:"exporters,omitempty"`
	Processors         *[]Processor        `json:"processors,omitempty"`
	Receivers          *[]Receiver         `json:"receivers,omitempty"`
	Replicas           *int64              `json:"replicas,omitempty"`
	Service            *ServiceUpdate      `json:"service,omitempty"`
	TlsConfigurations  *[]TlsConfiguration `json:"tlsConfigurations,omitempty"`
}
