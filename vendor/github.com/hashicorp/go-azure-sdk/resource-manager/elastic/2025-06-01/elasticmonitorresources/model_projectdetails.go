package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectDetails struct {
	ConfigurationType *ConfigurationType `json:"configurationType,omitempty"`
	ProjectType       *ProjectType       `json:"projectType,omitempty"`
}
