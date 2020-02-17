package kusto

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2019-05-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
			// TODO: confirm these
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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

			"client_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"fully_qualified_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// These must be computed as the values passed in are overwritten by what the `fqn` returns.
			// For more info: https://github.com/Azure/azure-sdk-for-go/issues/6547
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmKustoDatabasePrincipalCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Database Principal creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)
	role := d.Get("role").(string)
	principalType := d.Get("type").(string)

	clientID := d.Get("client_id").(string)
	objectID := d.Get("object_id").(string)
	fqn := ""
	if principalType == "User" {
		fqn = fmt.Sprintf("aaduser=%s;%s", objectID, clientID)
	} else if principalType == "Group" {
		fqn = fmt.Sprintf("aadgroup=%s;%s", objectID, clientID)
	} else if principalType == "App" {
		fqn = fmt.Sprintf("aadapp=%s;%s", objectID, clientID)
	}

	database, err := client.Get(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if utils.ResponseWasNotFound(database.Response) {
			return fmt.Errorf("Kusto Database %q (Resource Group %q) was not found", databaseName, resourceGroup)
		}

		return fmt.Errorf("Error loading Kusto Database %q (Resource Group %q): %+v", databaseName, resourceGroup, err)
	}
	resourceId := fmt.Sprintf("%s/Role/%s/FQN/%s", *database.ID, role, fqn)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Database Principals (Resource Group %q, Cluster %q): %s", resourceGroup, clusterName, err)
			}
		}

		if principals := resp.Value; principals != nil {
			for _, principal := range *principals {
				// kusto database principals are unique when looked at with name and role
				if string(principal.Role) == role && principal.Fqn != nil && *principal.Fqn == fqn {
					return tf.ImportAsExistsError("azurerm_kusto_database_principal", resourceId)
				}
			}
		}
	}

	kustoPrincipal := kusto.DatabasePrincipal{
		Type: kusto.DatabasePrincipalType(principalType),
		Role: kusto.DatabasePrincipalRole(role),
		Fqn:  utils.String(fqn),
		// These three must be specified or the api returns `The request is invalid.`
		// For more info: https://github.com/Azure/azure-sdk-for-go/issues/6547
		Email: utils.String(""),
		AppID: utils.String(""),
		Name:  utils.String(""),
	}

	principals := []kusto.DatabasePrincipal{kustoPrincipal}
	request := kusto.DatabasePrincipalListRequest{
		Value: &principals,
	}

	if _, err = client.AddPrincipals(ctx, resourceGroup, clusterName, databaseName, request); err != nil {
		return fmt.Errorf("Error creating Kusto Database Principal (Resource Group %q, Cluster %q): %+v", resourceGroup, clusterName, err)
	}

	resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for presence of existing Kusto Database Principals (Resource Group %q, Cluster %q): %s", resourceGroup, clusterName, err)
		}
	}

	d.SetId(resourceId)

	return resourceArmKustoDatabasePrincipalRead(d, meta)
}

func resourceArmKustoDatabasePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	databaseName := id.Path["Databases"]
	role := id.Path["Role"]
	fqn := id.Path["FQN"]

	databaseResponse, err := client.Get(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if utils.ResponseWasNotFound(databaseResponse.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", databaseName, resourceGroup, clusterName, err)
	}

	resp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for presence of existing Kusto Database Principals %q (Resource Group %q, Cluster %q): %s", id, resourceGroup, clusterName, err)
		}
	}

	principal := kusto.DatabasePrincipal{}
	found := false
	if principals := resp.Value; principals != nil {
		for _, currPrincipal := range *principals {
			// kusto database principals are unique when looked at with fqn and role
			if string(currPrincipal.Role) == role && currPrincipal.Fqn != nil && *currPrincipal.Fqn == fqn {
				principal = currPrincipal
				found = true
				break
			}
		}
	}

	if !found {
		log.Printf("[DEBUG] Kusto Database Principal %q was not found - removing from state", id)
		d.SetId("")
		return nil
	}

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
	if name := principal.Name; name != nil {
		d.Set("name", principal.Name)
	}

	splitFQN := strings.Split(fqn, "=")
	if len(splitFQN) != 2 {
		return fmt.Errorf("Expected `fqn` to be in the format aadtype=objectid:clientid but got: %q", fqn)
	}
	splitIDs := strings.Split(splitFQN[1], ";")
	if len(splitIDs) != 2 {
		return fmt.Errorf("Expected `fqn` to be in the format aadtype=objectid:clientid but got: %q", fqn)
	}
	d.Set("object_id", splitIDs[0])
	d.Set("client_id", splitIDs[1])

	return nil
}

func resourceArmKustoDatabasePrincipalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	clusterName := id.Path["Clusters"]
	databaseName := id.Path["Databases"]
	role := id.Path["Role"]
	fqn := id.Path["FQN"]

	kustoPrincipal := kusto.DatabasePrincipal{
		Role: kusto.DatabasePrincipalRole(role),
		Fqn:  utils.String(fqn),
		Type: kusto.DatabasePrincipalType(d.Get("type").(string)),
		// These three must be specified or the api returns `The request is invalid.`
		// For more info: https://github.com/Azure/azure-sdk-for-go/issues/6547
		Name:  utils.String(""),
		Email: utils.String(""),
		AppID: utils.String(""),
	}

	principals := []kusto.DatabasePrincipal{kustoPrincipal}
	request := kusto.DatabasePrincipalListRequest{
		Value: &principals,
	}

	if _, err = client.RemovePrincipals(ctx, resGroup, clusterName, databaseName, request); err != nil {
		return fmt.Errorf("Error deleting Kusto Database Principal %q (Resource Group %q, Cluster %q, Database %q): %+v", id, resGroup, clusterName, databaseName, err)
	}

	return nil
}
