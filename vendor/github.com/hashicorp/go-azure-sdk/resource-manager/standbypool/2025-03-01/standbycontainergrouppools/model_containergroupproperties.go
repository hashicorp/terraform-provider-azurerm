package standbycontainergrouppools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupProperties struct {
	ContainerGroupProfile ContainerGroupProfile `json:"containerGroupProfile"`
	SubnetIds             *[]Subnet             `json:"subnetIds,omitempty"`
}
