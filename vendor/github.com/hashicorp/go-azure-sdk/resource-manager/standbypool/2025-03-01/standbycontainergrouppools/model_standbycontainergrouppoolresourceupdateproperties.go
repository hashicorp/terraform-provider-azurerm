package standbycontainergrouppools

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StandbyContainerGroupPoolResourceUpdateProperties struct {
	ContainerGroupProperties *ContainerGroupProperties                   `json:"containerGroupProperties,omitempty"`
	ElasticityProfile        *StandbyContainerGroupPoolElasticityProfile `json:"elasticityProfile,omitempty"`
	Zones                    *zones.Schema                               `json:"zones,omitempty"`
}
