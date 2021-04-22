package migration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func RegistryV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    registrySchemaForV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: registryUpgradeV0ToV1,
		Version: 0,
	}
}

func RegistryV1ToV2() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    registrySchemaForV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: registryUpgradeV1ToV2,
		Version: 1,
	}
}

func registrySchemaForV0AndV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// lintignore:S018
			"storage_account": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},

			"login_server": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_username": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_password": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func registryUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState["sku"] = "Basic"
	return rawState, nil
}

func registryUpgradeV1ToV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	// Basic's been renamed Classic to allow for "ManagedBasic" ¯\_(ツ)_/¯
	rawState["sku"] = "Classic"

	storageAccountId := ""
	if v, ok := rawState["storage_account"]; ok {
		raw := v.(*schema.Set).List()
		rawVals := raw[0].(map[string]interface{})
		storageAccountName := rawVals["name"].(string)

		client := meta.(*clients.Client).Storage.AccountsClient
		ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, time.Minute*5)
		defer cancel()

		accounts, err := client.ListComplete(ctx)
		if err != nil {
			return rawState, fmt.Errorf("listing storage accounts")
		}

		for accounts.NotDone() {
			account := accounts.Value()
			if strings.EqualFold(*account.Name, storageAccountName) {
				storageAccountId = *account.ID
				break
			}

			if err := accounts.NextWithContext(ctx); err != nil {
				return rawState, fmt.Errorf("retrieving accounts: %+v", err)
			}
		}
	}

	if storageAccountId == "" {
		return rawState, fmt.Errorf("unable to determine storage account ID")
	}

	return rawState, nil
}
