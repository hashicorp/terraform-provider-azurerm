package servicefabricmesh

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/servicefabricmesh/mgmt/2018-09-01-preview/servicefabricmesh"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabricmesh/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceFabricMeshNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceFabricMeshNetworkCreateUpdate,
		Read:   resourceArmServiceFabricMeshNetworkRead,
		Update: resourceArmServiceFabricMeshNetworkCreateUpdate,
		Delete: resourceArmServiceFabricMeshNetworkDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServiceFabricMeshNetworkID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Follow casing issue here https://github.com/Azure/azure-rest-api-specs/issues/9330
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocation(),

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmServiceFabricMeshNetworkCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.NetworkClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Service Fabric Mesh Network: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_service_fabric_mesh_network", *existing.ID)
		}
	}

	parameters := servicefabricmesh.NetworkResourceDescription{
		Properties: &servicefabricmesh.NetworkResourceProperties{
			Description: utils.String(d.Get("description").(string)),
			Kind:        servicefabricmesh.KindNetworkResourcePropertiesBase,
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.Create(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("creating Service Fabric Mesh Network %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Service Fabric Mesh Network %q (Resource Group %q) to finish creating", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(servicefabricmesh.Creating), string(servicefabricmesh.Upgrading), string(servicefabricmesh.Deleting)},
		Target:                    []string{string(servicefabricmesh.Ready)},
		Refresh:                   serviceFabricMeshNetworkCreateRefreshFunc(ctx, client, resourceGroup, name),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Service Fabric Mesh Network %q (Resource Group %q) to become available: %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Service Fabric Mesh Network %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmServiceFabricMeshNetworkRead(d, meta)
}

func resourceArmServiceFabricMeshNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.NetworkClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceFabricMeshNetworkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Unable to find Service Fabric Mesh Network %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Service Fabric Mesh Network: %+v", err)
	}

	props, ok := resp.Properties.AsNetworkResourceProperties()
	if !ok {
		return fmt.Errorf("classifiying Service Fabric Mesh Network %q (Resource Group %q): Expected: %q Received: %q", id.Name, id.ResourceGroup, servicefabricmesh.KindNetworkResourceProperties, props.Kind)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("description", props.Description)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmServiceFabricMeshNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.NetworkClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceFabricMeshNetworkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Service Fabric Mesh Network %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func serviceFabricMeshNetworkCreateRefreshFunc(ctx context.Context, client *servicefabricmesh.NetworkClient, resourceGroup string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil, "Error", fmt.Errorf("issuing read request in serviceFabricMeshNetworkCreateRefreshFunc %q (Resource Group %q): %s", name, resourceGroup, err)
		}

		props, ok := res.Properties.AsNetworkResourceProperties()
		if !ok {
			return res, string(props.Status), fmt.Errorf("classifiying Service Fabric Mesh Network %q (Resource Group %q): Expected: %q Received: %q", name, resourceGroup, servicefabricmesh.KindNetworkResourceProperties, props.Kind)
		}

		return res, string(props.Status), nil
	}
}
