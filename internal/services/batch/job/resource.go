// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package job

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource           = Resource{}
	_ sdk.ResourceWithUpdate = Resource{}
)

const ResourceName = "azurerm_batch_job"

func (r Resource) ResourceType() string {
	return ResourceName
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.JobID
}
