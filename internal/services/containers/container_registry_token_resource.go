package containers

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/date"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2020-11-01-preview/containerregistry"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceContainerRegistryToken() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceContainerRegistryTokenCreate,
		Read:     resourceContainerRegistryTokenRead,
		Update:   resourceContainerRegistryTokenUpdate,
		Delete:   resourceContainerRegistryTokenDelete,
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryTokenName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},

			"scope_map_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryScopeMapID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"password": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 2,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"password1",
								"password2",
							}, false),
						},
						"expiry": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"value": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceContainerRegistryTokenCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry token creation.")
	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing token %q in Container Registry %q (Resource Group %q): %s", name, containerRegistryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_registry_token", *existing.ID)
		}
	}

	scopeMapID := d.Get("scope_map_id").(string)
	enabled := d.Get("enabled").(bool)
	status := containerregistry.TokenStatusEnabled

	if !enabled {
		status = containerregistry.TokenStatusDisabled
	}

	parameters := containerregistry.Token{
		TokenProperties: &containerregistry.TokenProperties{
			ScopeMapID: utils.String(scopeMapID),
			Status:     status,
		},
	}

	future, err := client.Create(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("creating token %q in Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of token %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		return fmt.Errorf("retrieving token %q for Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read token %q for Container Registry %q (resource group %q) ID", name, containerRegistryName, resourceGroup)
	}

	d.SetId(*read.ID)

	passwords, err := expandContainerRegistryTokenPassword(d.Get("password").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `password`: %v", err)
	}
	if passwords != nil {
		oldPasswords := *passwords
		for idx, password := range oldPasswords {
			param := containerregistry.GenerateCredentialsParameters{
				TokenID: read.ID,
				Expiry:  password.Expiry,
				Name:    password.Name,
			}
			regClient := meta.(*clients.Client).Containers.RegistriesClient
			future, err := regClient.GenerateCredentials(ctx, resourceGroup, containerRegistryName, param)
			if err != nil {
				return fmt.Errorf("generating password credential %s: %v", password.Name, err)
			}
			if err := future.WaitForCompletionRef(ctx, regClient.Client); err != nil {
				return fmt.Errorf("waiting for password credential generation for %s: %v", password.Name, err)
			}

			result, err := future.Result(*regClient)
			if err != nil {
				return fmt.Errorf("getting password credential after creation for %s: %v", password.Name, err)
			}

			if result.Passwords != nil && len(*result.Passwords) > idx {
				password.Value = (*result.Passwords)[idx].Value
				(*passwords)[idx] = password
			}
		}
		// The password is only known right after it is generated, therefore setting it to the resource data here.
		if err := d.Set("password", flattenContainerRegistryTokenPassword(passwords)); err != nil {
			return fmt.Errorf(`setting "passwords": %v`, err)
		}
	}

	return resourceContainerRegistryTokenRead(d, meta)
}

func resourceContainerRegistryTokenUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Container Registry token update.")
	resourceGroup := d.Get("resource_group_name").(string)
	containerRegistryName := d.Get("container_registry_name").(string)
	name := d.Get("name").(string)
	scopeMapID := d.Get("scope_map_id").(string)
	enabled := d.Get("enabled").(bool)
	status := containerregistry.TokenStatusEnabled

	if !enabled {
		status = containerregistry.TokenStatusDisabled
	}

	parameters := containerregistry.TokenUpdateParameters{
		TokenUpdateProperties: &containerregistry.TokenUpdateProperties{
			ScopeMapID: utils.String(scopeMapID),
			Status:     status,
		},
	}

	future, err := client.Update(ctx, resourceGroup, containerRegistryName, name, parameters)
	if err != nil {
		return fmt.Errorf("updating token %q for Container Registry %q (Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of token %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, containerRegistryName, name)
	if err != nil {
		return fmt.Errorf("retrieving token %q (Container Registry %q, Resource Group %q): %+v", name, containerRegistryName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read token %q (Container Registry %q, resource group %q) ID", name, containerRegistryName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceContainerRegistryTokenRead(d, meta)
}

func resourceContainerRegistryTokenRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerRegistryTokenID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Token %q was not found in Container Registry %q in Resource Group %q", id.TokenName, id.RegistryName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on token %q in Azure Container Registry %q (Resource Group %q): %+v", id.TokenName, id.RegistryName, id.ResourceGroup, err)
	}

	status := true
	if resp.Status == containerregistry.TokenStatusDisabled {
		status = false
	}

	d.Set("name", resp.Name)
	d.Set("container_registry_name", id.RegistryName)
	d.Set("enabled", status)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("scope_map_id", resp.ScopeMapID)

	return nil
}

func resourceContainerRegistryTokenDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.TokensClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ContainerRegistryTokenID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandContainerRegistryTokenPassword(input []interface{}) (*[]containerregistry.TokenPassword, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]containerregistry.TokenPassword, 0)

	for _, e := range input {
		e := e.(map[string]interface{})

		var dt date.Time
		if v := e["expiry"].(string); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
			dt.Time = t
		}
		result = append(result, containerregistry.TokenPassword{
			Expiry: &dt,
			Name:   containerregistry.TokenPasswordName(e["name"].(string)),
		})
	}
	return &result, nil
}

func flattenContainerRegistryTokenPassword(input *[]containerregistry.TokenPassword) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var expiry string
		if e.Expiry != nil {
			expiry = e.Expiry.String()
		}

		var value string
		if e.Value != nil {
			value = *e.Value
		}

		output = append(output, map[string]interface{}{
			"name":   string(e.Name),
			"expiry": expiry,
			"value":  value,
		})
	}
	return output
}
