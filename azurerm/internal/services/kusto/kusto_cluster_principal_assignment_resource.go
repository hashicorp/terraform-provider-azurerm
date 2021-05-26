package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2020-09-18/kusto"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/kusto/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKustoClusterPrincipalAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterPrincipalAssignmentCreateUpdate,
		Read:   resourceKustoClusterPrincipalAssignmentRead,
		Delete: resourceKustoClusterPrincipalAssignmentDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterPrincipalAssignmentName,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tenant_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"principal_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"principal_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"principal_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.PrincipalTypeApp),
					string(kusto.PrincipalTypeGroup),
					string(kusto.PrincipalTypeUser),
				}, false),
			},

			"role": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(kusto.AllDatabasesAdmin),
					string(kusto.AllDatabasesViewer),
				}, false),
			},
		},
	}
}

func resourceKustoClusterPrincipalAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Kusto Cluster Principal Assignment creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		principalAssignment, err := client.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(principalAssignment.Response) {
				return fmt.Errorf("Error checking for presence of existing Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
			}
		}

		if principalAssignment.ID != nil && *principalAssignment.ID != "" {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_principal_assignment", *principalAssignment.ID)
		}
	}

	tenantID := d.Get("tenant_id").(string)
	principalID := d.Get("principal_id").(string)
	principalType := d.Get("principal_type").(string)
	role := d.Get("role").(string)

	principalAssignment := kusto.ClusterPrincipalAssignment{
		ClusterPrincipalProperties: &kusto.ClusterPrincipalProperties{
			TenantID:      utils.String(tenantID),
			PrincipalID:   utils.String(principalID),
			PrincipalType: kusto.PrincipalType(principalType),
			Role:          kusto.ClusterPrincipalRole(role),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, clusterName, name, principalAssignment)
	if err != nil {
		return fmt.Errorf("Error creating or updating Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, clusterName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", name, resourceGroup, clusterName, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q)", name, resourceGroup, clusterName)
	}

	d.SetId(*resp.ID)

	return resourceKustoClusterPrincipalAssignmentRead(d, meta)
}

func resourceKustoClusterPrincipalAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterPrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.PrincipalAssignmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", id.PrincipalAssignmentName, id.ResourceGroup, id.ClusterName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	d.Set("name", id.PrincipalAssignmentName)

	tenantID := ""
	if resp.TenantID != nil {
		tenantID = *resp.TenantID
	}

	tenantName := ""
	if resp.TenantName != nil {
		tenantName = *resp.TenantName
	}

	principalID := ""
	if resp.PrincipalID != nil {
		principalID = *resp.PrincipalID
	}

	principalName := ""
	if resp.PrincipalName != nil {
		principalName = *resp.PrincipalName
	}

	principalType := string(resp.PrincipalType)
	role := string(resp.Role)

	d.Set("tenant_id", tenantID)
	d.Set("tenant_name", tenantName)
	d.Set("principal_id", principalID)
	d.Set("principal_name", principalName)
	d.Set("principal_type", principalType)
	d.Set("role", role)

	return nil
}

func resourceKustoClusterPrincipalAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterPrincipalAssignmentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.PrincipalAssignmentName)
	if err != nil {
		return fmt.Errorf("Error deleting Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", id.PrincipalAssignmentName, id.ResourceGroup, id.ClusterName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Kusto Cluster Principal Assignment %q (Resource Group %q, Cluster %q): %+v", id.PrincipalAssignmentName, id.ResourceGroup, id.ClusterName, err)
	}

	return nil
}
