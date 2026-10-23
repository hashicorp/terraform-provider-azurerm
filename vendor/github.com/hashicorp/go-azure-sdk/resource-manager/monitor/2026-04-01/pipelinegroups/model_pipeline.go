package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Pipeline struct {
	Exporters  []string     `json:"exporters"`
	Name       string       `json:"name"`
	Processors *[]string    `json:"processors,omitempty"`
	Receivers  []string     `json:"receivers"`
	Type       PipelineType `json:"type"`
}
