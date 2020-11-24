package appplatform

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2019-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/appplatform/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSpringCloudApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSpringCloudAppCreateUpdate,
		Read:   resourceArmSpringCloudAppRead,
		Update: resourceArmSpringCloudAppCreateUpdate,
		Delete: resourceArmSpringCloudAppDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SpringCloudAppID(id)
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudAppName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceName,
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							// other two enum: 'UserAssigned', 'SystemAssignedUserAssigned' are not supported for now
							ValidateFunc: validation.StringInSlice([]string{
								string(appplatform.SystemAssigned),
							}, false),
						},

						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmSpringCloudAppCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	servicesClient := meta.(*clients.Client).AppPlatform.ServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)

	serviceResp, err := servicesClient.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return fmt.Errorf("unable to retrieve Spring Cloud Service %q (Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	resourceId := parse.NewSpringCloudAppID(subscriptionId, resourceGroup, serviceName, name).ID("")
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_app", resourceId)
		}
	}

	app := appplatform.AppResource{
		Location: serviceResp.Location,
		Identity: expandSpringCloudAppIdentity(d.Get("identity").([]interface{})),
	}
	future, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, app)
	if err != nil {
		return fmt.Errorf("creating/update Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceArmSpringCloudAppRead(d, meta)
}

func resourceArmSpringCloudAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Spring Cloud App %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", id.AppName, id.SpringName, id.ResourceGroup, err)
	}

	d.Set("name", id.AppName)
	d.Set("service_name", id.SpringName)
	d.Set("resource_group_name", id.ResourceGroup)

	if err := d.Set("identity", flattenSpringCloudAppIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %s", err)
	}

	return nil
}

func resourceArmSpringCloudAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.AppsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudAppID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.AppName); err != nil {
		return fmt.Errorf("deleting Spring Cloud App %q (Spring Cloud Service %q / Resource Group %q): %+v", id.AppName, id.SpringName, id.ResourceGroup, err)
	}

	return nil
}

func expandSpringCloudAppIdentity(input []interface{}) *appplatform.ManagedIdentityProperties {
	if len(input) == 0 || input[0] == nil {
		return &appplatform.ManagedIdentityProperties{
			Type: appplatform.None,
		}
	}
	identity := input[0].(map[string]interface{})
	return &appplatform.ManagedIdentityProperties{
		Type: appplatform.ManagedIdentityType(identity["type"].(string)),
	}
}

func flattenSpringCloudAppIdentity(identity *appplatform.ManagedIdentityProperties) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	principalId := ""
	if identity.PrincipalID != nil {
		principalId = *identity.PrincipalID
	}

	tenantId := ""
	if identity.TenantID != nil {
		tenantId = *identity.TenantID
	}

	return []interface{}{
		map[string]interface{}{
			"principal_id": principalId,
			"tenant_id":    tenantId,
			"type":         string(identity.Type),
		},
	}
}
