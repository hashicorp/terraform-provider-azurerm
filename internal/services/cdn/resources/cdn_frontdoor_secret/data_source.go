// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_secret

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Front Door (standard/premium) Secret.",
		Read:        dataSourceCdnFrontDoorSecretRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Front Door Secret.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CdnFrontDoorSecretName,
			},

			"profile_name": {
				Description:  "The name of the Front Door Profile within which the Front Door Secret exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			// Computed
			"cdn_frontdoor_profile_id": {
				Description: "The ID of the Front Door Profile within which this Front Door Secret exists.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"secret": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"customer_certificate": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"expiration_date": {
										Description: "The Key Vault Certificate expiration date.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},

									"key_vault_certificate_id": {
										Description: "The Key Vault Certificate ID.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},

									"subject_alternative_names": {
										Description: "One or more `subject alternative names` contained within the Key Vault Certificate.",
										Type:        pluginsdk.TypeList,
										Computed:    true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCdnFrontDoorSecretRead(d *pluginsdk.ResourceData, meta any) error {
	client := cdnClient.From(meta).FrontDoorSecretsClient

	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := secrets.NewSecretID(account.From(meta).SubscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.SecretName)
	d.Set("profile_name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			var customerCertificate []any
			if customerCertificate, err = flattenCdnFrontDoorSecretParameters(ctx, props.Parameters, meta); err != nil {
				return fmt.Errorf("flattening `secret`: %+v", err)
			}

			if err := d.Set("secret", customerCertificate); err != nil {
				return fmt.Errorf("setting `secret`: %+v", err)
			}
		}
	}

	return nil
}
