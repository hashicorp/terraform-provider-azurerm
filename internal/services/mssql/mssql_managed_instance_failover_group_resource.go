package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlManagedInstanceFailoverGroupModel struct {
	Name                                  string `tfschema:"name"`
	Location                              string `tfschema:"location"`
	ManagedInstanceId                     string `tfschema:"managed_instance_id"`
	PartnerManagedInstanceId              string `tfschema:"partner_managed_instance_id"`
	ReadOnlyEndpointFailoverPolicyEnabled bool   `tfschema:"readonly_endpoint_failover_policy_enabled"`

	ReadWriteEndpointFailurePolicy []MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel `tfschema:"read_write_endpoint_failover_policy"`

	PartnerRegion []MsSqlManagedInstancePartnerRegionModel `tfschema:"partner_region"`
	Role          string                                   `tfschema:"role"`
}

type MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel struct {
	GraceMinutes int32  `tfschema:"grace_minutes"`
	Mode         string `tfschema:"mode"`
}

type MsSqlManagedInstancePartnerRegionModel struct {
	Location string `tfschema:"location"`
	Role     string `tfschema:"role"`
}

var _ sdk.Resource = MsSqlManagedInstanceFailoverGroupResource{}
var _ sdk.ResourceWithUpdate = MsSqlManagedInstanceFailoverGroupResource{}

type MsSqlManagedInstanceFailoverGroupResource struct{}

func (r MsSqlManagedInstanceFailoverGroupResource) ResourceType() string {
	return "azurerm_mssql_managed_instance_failover_group"
}

func (r MsSqlManagedInstanceFailoverGroupResource) ModelObject() interface{} {
	return &MsSqlManagedInstanceFailoverGroupModel{}
}

func (r MsSqlManagedInstanceFailoverGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.InstanceFailoverGroupID
}

func (r MsSqlManagedInstanceFailoverGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlFailoverGroupName,
		},

		"location": commonschema.Location(),

		"managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedInstanceID,
		},

		"partner_managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
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
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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

		"role": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.InstanceFailoverGroupsClient
			instancesClient := metadata.Client.MSSQL.ManagedInstancesClient

			var model MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			managedInstanceId, err := parse.ManagedInstanceID(model.ManagedInstanceId)
			if err != nil {
				return fmt.Errorf("parsing `managed_instance_id`: %v", err)
			}

			id := parse.NewInstanceFailoverGroupID(managedInstanceId.SubscriptionId,
				managedInstanceId.ResourceGroup, model.Location, model.Name)

			partnerId, err := parse.ManagedInstanceID(model.PartnerManagedInstanceId)
			if err != nil {
				return err
			}

			partner, err := instancesClient.Get(ctx, partnerId.ResourceGroup, partnerId.Name, "")
			if err != nil || partner.Location == nil || *partner.Location == "" {
				return fmt.Errorf("checking for existence and region of Partner of %q: %+v", id, err)
			}

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id.ResourceGroup, id.LocationName, id.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			readOnlyFailoverPolicy := sql.ReadOnlyEndpointFailoverPolicyDisabled
			if model.ReadOnlyEndpointFailoverPolicyEnabled {
				readOnlyFailoverPolicy = sql.ReadOnlyEndpointFailoverPolicyEnabled
			}

			parameters := sql.InstanceFailoverGroup{
				InstanceFailoverGroupProperties: &sql.InstanceFailoverGroupProperties{
					ReadOnlyEndpoint: &sql.InstanceFailoverGroupReadOnlyEndpoint{
						FailoverPolicy: readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: &sql.InstanceFailoverGroupReadWriteEndpoint{},
					PartnerRegions: &[]sql.PartnerRegionInfo{
						{
							Location: partner.Location,
						},
					},
					ManagedInstancePairs: &[]sql.ManagedInstancePairInfo{
						{
							PrimaryManagedInstanceID: utils.String(managedInstanceId.ID()),
							PartnerManagedInstanceID: utils.String(partnerId.ID()),
						},
					},
				},
			}

			if rwPolicy := model.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				parameters.InstanceFailoverGroupProperties.ReadWriteEndpoint.FailoverPolicy = sql.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(sql.ReadWriteEndpointFailoverPolicyAutomatic) {
					parameters.InstanceFailoverGroupProperties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = utils.Int32(rwPolicy[0].GraceMinutes)
				}
			}

			metadata.Logger.Infof("Creating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LocationName, id.Name, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.InstanceFailoverGroupsClient
			instancesClient := metadata.Client.MSSQL.ManagedInstancesClient

			id, err := parse.InstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			managedInstanceId, err := parse.ManagedInstanceID(state.ManagedInstanceId)
			if err != nil {
				return fmt.Errorf("parsing `managed_instance_id`: %v", err)
			}

			partnerId, err := parse.ManagedInstanceID(state.PartnerManagedInstanceId)
			if err != nil {
				return err
			}

			partner, err := instancesClient.Get(ctx, partnerId.ResourceGroup, partnerId.Name, "")
			if err != nil || partner.Location == nil || *partner.Location == "" {
				return fmt.Errorf("checking for existence and region of Partner of %q: %+v", id, err)
			}

			readOnlyFailoverPolicy := sql.ReadOnlyEndpointFailoverPolicyDisabled
			if state.ReadOnlyEndpointFailoverPolicyEnabled {
				readOnlyFailoverPolicy = sql.ReadOnlyEndpointFailoverPolicyEnabled
			}

			parameters := sql.InstanceFailoverGroup{
				InstanceFailoverGroupProperties: &sql.InstanceFailoverGroupProperties{
					ReadOnlyEndpoint: &sql.InstanceFailoverGroupReadOnlyEndpoint{
						FailoverPolicy: readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: &sql.InstanceFailoverGroupReadWriteEndpoint{},
					PartnerRegions: &[]sql.PartnerRegionInfo{
						{
							Location: partner.Location,
						},
					},
					ManagedInstancePairs: &[]sql.ManagedInstancePairInfo{
						{
							PrimaryManagedInstanceID: utils.String(managedInstanceId.ID()),
							PartnerManagedInstanceID: utils.String(partnerId.ID()),
						},
					},
				},
			}

			if rwPolicy := state.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				parameters.InstanceFailoverGroupProperties.ReadWriteEndpoint.FailoverPolicy = sql.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(sql.ReadWriteEndpointFailoverPolicyAutomatic) {
					parameters.InstanceFailoverGroupProperties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = utils.Int32(rwPolicy[0].GraceMinutes)
				}
			}

			metadata.Logger.Infof("Updating %s", id)

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LocationName, id.Name, parameters)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.InstanceFailoverGroupsClient

			id, err := parse.InstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			result, err := client.Get(ctx, id.ResourceGroup, id.LocationName, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(result.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceFailoverGroupModel{
				Name:     id.Name,
				Location: id.LocationName,
			}

			if props := result.InstanceFailoverGroupProperties; props != nil {
				model.Role = string(props.ReplicationRole)

				if instancePairs := props.ManagedInstancePairs; instancePairs != nil && len(*instancePairs) == 1 {
					if primaryId := (*instancePairs)[0].PrimaryManagedInstanceID; primaryId != nil {
						id, err := parse.ManagedInstanceID(*primaryId)
						if err != nil {
							return fmt.Errorf("parsing `PrimaryManagedInstanceID` from response: %v", err)
						}

						model.ManagedInstanceId = id.ID()
					}

					if partnerId := (*instancePairs)[0].PartnerManagedInstanceID; partnerId != nil {
						id, err := parse.ManagedInstanceID(*partnerId)
						if err != nil {
							return fmt.Errorf("parsing `PrimaryManagedInstanceID` from response: %v", err)
						}

						model.PartnerManagedInstanceId = id.ID()
					}
				}

				if partnerRegions := props.PartnerRegions; partnerRegions != nil {
					for _, partnerRegion := range *partnerRegions {
						var location string
						if partnerRegion.Location != nil {
							location = *partnerRegion.Location
						}

						model.PartnerRegion = append(model.PartnerRegion, MsSqlManagedInstancePartnerRegionModel{
							Location: location,
							Role:     string(partnerRegion.ReplicationRole),
						})
					}
				}

				if readOnlyEndpoint := props.ReadOnlyEndpoint; readOnlyEndpoint != nil {
					if readOnlyEndpoint.FailoverPolicy == sql.ReadOnlyEndpointFailoverPolicyEnabled {
						model.ReadOnlyEndpointFailoverPolicyEnabled = true
					}
				}

				if readWriteEndpoint := props.ReadWriteEndpoint; readWriteEndpoint != nil {
					var graceMinutes int32
					if readWriteEndpoint.FailoverWithDataLossGracePeriodMinutes != nil {
						graceMinutes = *readWriteEndpoint.FailoverWithDataLossGracePeriodMinutes
					}

					model.ReadWriteEndpointFailurePolicy = []MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel{
						{
							Mode:         string(readWriteEndpoint.FailoverPolicy),
							GraceMinutes: graceMinutes,
						},
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r MsSqlManagedInstanceFailoverGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.InstanceFailoverGroupsClient

			id, err := parse.InstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.LocationName, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
			}

			return nil
		},
	}
}
