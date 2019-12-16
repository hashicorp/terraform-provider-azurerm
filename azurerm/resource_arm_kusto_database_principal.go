package azurerm

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-05-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoDatabasePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoDatabasePrincipalCreate,
		Read:   resourceArmKustoDatabasePrincipalRead,
		Delete: resourceArmKustoDatabasePrincipalDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoClusterName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMKustoDatabaseName,
			},

			"role": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.Admin),
					string(kusto.Ingestor),
					string(kusto.Monitor),
					string(kusto.User),
					string(kusto.UnrestrictedViewers),
					string(kusto.Viewer),
				}, false),
			},

			"type": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.DatabasePrincipalTypeApp),
					string(kusto.DatabasePrincipalTypeGroup),
					string(kusto.DatabasePrincipalTypeUser),
				}, false),
			},

			"fully_qualified_domain_name": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"email": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"application_id": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmKustoDatabasePrincipalCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Database Principal creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Database Principals %q (Resource Group %q, Cluster %q): %s", name, resourceGroup, clusterName, err)
			}
		}

		if principals := resp.Value; principals != nil {
			for _, principal := range(*principals) {
				if principal.Name != nil && *principal.Name == name {
					return tf.ImportAsExistsError("azurerm_kusto_database_principal", *principal.Name)
				}
			}
		}
	}

	kustoPrincipal := kusto.DatabasePrincipal{
		Type: kusto.DatabasePrincipalType(d.Get("type").(string)),
	}

	principals := []kusto.DatabasePrincipal{kustoPrincipal}
	request := kusto.DatabasePrincipalListRequest{
		Value: &principals,
	}

	resp, err := client.AddPrincipals(ctx, resourceGroup, clusterName, databaseName, request)
	if err != nil {
		return fmt.Errorf("Error creating Kusto Database Principal %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	d.SetId(*resp.)

	return resourceArmKustoDatabaseRead(d, meta)
}

func resourceArmKustoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	name := id.Path["Databases"]

	databaseResponse, err := client.Get(ctx, resourceGroup, clusterName, name)

	if err != nil {
		if utils.ResponseWasNotFound(databaseResponse.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("cluster_name", clusterName)

	if location := databaseResponse.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := databaseResponse.DatabaseProperties; props != nil {
		d.Set("hot_cache_period", props.HotCachePeriod)
		d.Set("soft_delete_period", props.SoftDeletePeriod)

		if statistics := props.Statistics; statistics != nil {
			d.Set("size", statistics.Size)
		}
	}

	return nil
}

func resourceArmKustoDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	name := id.Path["Databases"]

	future, err := client.Delete(ctx, resGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Database Principal %q (Resource Group %q, Cluster %q): %+v", name, resGroup, clusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Database Principal %q (Resource Group %q, Cluster %q): %+v", name, resGroup, clusterName, err)
	}

	return nil
}

func validateAzureRMKustoDatabasePrincipalName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespaces only", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s.-]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, whitespaces, dashes and dots: %q", k, name))
	}

	if len(name) > 260 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 4 and 22 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}

func expandKustoDatabasePrincipalProperties(d *schema.ResourceData) *kusto.DatabaseProperties {
	databaseProperties := &kusto.DatabaseProperties{}

	if softDeletePeriod, ok := d.GetOk("soft_delete_period"); ok {
		databaseProperties.SoftDeletePeriod = utils.String(softDeletePeriod.(string))
	}

	if hotCachePeriod, ok := d.GetOk("hot_cache_period"); ok {
		databaseProperties.HotCachePeriod = utils.String(hotCachePeriod.(string))
	}

	return databaseProperties
}
