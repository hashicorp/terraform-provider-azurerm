// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),
		"default_maximum_expiration_time_in_days": &schema.Schema{
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 365),
			Default:      7,
		},
		"partner_authorization": &schema.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"partner_registration_id": &schema.Schema{
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
					"partner_name": &schema.Schema{
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"authorization_expiration_time_in_utc": &schema.Schema{
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			config := metadata.ResourceDiff.GetRawConfig()
			if config.IsNull() || !config.IsKnown() {
				return nil
			}

			return validatePartnerAuthorizations(config.GetAttr("partner_authorization"))
		},
	}
}

func validatePartnerAuthorizations(partners cty.Value) error {
	if partners.IsNull() || !partners.IsKnown() {
		return nil
	}

	for iterator := partners.ElementIterator(); iterator.Next(); {
		index, partner := iterator.Element()
		if !partner.IsKnown() || partner.IsNull() {
			continue
		}

		if partner.GetAttr("partner_registration_id").IsNull() && partner.GetAttr("partner_name").IsNull() {
			blockIndex, _ := index.AsBigFloat().Int64()
			return fmt.Errorf("at least one of `partner_registration_id` or `partner_name` must be provided in `partner_authorization.%d`", blockIndex)
		}
	}

	return nil
}
