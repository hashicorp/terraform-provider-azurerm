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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceproviders"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceResourceProviderRegistration() *schema.Resource {
	return &schema.Resource{
		Create: resourceResourceProviderCreate,
		Read:   resourceResourceProviderRead,
		Delete: resourceResourceProviderDelete,
		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.ResourceProviderID(id)
			return err
		}, importResourceProviderRegistration),

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
				ValidateFunc: resourceproviders.EnhancedValidate,
			},
		},
	}
}

func resourceResourceProviderCreate(d *schema.ResourceData, meta interface{}) error {
	account := meta.(*clients.Client).Account

	client := meta.(*clients.Client).Resource.ProvidersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewResourceProviderID(account.SubscriptionId, d.Get("name").(string))
	if err := checkIfResourceProviderManagedByTerraform(resourceId.ResourceProvider, account); err != nil {
		return err
	}

	provider, err := client.Get(ctx, resourceId.ResourceProvider, "")
	if err != nil {
		if utils.ResponseWasNotFound(provider.Response) {
			return fmt.Errorf("the Resource Provider %q was not found", resourceId.ResourceProvider)
		}

		return fmt.Errorf("retrieving Resource Provider %q: %+v", resourceId.ResourceProvider, err)
	}
	if provider.RegistrationState == nil {
		return fmt.Errorf("retrieving Resource Provider %q: `registrationState` was nil", resourceId.ResourceProvider)
	}

	if strings.EqualFold(*provider.RegistrationState, "Registered") {
		return tf.ImportAsExistsError("azurerm_resource_provider_registration", resourceId.ID(""))
	}

	log.Printf("[DEBUG] Registering Resource Provider %q..", resourceId.ResourceProvider)
	if _, err := client.Register(ctx, resourceId.ResourceProvider); err != nil {
		return fmt.Errorf("registering Resource Provider %q: %+v", resourceId.ResourceProvider, err)
	}

	log.Printf("[DEBUG] Waiting for Resource Provider %q to finish registering..", resourceId.ResourceProvider)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Processing"},
		Target:       []string{"Registered"},
		Refresh:      resourceProviderNamespaceRegisterRefreshFunc(ctx, client, resourceId.ResourceProvider),
		MinTimeout:   15 * time.Second,
		PollInterval: 30 * time.Second,
		Timeout:      d.Timeout(schema.TimeoutCreate),
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Resource Provider Namespace %q to be registered: %s", resourceId.ResourceProvider, err)
	}
	log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId.ResourceProvider)

	d.SetId(resourceId.ID(""))
	return resourceResourceProviderRead(d, meta)
}

func resourceResourceProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ProvidersClient
	account := meta.(*clients.Client).Account
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	if err := checkIfResourceProviderManagedByTerraform(id.ResourceProvider, account); err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceProvider, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Resource Provider %q was not found - removing from state!", id.ResourceProvider)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Resource Provider %q: %+v", id.ResourceProvider, err)
	}

	if resp.RegistrationState != nil && !strings.EqualFold(*resp.RegistrationState, "Registered") {
		log.Printf("[WARN] Resource Provider Namespace '%s' was not registered", id.ResourceProvider)
		d.SetId("")
		return nil
	}

	d.Set("name", id.ResourceProvider)

	return nil
}

func resourceResourceProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Resource.ProvidersClient
	account := meta.(*clients.Client).Account
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	if err := checkIfResourceProviderManagedByTerraform(id.ResourceProvider, account); err != nil {
		return err
	}

	if _, err := client.Unregister(ctx, id.ResourceProvider); err != nil {
		return fmt.Errorf("unregistering Resource Provider %q: %+v", id.ResourceProvider, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Processing"},
		Target:     []string{"Unregistered"},
		Refresh:    resourceProviderNamespaceUnregisterRefreshFunc(ctx, client, id.ResourceProvider),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Resource Provider %q to become unregistered: %+v", id.ResourceProvider, err)
	}

	return nil
}

func resourceProviderNamespaceRegisterRefreshFunc(ctx context.Context, client *resources.ProvidersClient, resourceProviderNamespace string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceProviderNamespace, "")
		if err != nil {
			return resp, "Failed", err
		}

		if resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Registered") {
			return resp, "Registered", nil
		}

		return resp, "Processing", nil
	}
}

func resourceProviderNamespaceUnregisterRefreshFunc(ctx context.Context, client *resources.ProvidersClient, resourceProvider string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, resourceProvider, "")
		if err != nil {
			return resp, "Failed", err
		}

		if resp.RegistrationState != nil && strings.EqualFold(*resp.RegistrationState, "Unregistered") {
			return resp, "Unregistered", nil
		}

		return resp, "Processing", nil
	}
}

func checkIfResourceProviderManagedByTerraform(name string, account *clients.ResourceManagerAccount) error {
	if account.SkipResourceProviderRegistration {
		return nil
	}

	for resourceProvider := range resourceproviders.Required() {
		if resourceProvider == name {
			fmtStr := `The Resource Provider %q is automatically registered by Terraform.

To manage this Resource Provider Registration with Terraform you need to opt-out
of Automatic Resource Provider Registration (by setting 'skip_provider_registration'
to 'true' in the Provider block) to avoid conflicting with Terraform.`
			return fmt.Errorf(fmtStr, name)
		}
	}

	return nil
}

func importResourceProviderRegistration(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	id, err := parse.ResourceProviderID(d.Id())
	if err != nil {
		return nil, err
	}

	client := meta.(*clients.Client).Resource.ProvidersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	provider, err := client.Get(ctx, id.ResourceProvider, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Resource Provider %q: %+v", id.ResourceProvider, err)
	}

	if provider.Namespace == nil {
		return nil, fmt.Errorf("retrieving Resource Provider %q: `namespace` was nil", id.ResourceProvider)
	}

	if *provider.Namespace != id.ResourceProvider {
		return nil, fmt.Errorf("importing Resource Provider %q: expected %q", id.ResourceProvider, provider.Namespace)
	}

	if provider.RegistrationState == nil || !strings.EqualFold(*provider.RegistrationState, "Registered") {
		return nil, fmt.Errorf("importing Resource Provider %q: Resource Provider must be registered to be imported", id.ResourceProvider)
	}

	account := meta.(*clients.Client).Account
	if err := checkIfResourceProviderManagedByTerraform(id.ResourceProvider, account); err != nil {
		return nil, fmt.Errorf("importing Resource Provider %q: %+v", id.ResourceProvider, err)
	}

	return []*schema.ResourceData{d}, nil
}
