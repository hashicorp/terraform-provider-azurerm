// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryIntegrationRuntimeSelfHosted() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate,
		Read:   resourceDataFactoryIntegrationRuntimeSelfHostedRead,
		Update: resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeSelfHostedDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IntegrationRuntimeID(id)
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
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
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
}

func resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	subscriptionId := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewIntegrationRuntimeID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_self_hosted", id.ID())
		}
	}

	description := d.Get("description").(string)

	selfHostedIntegrationRuntime := datafactory.SelfHostedIntegrationRuntime{
		Description: &description,
		Type:        datafactory.TypeBasicIntegrationRuntimeTypeSelfHosted,
	}

	properties := expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypeProperties(d)
	if properties != nil {
		selfHostedIntegrationRuntime.SelfHostedIntegrationRuntimeTypeProperties = properties
	}

	basicIntegrationRuntime, _ := selfHostedIntegrationRuntime.AsBasicIntegrationRuntime()

	integrationRuntime := datafactory.IntegrationRuntimeResource{
		Name:       &id.Name,
		Properties: basicIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, integrationRuntime, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Self-Hosted %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryIntegrationRuntimeSelfHostedRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeSelfHostedRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	selfHostedIntegrationRuntime, convertSuccess := resp.Properties.AsSelfHostedIntegrationRuntime()

	if !convertSuccess {
		return fmt.Errorf("converting Integration Runtime to Self-Hosted %s", *id)
	}

	if selfHostedIntegrationRuntime.Description != nil {
		d.Set("description", selfHostedIntegrationRuntime.Description)
	}

	if props := selfHostedIntegrationRuntime.SelfHostedIntegrationRuntimeTypeProperties; props != nil {
		// LinkedInfo BasicLinkedIntegrationRuntimeType
		if linkedInfo := props.LinkedInfo; linkedInfo != nil {
			rbacAuthorization, _ := linkedInfo.AsLinkedIntegrationRuntimeRbacAuthorization()
			if rbacAuthorization != nil {
				if err := d.Set("rbac_authorization", pluginsdk.NewSet(resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash, flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(rbacAuthorization))); err != nil {
					return fmt.Errorf("setting `rbac_authorization`: %#v", err)
				}
			}
		}
		return nil
	}

	respKey, errKey := client.ListAuthKeys(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if errKey != nil {
		if utils.ResponseWasNotFound(respKey.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Auth Keys for Data Factory Self-Hosted %s: %+v", *id, errKey)
	}

	d.Set("primary_authorization_key", respKey.AuthKey1)
	d.Set("secondary_authorization_key", respKey.AuthKey2)

	return nil
}

func resourceDataFactoryIntegrationRuntimeSelfHostedDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}
	return nil
}

func expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypeProperties(d *pluginsdk.ResourceData) *datafactory.SelfHostedIntegrationRuntimeTypeProperties {
	if _, ok := d.GetOk("rbac_authorization"); ok {
		rbacAuthorization := d.Get("rbac_authorization").(*pluginsdk.Set).List()
		rbacConfig := rbacAuthorization[0].(map[string]interface{})
		rbac := rbacConfig["resource_id"].(string)
		linkedInfo := &datafactory.SelfHostedIntegrationRuntimeTypeProperties{
			LinkedInfo: &datafactory.LinkedIntegrationRuntimeRbacAuthorization{
				ResourceID:        &rbac,
				AuthorizationType: datafactory.AuthorizationTypeRBAC,
			},
		}
		return linkedInfo
	}
	return nil
}

func flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(input *datafactory.LinkedIntegrationRuntimeRbacAuthorization) []interface{} {
	result := make(map[string]interface{})
	result["resource_id"] = *input.ResourceID

	return []interface{}{result}
}

func resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["resource_id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}
