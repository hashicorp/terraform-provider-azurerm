package servicefabricmesh

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/servicefabricmesh/mgmt/2018-09-01-preview/servicefabricmesh"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicefabricmesh/parse"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceFabricMeshApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceFabricMeshApplicationCreateUpdate,
		Read:   resourceArmServiceFabricMeshApplicationRead,
		Update: resourceArmServiceFabricMeshApplicationCreateUpdate,
		Delete: resourceArmServiceFabricMeshApplicationDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServiceFabricMeshApplicationID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"service": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"os_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(servicefabricmesh.Linux),
								string(servicefabricmesh.Windows),
							}, false),
						},
						"code_package": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmServiceFabricMeshApplicationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.ApplicationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing service fabric mesh application: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_service_fabric_mesh_application", *existing.ID)
		}
	}

	parameters := servicefabricmesh.ApplicationResourceDescription{
		ApplicationResourceProperties: &servicefabricmesh.ApplicationResourceProperties{
			Description: utils.String(d.Get("description").(string)),
			Services:    expandServiceFabricMeshApplicationServices(d.Get("service").(*schema.Set).List()),
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.Create(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("creating Service Fabric Mesh Application %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Service Fabric Mesh Application %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmServiceFabricMeshApplicationRead(d, meta)
}

func resourceArmServiceFabricMeshApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.ApplicationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceFabricMeshApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Unable to find Service Fabric Mesh Application %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Service Fabric Mesh Application: %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("service", flattenServiceFabricMeshApplicationServices(resp.Services)); err != nil {
		return fmt.Errorf("setting `service`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmServiceFabricMeshApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceFabricMeshApplicationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Service Fabric Mesh Application %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandServiceFabricMeshApplicationServices(input []interface{}) *[]servicefabricmesh.ServiceResourceDescription {
	services := make([]servicefabricmesh.ServiceResourceDescription, 0)

	for _, serviceConfig := range input {
		if serviceConfig == nil {
			continue
		}
		config := serviceConfig.(map[string]interface{})

		services = append(services, servicefabricmesh.ServiceResourceDescription{
			Name: utils.String(config["name"].(string)),
			ServiceResourceProperties: &servicefabricmesh.ServiceResourceProperties{
				OsType:       servicefabricmesh.OperatingSystemType(config["os_type"].(string)),
				CodePackages: expandServiceFabricMeshCodePackages(config["code_package"].(*schema.Set).List()),
			},
		})
	}

	return &services
}

func expandServiceFabricMeshCodePackages(input []interface{}) *[]servicefabricmesh.ContainerCodePackageProperties {
	codePackages := make([]servicefabricmesh.ContainerCodePackageProperties, 0, len(input))

	for _, codePackageConfig := range input {
		if codePackageConfig == nil {
			continue
		}
		config := codePackageConfig.(map[string]interface{})

		codePackages = append(codePackages, servicefabricmesh.ContainerCodePackageProperties{
			Name: utils.String(config["name"].(string)),
		})
	}

	return &codePackages
}

func flattenServiceFabricMeshApplicationServices(input *[]servicefabricmesh.ServiceResourceDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if input == nil {
		return result
	}

	for _, service := range *input {
		attr := make(map[string]interface{}, 0)
		if service.Name != nil {
			attr["name"] = *service.Name
		}

		attr["os_type"] = string(service.OsType)

		attr["code_package"] = flattenServiceFabricMeshApplicationCodePackage(service.CodePackages)

		result = append(result, attr)
	}

	return result
}

func flattenServiceFabricMeshApplicationCodePackage(input *[]servicefabricmesh.ContainerCodePackageProperties) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if input == nil {
		return result
	}

	for _, codePackage := range *input {
		attr := make(map[string]interface{}, 0)
		if codePackage.Name != nil {
			attr["name"] = *codePackage.Name
		}
	}

	return result
}
