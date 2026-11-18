package appserviceplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerFarmInstance struct {
	IPAddress    *string `json:"ipAddress,omitempty"`
	InstanceName *string `json:"instanceName,omitempty"`
	Status       *string `json:"status,omitempty"`
}
