package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeInstanceContainer struct {
	Autosave    *Autosave                       `json:"autosave,omitempty"`
	Environment *ComputeInstanceEnvironmentInfo `json:"environment,omitempty"`
	Gpu         *string                         `json:"gpu,omitempty"`
	Name        *string                         `json:"name,omitempty"`
	Network     *Network                        `json:"network,omitempty"`
	Services    *[]interface{}                  `json:"services,omitempty"`
}
