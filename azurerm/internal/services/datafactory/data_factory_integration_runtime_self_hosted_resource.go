package datafactory

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryIntegrationRuntimeSelfHosted() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate,
		Read:   resourceDataFactoryIntegrationRuntimeSelfHostedRead,
		Update: resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate,
		Delete: resourceDataFactoryIntegrationRuntimeSelfHostedDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.IntegrationRuntimeID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Self-Hosted Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
				),
			},

			"data_factory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"rbac_authorization": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"auth_key_1": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"auth_key_2": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDataFactoryIntegrationRuntimeSelfHostedCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	factoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, factoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Self-Hosted Integration Runtime %q (Resource Group %q, Data Factory %q): %s", name, resourceGroup, factoryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_integration_runtime_self_hosted", *existing.ID)
		}
	}

	description := d.Get("description").(string)

	selfHostedIntegrationRuntime := datafactory.SelfHostedIntegrationRuntime{
		Description: &description,
		Type:        datafactory.TypeSelfHosted,
	}

	properties := expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypeProperties(d)
	if properties != nil {
		selfHostedIntegrationRuntime.SelfHostedIntegrationRuntimeTypeProperties = properties
	}

	basicIntegrationRuntime, _ := selfHostedIntegrationRuntime.AsBasicIntegrationRuntime()

	integrationRuntime := datafactory.IntegrationRuntimeResource{
		Name:       &name,
		Properties: basicIntegrationRuntime,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, factoryName, name, integrationRuntime, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Self-Hosted Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Self-Hosted Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Self-Hosted Integration Runtime %q (Resource Group %q, Data Factory %q) ID", name, resourceGroup, factoryName)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryIntegrationRuntimeSelfHostedRead(d, meta)
}

func resourceDataFactoryIntegrationRuntimeSelfHostedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	factoryName := id.FactoryName
	name := id.Name

	resp, err := client.Get(ctx, resourceGroup, factoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Self-Hosted Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
	}

	d.Set("name", name)
	d.Set("data_factory_name", factoryName)
	d.Set("resource_group_name", resourceGroup)

	selfHostedIntegrationRuntime, convertSuccess := resp.Properties.AsSelfHostedIntegrationRuntime()

	if !convertSuccess {
		return fmt.Errorf("Error converting integration runtime to Self-Hosted integration runtime %q (Resource Group %q, Data Factory %q)", name, resourceGroup, factoryName)
	}

	if selfHostedIntegrationRuntime.Description != nil {
		d.Set("description", selfHostedIntegrationRuntime.Description)
	}

	if props := selfHostedIntegrationRuntime.SelfHostedIntegrationRuntimeTypeProperties; props != nil {
		// LinkedInfo BasicLinkedIntegrationRuntimeType
		if linkedInfo := props.LinkedInfo; linkedInfo != nil {
			rbacAuthorization, _ := linkedInfo.AsLinkedIntegrationRuntimeRbacAuthorization()
			if rbacAuthorization != nil {
				if err := d.Set("rbac_authorization", schema.NewSet(resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash, flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(rbacAuthorization))); err != nil {
					return fmt.Errorf("Error setting `rbac_authorization`: %#v", err)
				}
			}
		}
		return nil
	}

	respKey, errKey := client.ListAuthKeys(ctx, resourceGroup, factoryName, name)
	if errKey != nil {
		if utils.ResponseWasNotFound(respKey.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Self-Hosted Integration Runtime %q Auth Keys (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, errKey)
	}

	d.Set("auth_key_1", respKey.AuthKey1)
	d.Set("auth_key_1", respKey.AuthKey2)

	return nil
}

func resourceDataFactoryIntegrationRuntimeSelfHostedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	factoryName := id.FactoryName
	name := id.Name

	response, err := client.Delete(ctx, resourceGroup, factoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory SelfHosted Integration Runtime %q (Resource Group %q, Data Factory %q): %+v", name, resourceGroup, factoryName, err)
		}
	}
	return nil
}

func expandAzureRmDataFactoryIntegrationRuntimeSelfHostedTypeProperties(d *schema.ResourceData) *datafactory.SelfHostedIntegrationRuntimeTypeProperties {
	if _, ok := d.GetOk("rbac_authorization"); ok {
		rbacAuthorization := d.Get("rbac_authorization").(*schema.Set).List()
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

	return hashcode.String(buf.String())
}
