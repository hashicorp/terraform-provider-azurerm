package sqldedicatedgateway

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceResourceCreateUpdateProperties struct {
	InstanceCount *int64       `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize `json:"instanceSize,omitempty"`
	ServiceType   *ServiceType `json:"serviceType,omitempty"`
}
