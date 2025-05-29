// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeSelfHosted() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate,
		Read:   resourceDataFactoryIntegrationRuntimeSelfHostedRead,
		Update: resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeSelfHostedDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DataFactoryIntegrationRuntimeSelfHostedV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationruntimes.ParseIntegrationRuntimeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Self-Hosted Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"rbac_authorization": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				ForceNew: true,
				Set:      resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: integrationruntimes.ValidateIntegrationRuntimeID,
						},
					},
				},
			},

			"self_contained_interactive_authoring_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"primary_authorization_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_authorization_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	if !features.FivePointOh() {
		r.Schema["rbac_authorization"].Elem.(*pluginsdk.Resource).Schema["resource_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		}
	}

	return r
}

func resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := integrationruntimes.NewIntegrationRuntimeID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, integrationruntimes.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_self_hosted", id.ID())
		}
	}

	selfHostedIntegrationRuntime := integrationruntimes.SelfHostedIntegrationRuntime{
		Description: pointer.To(d.Get("description").(string)),
		Type:        integrationruntimes.IntegrationRuntimeTypeSelfHosted,
		TypeProperties: &integrationruntimes.SelfHostedIntegrationRuntimeTypeProperties{
			SelfContainedInteractiveAuthoringEnabled: pointer.To(d.Get("self_contained_interactive_authoring_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("rbac_authorization"); ok {
		if linkedInfo := expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesLinkedInfo(v.(*pluginsdk.Set).List()); linkedInfo != nil {
			selfHostedIntegrationRuntime.TypeProperties.LinkedInfo = linkedInfo
		}
	}

	integrationRuntime := integrationruntimes.IntegrationRuntimeResource{
		Name:       pointer.To(id.IntegrationRuntimeName),
		Properties: selfHostedIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id, integrationRuntime, integrationruntimes.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryIntegrationRuntimeSelfHostedRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeSelfHostedRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)

	resp, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("data_factory_id", dataFactoryId.ID())

	if model := resp.Model; model != nil {
		runTime, ok := model.Properties.(integrationruntimes.SelfHostedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("asserting `IntegrationRuntime` as `SelfHostedIntegrationRuntime` for %s", *id)
		}

		d.Set("description", runTime.Description)

		if props := runTime.TypeProperties; props != nil {
			d.Set("self_contained_interactive_authoring_enabled", pointer.From(props.SelfContainedInteractiveAuthoringEnabled))
			rbacAuthorization, ok := props.LinkedInfo.(integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization)
			if ok {
				if err := d.Set("rbac_authorization", pluginsdk.NewSet(resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash, flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(rbacAuthorization))); err != nil {
					return fmt.Errorf("setting `rbac_authorization`: %#v", err)
				}
			}

			// The ListAuthenticationKeys on integration runtime type Linked is not supported.
			// Only skip the call to ListAuthKeys if the linkedInfo is valid.
			if _, ok := props.LinkedInfo.(integrationruntimes.RawLinkedIntegrationRuntimeTypeImpl); !ok {
				return nil
			}
		}
	}

	keyResp, keyErr := client.ListAuthKeys(ctx, *id)
	if keyErr != nil {
		if response.WasNotFound(keyResp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Authorization Keys for %s: %+v", *id, keyErr)
	}

	if keyResp.Model == nil {
		return fmt.Errorf("retrieving Authorization Keys for %s: `model` was nil", *id)
	}
	keyModel := keyResp.Model

	d.Set("primary_authorization_key", keyModel.AuthKey1)
	d.Set("secondary_authorization_key", keyModel.AuthKey2)

	return nil
}

func resourceDataFactoryIntegrationRuntimeSelfHostedDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesLinkedInfo(input []interface{}) *integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization {
	if len(input) == 0 {
		return nil
	}

	rbacConfig := input[0].(map[string]interface{})
	rbac := rbacConfig["resource_id"].(string)
	return &integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization{
		ResourceId:        rbac,
		AuthorizationType: string(helper.AuthorizationTypeRBAC),
	}
}

func flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(input integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization) []interface{} {
	result := make(map[string]interface{})
	result["resource_id"] = input.ResourceId

	return []interface{}{result}
}

func resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["resource_id"]; ok {
			if !features.FivePointOh() {
				buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
			} else {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}
