// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type CognitiveAccountProjectModel struct {
	CognitiveAccountId string                                     `tfschema:"cognitive_account_id"`
	Default            bool                                       `tfschema:"default"`
	Description        string                                     `tfschema:"description"`
	DisplayName        string                                     `tfschema:"display_name"`
	Endpoints          map[string]string                          `tfschema:"endpoints"`
	Identity           []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location           string                                     `tfschema:"location"`
	Name               string                                     `tfschema:"name"`
	Tags               map[string]string                          `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &CognitiveAccountProjectModel{}
}
