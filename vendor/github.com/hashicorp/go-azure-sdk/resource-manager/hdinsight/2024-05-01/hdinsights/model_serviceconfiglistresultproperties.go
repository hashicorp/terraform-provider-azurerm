package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceConfigListResultProperties struct {
	ComponentName string                                         `json:"componentName"`
	Content       *string                                        `json:"content,omitempty"`
	CustomKeys    *map[string]string                             `json:"customKeys,omitempty"`
	DefaultKeys   *map[string]ServiceConfigListResultValueEntity `json:"defaultKeys,omitempty"`
	FileName      string                                         `json:"fileName"`
	Path          *string                                        `json:"path,omitempty"`
	ServiceName   string                                         `json:"serviceName"`
	Type          *string                                        `json:"type,omitempty"`
}
