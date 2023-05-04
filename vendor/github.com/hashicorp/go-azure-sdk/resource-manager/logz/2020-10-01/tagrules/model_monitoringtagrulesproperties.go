package tagrules

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitoringTagRulesProperties struct {
	LogRules          *LogRules              `json:"logRules,omitempty"`
	ProvisioningState *ProvisioningState     `json:"provisioningState,omitempty"`
	SystemData        *systemdata.SystemData `json:"systemData,omitempty"`
}
