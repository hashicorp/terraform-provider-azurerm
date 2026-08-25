// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DataFactoryPipelineAndTriggerName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_][^<>*#.%&:\\+?/]*$`), "invalid name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules")
}

func DataFactoryName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`), "invalid data_factory_name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules")
}

func DataFactoryManagedPrivateEndpointName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^([[:alnum:]][-._[:alnum:]]{0,78}[_[:alnum:]])$`), "invalid Data Factory Managed Private Endpoint name, must match the regular expression ^^([[:alnum:]][-._[:alnum:]]{0,78}[_[:alnum:]])$")
}

func CMKIdentityIdRequiredAtCreation(ctx context.Context, d *pluginsdk.ResourceDiff, meta interface{}) error {
	if d.Id() == "" {
		rawConfig := d.GetRawConfig().AsValueMap()

		rawCMK := rawConfig["customer_managed_key_id"]
		rawCMKIdentity := rawConfig["customer_managed_key_identity_id"]

		if !rawCMK.IsNull() && rawCMKIdentity.IsNull() {
			return fmt.Errorf("`customer_managed_key_identity_id` is required when creating a new Data Factory with `customer_managed_key_id`")
		}
	}
	return nil
}
