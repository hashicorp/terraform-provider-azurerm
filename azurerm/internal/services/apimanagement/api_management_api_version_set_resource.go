package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiVersionSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiVersionSetCreateUpdate,
		Read:   resourceApiManagementApiVersionSetRead,
		Update: resourceApiManagementApiVersionSetCreateUpdate,
		Delete: resourceApiManagementApiVersionSetDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ApiVersionSetV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"versioning_scheme": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.VersioningSchemeHeader),
					string(apimanagement.VersioningSchemeQuery),
					string(apimanagement.VersioningSchemeSegment),
				}, false),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"version_header_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"version_query_name"},
			},

			"version_query_name": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"version_header_name"},
			},
		},
	}
}

func resourceApiManagementApiVersionSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Api Version Set %q (Api Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_version_set", *existing.ID)
		}
	}

	versioningScheme := apimanagement.VersioningScheme(d.Get("versioning_scheme").(string))
	parameters := apimanagement.APIVersionSetContract{
		APIVersionSetContractProperties: &apimanagement.APIVersionSetContractProperties{
			DisplayName:      utils.String(d.Get("display_name").(string)),
			VersioningScheme: versioningScheme,
			Description:      utils.String(d.Get("description").(string)),
		},
	}

	var headerSet, querySet bool
	if v, ok := d.GetOk("version_header_name"); ok {
		headerSet = v.(string) != ""
		parameters.APIVersionSetContractProperties.VersionHeaderName = utils.String(v.(string))
	}
	if v, ok := d.GetOk("version_query_name"); ok {
		querySet = v.(string) != ""
		parameters.APIVersionSetContractProperties.VersionQueryName = utils.String(v.(string))
	}

	switch schema := versioningScheme; schema {
	case apimanagement.VersioningSchemeHeader:
		if !headerSet {
			return fmt.Errorf("`version_header_name` must be set if `versioning_schema` is `Header`")
		}
		if querySet {
			return fmt.Errorf("`version_query_name` can not be set if `versioning_schema` is `Header`")
		}

	case apimanagement.VersioningSchemeQuery:
		if headerSet {
			return fmt.Errorf("`version_header_name` can not be set if `versioning_schema` is `Query`")
		}
		if !querySet {
			return fmt.Errorf("`version_query_name` must be set if `versioning_schema` is `Query`")
		}

	case apimanagement.VersioningSchemeSegment:
		if headerSet {
			return fmt.Errorf("`version_header_name` can not be set if `versioning_schema` is `Segment`")
		}
		if querySet {
			return fmt.Errorf("`version_query_name` can not be set if `versioning_schema` is `Segment`")
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("creating/updating Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("retrieving Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Api Version Set %q (Resource Group %q / Api Management Service %q)", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementApiVersionSetRead(d, meta)
}

func resourceApiManagementApiVersionSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiVersionSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Api Version Set %q (Resource Group %q / Api Management Service %q) was not found - removing from state!", id.Name, id.ResourceGroup, id.ServiceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", id.Name, id.ResourceGroup, id.ServiceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	if props := resp.APIVersionSetContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("versioning_scheme", string(props.VersioningScheme))
		d.Set("version_header_name", props.VersionHeaderName)
		d.Set("version_query_name", props.VersionQueryName)
	}

	return nil
}

func resourceApiManagementApiVersionSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiVersionSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiVersionSetID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", id.Name, id.ResourceGroup, id.ServiceName, err)
		}
	}

	return nil
}
