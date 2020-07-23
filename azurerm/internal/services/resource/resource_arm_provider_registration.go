package resource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceproviders"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmResourceProviderRegistration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmResourceProviderCreate,
		Read:   resourceArmResourceProviderRead,
		Delete: resourceArmResourceProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmResourceProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ProvidersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceProviderNamespace := d.Get("name").(string)
	expand := ""

	isSkipResourceProviderRegistration := meta.(*clients.Client).Account.SkipResourceProviderRegistration

	if !isSkipResourceProviderRegistration {
		for resourceProvider := range resourceproviders.Required() {
			if strings.EqualFold(resourceProviderNamespace, resourceProvider) {
				return fmt.Errorf(`The Resource Provider %q is automatically registered by Terraform. To manage this Resource Provider Registration with Terraform you need to opt-out of Automatic Resource Provider Registration (by setting 'skip_provider_registration' to 'true' in the Provider block) to avoid conflicting with Terraform`, resourceProviderNamespace)
			}
		}
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceProviderNamespace, expand)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Resource Provider Namespace %q: %+v",
					resourceProviderNamespace, err)
			}
		}

		if existing.ID != nil && existing.RegistrationState != nil && *existing.RegistrationState == "Registered" {
			return tf.ImportAsExistsError("azurerm_resource_provider_registration", *existing.ID)
		}
	}

	if _, err := client.Register(ctx, resourceProviderNamespace); err != nil {
		return fmt.Errorf("Cannot register provider %q with Azure Resource Manager: %+v",
			resourceProviderNamespace, err)
	}

	read, err := client.Get(ctx, resourceProviderNamespace, expand)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Resource Provider Namespace %q", resourceProviderNamespace)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Registered"},
		Refresh:    resourceProviderNamespaceRegisterRefreshFunc(ctx, client, resourceProviderNamespace, expand),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Resource Provider Namespace %q to be registered: %s",
			resourceProviderNamespace, err)
	}

	d.SetId(*read.ID)

	return resourceArmResourceProviderRead(d, meta)
}

func resourceArmResourceProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ProvidersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceProviderNamespace := id.Provider
	expand := ""

	isSkipResourceProviderRegistration := meta.(*clients.Client).Account.SkipResourceProviderRegistration

	if !isSkipResourceProviderRegistration {
		for resourceProvider := range resourceproviders.Required() {
			if strings.EqualFold(resourceProviderNamespace, resourceProvider) {
				return fmt.Errorf(`The Resource Provider %q is automatically registered by Terraform. To manage this Resource Provider Registration with Terraform you need to opt-out of Automatic Resource Provider Registration (by setting 'skip_provider_registration' to 'true' in the Provider block) to avoid conflicting with Terraform`, resourceProviderNamespace)
			}
		}
	}

	resp, err := client.Get(ctx, resourceProviderNamespace, expand)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Resource Provider Namespace '%s' was not found",
				resourceProviderNamespace)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Resource Provider Namespace %q: %+v",
			resourceProviderNamespace, err)
	}

	if resp.RegistrationState != nil && *resp.RegistrationState != "Registered" {
		log.Printf("[WARN] Resource Provider Namespace '%s' was not registered",
			resourceProviderNamespace)
		d.SetId("")
		return nil
	}

	d.Set("name", resp.Namespace)

	return nil
}

func resourceArmResourceProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ProvidersClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceProviderNamespace := id.Provider
	expand := ""

	isSkipResourceProviderRegistration := meta.(*clients.Client).Account.SkipResourceProviderRegistration

	if !isSkipResourceProviderRegistration {
		for resourceProvider := range resourceproviders.Required() {
			if strings.EqualFold(resourceProviderNamespace, resourceProvider) {
				return fmt.Errorf(`The Resource Provider %q is automatically registered by Terraform. To manage this Resource Provider Registration with Terraform you need to opt-out of Automatic Resource Provider Registration (by setting 'skip_provider_registration' to 'true' in the Provider block) to avoid conflicting with Terraform`, resourceProviderNamespace)
			}
		}
	}

	if _, err := client.Unregister(ctx, resourceProviderNamespace); err != nil {
		return fmt.Errorf("Error unregistering Resource Provider Namespace %q: %+v",
			resourceProviderNamespace, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Unregistered"},
		Refresh:    resourceProviderNamespaceUnregisterRefreshFunc(ctx, client, resourceProviderNamespace, expand),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Resource Provider Namespace %q to be unregistered: %s",
			resourceProviderNamespace, err)
	}

	return nil
}

func resourceProviderNamespaceRegisterRefreshFunc(ctx context.Context, client *resources.ProvidersClient, resourceProviderNamespace string, expand string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceProviderNamespace, expand)
		if err != nil {
			return resp, "Failed", err
		}
		if resp.RegistrationState != nil && *resp.RegistrationState == "Registered" {
			return resp, "Registered", nil
		}
		return resp, "Processing", nil
	}
}

func resourceProviderNamespaceUnregisterRefreshFunc(ctx context.Context, client *resources.ProvidersClient, resourceProviderNamespace string, expand string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceProviderNamespace, expand)
		if err != nil {
			return resp, "Failed", err
		}
		if resp.RegistrationState != nil && *resp.RegistrationState == "Unregistered" {
			return resp, "Unregistered", nil
		}
		return resp, "Processing", nil
	}
}
