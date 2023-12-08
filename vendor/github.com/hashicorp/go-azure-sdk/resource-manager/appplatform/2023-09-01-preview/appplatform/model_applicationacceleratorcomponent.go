package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationAcceleratorComponent struct {
	Instances        *[]ApplicationAcceleratorInstance       `json:"instances,omitempty"`
	Name             *string                                 `json:"name,omitempty"`
	ResourceRequests *ApplicationAcceleratorResourceRequests `json:"resourceRequests,omitempty"`
}
