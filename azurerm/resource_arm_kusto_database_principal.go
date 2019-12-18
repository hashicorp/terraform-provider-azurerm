package azurerm

import (
	"fmt"
	"log"
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
				Computed:     true,
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
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.DatabasePrincipalTypeApp),
					string(kusto.DatabasePrincipalTypeGroup),
					string(kusto.DatabasePrincipalTypeUser),
				}, false),
			},

			// TODO try all of these together to see if I need to do ExactlyOneOf
			"fully_qualified_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"email": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
				ConflictsWith: []string{"app_id"},
			},

			"app_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
				ConflictsWith: []string{"email"},
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
	role := d.Get("role").(string)

	database, err := client.Get(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if utils.ResponseWasNotFound(database.Response) {
			return fmt.Errorf("Kusto Database %q (Resource Group %q) was not found", databaseName, resourceGroup)
		}

		return fmt.Errorf("Error loading Kusto Database %q (Resource Group %q): %+v", databaseName, resourceGroup, err)
	}
	resourceId := fmt.Sprintf("%s/Principals/%s/Role/%s/Type/%s", *database.ID, name, role)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Database Principals %q (Resource Group %q, Cluster %q): %s", name, resourceGroup, clusterName, err)
			}
		}

		if principals := resp.Value; principals != nil {
			for _, principal := range *principals {
				// kusto database principals are unique when looked at with name and role
				if principal.Name != nil && *principal.Name == name && string(principal.Role) == role {
					return tf.ImportAsExistsError("azurerm_kusto_database_principal", resourceId)
				}
			}
		}
	}

	kustoPrincipal := kusto.DatabasePrincipal{
		Type:  kusto.DatabasePrincipalType(d.Get("type").(string)),
		Role:  kusto.DatabasePrincipalRole(role),
		Name:  utils.String(""),
		Fqn:   utils.String(d.Get("fully_qualified_name").(string)),
		Email: utils.String(d.Get("email").(string)),
		AppID: utils.String(d.Get("app_id").(string)),
	}

	principals := []kusto.DatabasePrincipal{kustoPrincipal}
	request := kusto.DatabasePrincipalListRequest{
		Value: &principals,
	}

	_, err = client.AddPrincipals(ctx, resourceGroup, clusterName, databaseName, request)
	if err != nil {
		return fmt.Errorf("Error creating Kusto Database Principal %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	d.SetId(resourceId)

	return resourceArmKustoDatabasePrincipalRead(d, meta)
}

func resourceArmKustoDatabasePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	databaseName := id.Path["Databases"]
	name := id.Path["Principals"]
	role := id.Path["Role"]

	databaseResponse, err := client.Get(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if utils.ResponseWasNotFound(databaseResponse.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for presence of existing Kusto Database Principals %q (Resource Group %q, Cluster %q): %s", name, resourceGroup, clusterName, err)
		}
	}

	principal := kusto.DatabasePrincipal{}
	found := false
	if principals := resp.Value; principals != nil {
		for _, currPrincipal := range *principals {
			// kusto database principals are unique when looked at with name and role
			if currPrincipal.Name != nil && *currPrincipal.Name == name && string(currPrincipal.Role) == role {
				principal = currPrincipal
				found = true
				break
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] Kusto Database Principal %q was not found - removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", principal.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("cluster_name", clusterName)
	d.Set("database_name", databaseName)

	d.Set("role", string(principal.Role))
	d.Set("type", string(principal.Type))
	if fqn := principal.Fqn; fqn != nil {
		d.Set("fully_qualified_name", fqn)
	}
	if email := principal.Email; email != nil {
		d.Set("email", email)
	}
	if appID := principal.AppID; appID != nil {
		d.Set("app_id", appID)
	}

	return nil
}

func resourceArmKustoDatabasePrincipalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	databaseName := id.Path["Databases"]
	name := id.Path["Principals"]
	role := id.Path["Role"]

	kustoPrincipal := kusto.DatabasePrincipal{
		Role:  kusto.DatabasePrincipalRole(role),
		Name:  utils.String(""),
		Fqn:   utils.String(d.Get("fully_qualified_name").(string)),
		Email: utils.String(d.Get("email").(string)),
		AppID: utils.String(d.Get("app_id").(string)),
		Type:  kusto.DatabasePrincipalType(d.Get("type").(string)),
	}

	principals := []kusto.DatabasePrincipal{kustoPrincipal}
	request := kusto.DatabasePrincipalListRequest{
		Value: &principals,
	}

	_, err = client.RemovePrincipals(ctx, resGroup, clusterName, databaseName, request)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Database Principal %q (Resource Group %q, Cluster %q, Database %q): %+v", name, resGroup, clusterName, databaseName, err)
	}

	return nil
}
