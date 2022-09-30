package kusto

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoClusterPrincipalAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterPrincipalAssignmentCreateUpdate,
		Read:   resourceKustoClusterPrincipalAssignmentRead,
		Delete: resourceKustoClusterPrincipalAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterPrincipalAssignmentID(id)
			return err
		}),

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
					string(kusto.ClusterPrincipalRoleAllDatabasesAdmin),
					string(kusto.ClusterPrincipalRoleAllDatabasesViewer),
				}, false),
			},
		},
	}
}

func resourceKustoClusterPrincipalAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterPrincipalAssignmentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewClusterPrincipalAssignmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		principalAssignment, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.PrincipalAssignmentName)
		if err != nil {
			if !utils.ResponseWasNotFound(principalAssignment.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(principalAssignment.Response) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_principal_assignment", id.ID())
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ClusterName, id.PrincipalAssignmentName, principalAssignment)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
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
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrincipalAssignmentName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroup)

	principalID := ""
	if resp.PrincipalID != nil {
		principalID = *resp.PrincipalID
	}
	d.Set("principal_id", principalID)

	principalName := ""
	if resp.PrincipalName != nil {
		principalName = *resp.PrincipalName
	}
	d.Set("principal_name", principalName)

	d.Set("principal_type", string(resp.PrincipalType))
	d.Set("role", string(resp.Role))

	tenantID := ""
	if resp.TenantID != nil {
		tenantID = *resp.TenantID
	}
	d.Set("tenant_id", tenantID)

	tenantName := ""
	if resp.TenantName != nil {
		tenantName = *resp.TenantName
	}
	d.Set("tenant_name", tenantName)

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
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
