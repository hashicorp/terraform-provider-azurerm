package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineGroupProperties struct {
	ExecutionPlacement *ExecutionPlacement `json:"executionPlacement,omitempty"`
	Exporters          []Exporter          `json:"exporters"`
	Processors         []Processor         `json:"processors"`
	ProvisioningState  *ProvisioningState  `json:"provisioningState,omitempty"`
	Receivers          []Receiver          `json:"receivers"`
	Replicas           *int64              `json:"replicas,omitempty"`
	Service            Service             `json:"service"`
	TlsConfigurations  *[]TlsConfiguration `json:"tlsConfigurations,omitempty"`
}
