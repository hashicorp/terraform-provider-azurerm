package kusto

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-02-15/kusto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKustoDatabasePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKustoDatabasePrincipalCreate,
		Read:   resourceArmKustoDatabasePrincipalRead,
		Delete: resourceArmKustoDatabasePrincipalDelete,

		DeprecationMessage: "This resource has been superseded by `azurerm_kusto_database_principal_assignment` to reflects changes in the API/SDK and will be removed in version 3.0 of the provider.",

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
	switch principalType {
	case "User":
		fqn = fmt.Sprintf("aaduser=%s;%s", objectID, clientID)
	case "Group":
		fqn = fmt.Sprintf("aadgroup=%s;%s", objectID, clientID)
	case "App":
		fqn = fmt.Sprintf("aadapp=%s;%s", objectID, clientID)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Kusto Database %q (Resource Group %q, Cluster %q) was not found", databaseName, resourceGroup, clusterName)
		}

		return fmt.Errorf("Error loading Kusto Database %q (Resource Group %q, Cluster %q): %+v", databaseName, resourceGroup, clusterName, err)
	}
	if resp.Value == nil {
		return fmt.Errorf("Error loading Kusto Database %q (Resource Group %q, Cluster %q): Invalid resource response", databaseName, resourceGroup, clusterName)
	}

	database, ok := resp.Value.AsReadWriteDatabase()
	if !ok {
		return fmt.Errorf("Exisiting resource is not a Kusto Read/Write Database %q (Resource Group %q, Cluster %q)", databaseName, resourceGroup, clusterName)
	}

	resourceID := fmt.Sprintf("%s/Role/%s/FQN/%s", *database.ID, role, fqn)

	if d.IsNewResource() {
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
					return tf.ImportAsExistsError("azurerm_kusto_database_principal", resourceID)
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
		Name: utils.String(""),
	}

	principals := []kusto.DatabasePrincipal{kustoPrincipal}
	request := kusto.DatabasePrincipalListRequest{
		Value: &principals,
	}

	if _, err = client.AddPrincipals(ctx, resourceGroup, clusterName, databaseName, request); err != nil {
		return fmt.Errorf("Error creating Kusto Database Principal (Resource Group %q, Cluster %q): %+v", resourceGroup, clusterName, err)
	}

	principalsResp, err := client.ListPrincipals(ctx, resourceGroup, clusterName, databaseName)
	if err != nil {
		if !utils.ResponseWasNotFound(principalsResp.Response) {
			return fmt.Errorf("Error checking for presence of existing Kusto Database Principals (Resource Group %q, Cluster %q): %s", resourceGroup, clusterName, err)
		}
	}

	d.SetId(resourceID)

	return resourceArmKustoDatabasePrincipalRead(d, meta)
}

func resourceArmKustoDatabasePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabasePrincipalID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Database %q (Resource Group %q, Cluster %q): %+v", id.DatabaseName, id.ResourceGroup, id.ClusterName, err)
	}

	databasePrincipals, err := client.ListPrincipals(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName)
	if err != nil {
		if !utils.ResponseWasNotFound(databasePrincipals.Response) {
			return fmt.Errorf("Error checking for presence of existing Kusto Database Principals %q (Resource Group %q, Cluster %q): %s", id, id.ResourceGroup, id.ClusterName, err)
		}
	}

	principal := kusto.DatabasePrincipal{}
	found := false
	if principals := databasePrincipals.Value; principals != nil {
		for _, currPrincipal := range *principals {
			// kusto database principals are unique when looked at with fqn and role
			if string(currPrincipal.Role) == id.RoleName && currPrincipal.Fqn != nil && *currPrincipal.Fqn == id.FQNName {
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

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	d.Set("database_name", id.DatabaseName)

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

	splitFQN := strings.Split(id.FQNName, "=")
	if len(splitFQN) != 2 {
		return fmt.Errorf("Expected `fqn` to be in the format aadtype=objectid:clientid but got: %q", id.FQNName)
	}
	splitIDs := strings.Split(splitFQN[1], ";")
	if len(splitIDs) != 2 {
		return fmt.Errorf("Expected `fqn` to be in the format aadtype=objectid:clientid but got: %q", id.FQNName)
	}
	d.Set("object_id", splitIDs[0])
	d.Set("client_id", splitIDs[1])

	return nil
}

func resourceArmKustoDatabasePrincipalDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabasePrincipalID(d.Id())
	if err != nil {
		return err
	}

	kustoPrincipal := kusto.DatabasePrincipal{
		Role: kusto.DatabasePrincipalRole(id.RoleName),
		Fqn:  utils.String(id.FQNName),
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

	if _, err = client.RemovePrincipals(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, request); err != nil {
		return fmt.Errorf("Error deleting Kusto Database Principal %q (Resource Group %q, Cluster %q, Database %q): %+v", id, id.ResourceGroup, id.ClusterName, id.DatabaseName, err)
	}

	return nil
}
