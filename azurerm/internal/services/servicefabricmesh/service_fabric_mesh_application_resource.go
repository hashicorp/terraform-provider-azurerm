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

func resourceServiceFabricMeshApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceFabricMeshApplicationCreateUpdate,
		Read:   resourceServiceFabricMeshApplicationRead,
		Update: resourceServiceFabricMeshApplicationCreateUpdate,
		Delete: resourceServiceFabricMeshApplicationDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ApplicationID(id)
			return err
		}),

		DeprecationMessage: deprecationMessage("azurerm_service_fabric_mesh"),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Follow casing issue here https://github.com/Azure/azure-rest-api-specs/issues/9330
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"location": azure.SchemaLocation(),

			"service": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"image_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"resources": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"requests": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"memory": {
																Type:         schema.TypeFloat,
																Required:     true,
																ValidateFunc: validation.FloatAtLeast(0),
															},
															"cpu": {
																Type:         schema.TypeFloat,
																Required:     true,
																ValidateFunc: validation.FloatAtLeast(0),
															},
														},
													},
												},
												"limits": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"memory": {
																Type:         schema.TypeFloat,
																Required:     true,
																ValidateFunc: validation.FloatAtLeast(0),
															},
															"cpu": {
																Type:         schema.TypeFloat,
																Required:     true,
																ValidateFunc: validation.FloatAtLeast(0),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceServiceFabricMeshApplicationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.ApplicationClient
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
				return fmt.Errorf("checking for presence of existing service fabric mesh application: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_service_fabric_mesh_application", *existing.ID)
		}
	}

	parameters := servicefabricmesh.ApplicationResourceDescription{
		ApplicationResourceProperties: &servicefabricmesh.ApplicationResourceProperties{
			Services: expandServiceFabricMeshApplicationServices(d.Get("service").(*schema.Set).List()),
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.Create(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("creating Service Fabric Mesh Application %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Service Fabric Mesh Application %q (Resource Group %q) to finish creating", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(servicefabricmesh.Creating), string(servicefabricmesh.Upgrading), string(servicefabricmesh.Deleting)},
		Target:                    []string{string(servicefabricmesh.Ready)},
		Refresh:                   serviceFabricMeshApplicationCreateRefreshFunc(ctx, client, resourceGroup, name),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(schema.TimeoutCreate),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Service Fabric Mesh Application %q (Resource Group %q) to become available: %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Service Fabric Mesh Application %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceServiceFabricMeshApplicationRead(d, meta)
}

func resourceServiceFabricMeshApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.ApplicationClient
	serviceClient := meta.(*clients.Client).ServiceFabricMesh.ServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
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

	serviceResp, err := serviceClient.List(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Unable to list Service Fabric Mesh Application Services %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Service Fabric Mesh Application Services: %+v", err)
	}

	if err := d.Set("service", flattenServiceFabricMeshApplicationServices(serviceResp.Values())); err != nil {
		return fmt.Errorf("setting `service`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceServiceFabricMeshApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceFabricMesh.ApplicationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationID(d.Id())
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
			Name:      utils.String(config["name"].(string)),
			Image:     utils.String(config["image_name"].(string)),
			Resources: expandServiceFabricMeshCodePackageResources(config["resources"].([]interface{})),
		})
	}

	return &codePackages
}

func expandServiceFabricMeshCodePackageResources(input []interface{}) *servicefabricmesh.ResourceRequirements {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	attr := input[0].(map[string]interface{})

	return &servicefabricmesh.ResourceRequirements{
		Limits:   expandServiceFabricMeshCodePackageResourceLimits(attr["limits"].([]interface{})),
		Requests: expandServiceFabricMeshCodePackageResourceRequests(attr["requests"].([]interface{})),
	}
}

func expandServiceFabricMeshCodePackageResourceLimits(input []interface{}) *servicefabricmesh.ResourceLimits {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	attr := input[0].(map[string]interface{})

	return &servicefabricmesh.ResourceLimits{
		MemoryInGB: utils.Float(attr["memory"].(float64)),
		CPU:        utils.Float(attr["cpu"].(float64)),
	}
}

func expandServiceFabricMeshCodePackageResourceRequests(input []interface{}) *servicefabricmesh.ResourceRequests {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	attr := input[0].(map[string]interface{})

	return &servicefabricmesh.ResourceRequests{
		MemoryInGB: utils.Float(attr["memory"].(float64)),
		CPU:        utils.Float(attr["cpu"].(float64)),
	}
}

func flattenServiceFabricMeshApplicationServices(input []servicefabricmesh.ServiceResourceDescription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, service := range input {
		attr := make(map[string]interface{})
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
		attr := make(map[string]interface{})
		if codePackage.Name != nil {
			attr["name"] = *codePackage.Name
		}
		if codePackage.Image != nil {
			attr["image_name"] = *codePackage.Image
		}
		attr["resources"] = flattenServiceFabricMeshApplicationCodePackageResources(codePackage.Resources)

		result = append(result, attr)
	}

	return result
}

func flattenServiceFabricMeshApplicationCodePackageResources(input *servicefabricmesh.ResourceRequirements) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if input == nil {
		return result
	}
	attr := make(map[string]interface{})
	attr["requests"] = flattenServiceFabricMeshApplicationCodePackageResourceRequests(input.Requests)
	attr["limits"] = flattenServiceFabricMeshApplicationCodePackageResourceLimits(input.Limits)

	result = append(result, attr)

	return result
}

func flattenServiceFabricMeshApplicationCodePackageResourceRequests(input *servicefabricmesh.ResourceRequests) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if input == nil {
		return result
	}
	attr := make(map[string]interface{})
	if input.MemoryInGB != nil {
		attr["memory"] = *input.MemoryInGB
	}
	if input.CPU != nil {
		attr["cpu"] = *input.CPU
	}
	result = append(result, attr)

	return result
}

func flattenServiceFabricMeshApplicationCodePackageResourceLimits(input *servicefabricmesh.ResourceLimits) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if input == nil {
		return result
	}
	attr := make(map[string]interface{})
	if input.MemoryInGB != nil {
		attr["memory"] = *input.MemoryInGB
	}
	if input.CPU != nil {
		attr["cpu"] = *input.CPU
	}
	result = append(result, attr)

	return result
}

func serviceFabricMeshApplicationCreateRefreshFunc(ctx context.Context, client *servicefabricmesh.ApplicationClient, resourceGroup string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil, "Error", fmt.Errorf("issuing read request in serviceFabricMeshApplicationCreateRefreshFunc %q (Resource Group %q): %s", name, resourceGroup, err)
		}

		return res, string(res.Status), nil
	}
}
