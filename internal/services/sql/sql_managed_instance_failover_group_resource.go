package sql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSqlInstanceFailoverGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSqlInstanceFailoverGroupCreateUpdate,
		Read:   resourceSqlInstanceFailoverGroupRead,
		Update: resourceSqlInstanceFailoverGroupCreateUpdate,
		Delete: resourceSqlInstanceFailoverGroupDelete,

		DeprecationMessage: "The `azurerm_sql_managed_instance_failover_group` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_managed_instance_failover_group` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.InstanceFailoverGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlFailoverGroupName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"managed_instance_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlServerName,
			},

			"partner_managed_instance_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"partner_region": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"location": commonschema.LocationComputed(),

						"role": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"readonly_endpoint_failover_policy_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"read_write_endpoint_failover_policy": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.ReadWriteEndpointFailoverPolicyAutomatic),
								string(sql.ReadWriteEndpointFailoverPolicyManual),
							}, false),
						},
						"grace_minutes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(60),
						},
					},
				},
			},

			"role": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSqlInstanceFailoverGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.InstanceFailoverGroupsClient
	instanceClient := meta.(*clients.Client).Sql.ManagedInstancesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewInstanceFailoverGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("location").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.LocationName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_failover_group", id.ID())
		}
	}

	partnerRegions := make([]sql.PartnerRegionInfo, 0)
	partnerId, err := parse.ManagedInstanceID(d.Get("partner_managed_instance_id").(string))
	if err != nil {
		return err
	}
	resp, err := instanceClient.Get(ctx, partnerId.ResourceGroup, partnerId.Name, "")
	if err != nil || resp.Location == nil || *resp.Location == "" {
		return fmt.Errorf("checking for existence and region of Partner of %q: %+v", id, err)
	}

	regionInfo := sql.PartnerRegionInfo{
		Location: utils.String(*resp.Location),
	}
	partnerRegions = append(partnerRegions, regionInfo)

	primaryInstanceId := parse.NewManagedInstanceID(subscriptionId, id.ResourceGroup, d.Get("managed_instance_name").(string))
	properties := sql.InstanceFailoverGroup{
		InstanceFailoverGroupProperties: &sql.InstanceFailoverGroupProperties{
			ReadOnlyEndpoint:     expandSqlInstanceFailoverGroupReadOnlyPolicy(d),
			ReadWriteEndpoint:    expandSqlInstanceFailoverGroupReadWritePolicy(d),
			ManagedInstancePairs: expandSqlInstanceFailoverGroupManagedInstanceId(d, primaryInstanceId),
			PartnerRegions:       &partnerRegions,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LocationName, id.Name, properties)
	if err != nil {
		return fmt.Errorf("issuing create/update request for %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on create/update future for %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSqlInstanceFailoverGroupRead(d, meta)
}

func resourceSqlInstanceFailoverGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.InstanceFailoverGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.InstanceFailoverGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.LocationName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Failover Group %q (Location %q / Resource Group %q): %+v", id.Name, id.LocationName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("location", azure.NormalizeLocation(id.LocationName))
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.InstanceFailoverGroupProperties; props != nil {
		name, err := flattenSqlInstanceFailoverGroupPrimaryInstance(props.ManagedInstancePairs, props.ReplicationRole)
		if err != nil {
			return fmt.Errorf("flatten `managed_instance_name`: %+v", err)
		}

		if err := d.Set("managed_instance_name", name); err != nil {
			return fmt.Errorf("setting `managed_instance_name`: %+v", err)
		}

		if err := d.Set("read_write_endpoint_failover_policy", flattenSqlInstanceFailoverGroupReadWritePolicy(props.ReadWriteEndpoint)); err != nil {
			return fmt.Errorf("setting `read_write_endpoint_failover_policy`: %+v", err)
		}

		if err := d.Set("readonly_endpoint_failover_policy_enabled", props.ReadOnlyEndpoint.FailoverPolicy == sql.ReadOnlyEndpointFailoverPolicyEnabled); err != nil {
			return fmt.Errorf("setting `readonly_endpoint_failover_policy_enabled`: %+v", err)
		}

		d.Set("role", string(props.ReplicationRole))

		partnerManagedInstanceId, err := flattenSqlInstanceFailoverGroupManagedInstance(props.ManagedInstancePairs)
		if err != nil {
			return fmt.Errorf("flatten `partner_managed_instance_id`: %+v", err)
		}
		if err := d.Set("partner_managed_instance_id", partnerManagedInstanceId); err != nil {
			return fmt.Errorf("setting `partner_managed_instance_id`: %+v", err)
		}

		if err := d.Set("partner_region", flattenSqlInstanceFailoverGroupPartnerRegions(props.PartnerRegions)); err != nil {
			return fmt.Errorf("setting `partner_region`: %+v", err)
		}
	}

	return nil
}

func resourceSqlInstanceFailoverGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.InstanceFailoverGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.InstanceFailoverGroupID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.LocationName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return err
}

func expandSqlInstanceFailoverGroupReadWritePolicy(d *pluginsdk.ResourceData) *sql.InstanceFailoverGroupReadWriteEndpoint {
	vs := d.Get("read_write_endpoint_failover_policy").([]interface{})
	v := vs[0].(map[string]interface{})

	mode := sql.ReadWriteEndpointFailoverPolicy(v["mode"].(string))
	graceMins := int32(v["grace_minutes"].(int))

	policy := &sql.InstanceFailoverGroupReadWriteEndpoint{
		FailoverPolicy: mode,
	}

	if mode != sql.ReadWriteEndpointFailoverPolicyManual {
		policy.FailoverWithDataLossGracePeriodMinutes = utils.Int32(graceMins)
	}

	return policy
}

func flattenSqlInstanceFailoverGroupReadWritePolicy(input *sql.InstanceFailoverGroupReadWriteEndpoint) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	policy := make(map[string]interface{})

	policy["mode"] = string(input.FailoverPolicy)

	if input.FailoverWithDataLossGracePeriodMinutes != nil {
		policy["grace_minutes"] = *input.FailoverWithDataLossGracePeriodMinutes
	}
	return []interface{}{policy}
}

func expandSqlInstanceFailoverGroupReadOnlyPolicy(d *pluginsdk.ResourceData) *sql.InstanceFailoverGroupReadOnlyEndpoint {
	mode := sql.ReadOnlyEndpointFailoverPolicyDisabled
	if d.Get("readonly_endpoint_failover_policy_enabled").(bool) {
		mode = sql.ReadOnlyEndpointFailoverPolicyEnabled
	}
	return &sql.InstanceFailoverGroupReadOnlyEndpoint{
		FailoverPolicy: mode,
	}
}

func expandSqlInstanceFailoverGroupManagedInstanceId(d *pluginsdk.ResourceData, primaryID parse.ManagedInstanceId) *[]sql.ManagedInstancePairInfo {
	instanceId := d.Get("partner_managed_instance_id").(string)
	partners := make([]sql.ManagedInstancePairInfo, 0)

	partners = append(partners, sql.ManagedInstancePairInfo{
		PrimaryManagedInstanceID: utils.String(primaryID.ID()),
		PartnerManagedInstanceID: &instanceId,
	})

	return &partners
}

func flattenSqlInstanceFailoverGroupPrimaryInstance(input *[]sql.ManagedInstancePairInfo, role sql.InstanceFailoverGroupReplicationRole) (string, error) {
	id := ""
	if input != nil && len(*input) >= 1 {
		if managedInstancePairs := *input; managedInstancePairs[0].PrimaryManagedInstanceID != nil && role == sql.InstanceFailoverGroupReplicationRolePrimary {
			id = *managedInstancePairs[0].PrimaryManagedInstanceID
		}
		if managedInstancePairs := *input; managedInstancePairs[0].PartnerManagedInstanceID != nil && role == sql.InstanceFailoverGroupReplicationRoleSecondary {
			id = *managedInstancePairs[0].PartnerManagedInstanceID
		}
	}

	managedInstanceId, err := parse.ManagedInstanceID(id)
	if err != nil {
		return "", err
	}

	return managedInstanceId.Name, nil
}

func flattenSqlInstanceFailoverGroupManagedInstance(input *[]sql.ManagedInstancePairInfo) (string, error) {
	if input == nil || len(*input) != 1 || (*input)[0].PartnerManagedInstanceID == nil {
		return "", fmt.Errorf("invalid number of `partner_managed_instance_id` instances found")
	}

	return *(*input)[0].PartnerManagedInstanceID, nil
}

func flattenSqlInstanceFailoverGroupPartnerRegions(input *[]sql.PartnerRegionInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if input != nil {
		for _, region := range *input {
			result = append(result, map[string]interface{}{
				"location": location.NormalizeNilable(region.Location),
				"role":     string(region.ReplicationRole),
			})
		}
	}
	return result
}
