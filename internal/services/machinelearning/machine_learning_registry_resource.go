package machinelearning

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/registrymanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// PublicNetworkAccess values not defined as enum, but given in the description
const (
	PublicNetworkAccessEnabled  = "Enabled"
	PublicNetworkAccessDisabled = "Disabled"
)

func publicNetworkAccessFromBool(b bool) *string {
	if b {
		return pointer.To(PublicNetworkAccessEnabled)
	}

	return pointer.To(PublicNetworkAccessDisabled)
}

const (
	ACRAccountSKUPremium = "Premium"
)

//	type AcrDetails struct {
//		UserCreatedACRAccountARMResourceID string `tfschema:"user_created_acr_account_arm_resource_id"`
//	}
type ArmResourceId string

//
// type RegionDetails struct {
// 	Location string `tfschema:"location"`
//
// 	AcrDetails            []AcrDetails            `tfschema:"acr_details"`
// 	StorageAccountDetails []StorageAccountDetails `tfschema:"storage_account_details"`
// }
//
// type StorageAccountDetails struct {
// 	UserCreatedStorageAccountARMResourceID string `tfschema:"user_created_storage_account_arm_resource_id"`
//
// 	StorageAccountType       string        `tfschema:"system_created_storage_account_type"`
// 	StorageAccountHnsEnabled bool          `tfschema:"system_created_storage_account_hns_enabled"`
// 	StorageAccountName       string        `tfschema:"system_created_storage_account_name"`
// 	ArmResourceId            ArmResourceId `tfschema:"system_created_storage_account_arm_resource_id"`
// }
//
// type SystemCreatedAcrAccount struct {
// 	AcrAccountName string        `tfschema:"acr_account_name"`
// 	AcrAccountSku  string        `tfschema:"acr_account_sku"`
// 	ArmResourceId  ArmResourceId `tfschema:"arm_resource_id"`
// }
//
// type SystemCreatedStorageAccount struct {
// 	StorageAccountName       string        `tfschema:"name"`
// 	StorageAccountType       string        `tfschema:"type"`
// 	StorageAccountHnsEnabled bool          `tfschema:"hns_enabled"`
// 	ArmResourceId            ArmResourceId `tfschema:"arm_resource_id"`
// }

type ReplicationRegion struct {
	Location                 string        `tfschema:"location"`
	ContainerRegistryID      ArmResourceId `tfschema:"container_registry_id"`
	StorageAccountID         ArmResourceId `tfschema:"storage_account_id"`
	StorageAccountType       string        `tfschema:"storage_account_type"`
	StorageAccountHNSEnabled bool          `tfschema:"storage_account_hns_enabled"`
}

type MachineLearningRegistryModel struct {
	ResourceGroupName             string                                     `tfschema:"resource_group_name"`
	Name                          string                                     `tfschema:"name"`
	Location                      string                                     `tfschema:"location"`
	DiscoveryUrl                  string                                     `tfschema:"discovery_url"`
	ContainerRegistryID           ArmResourceId                              `tfschema:"container_registry_id"`
	StorageAccountID              ArmResourceId                              `tfschema:"storage_account_id"`
	StorageAccountType            string                                     `tfschema:"storage_account_type"`
	StorageAccountHNSEnabled      bool                                       `tfschema:"storage_account_hns_enabled"`
	ReplicationRegions            []ReplicationRegion                        `tfschema:"replication_regions"`
	MlFlowRegistryUri             string                                     `tfschema:"ml_flow_registry_uri"`
	IntellectualPropertyPublisher string                                     `tfschema:"intellectual_property_publisher"`
	PublicNetworkAccess           bool                                       `tfschema:"public_network_access_enabled"`
	Tags                          map[string]string                          `tfschema:"tags"`
	Identity                      []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	ManagedResourceGroup          ArmResourceId                              `tfschema:"managed_resource_group"`
	// RegionDetails                 []RegionDetails                            `tfschema:"region_details"`
}

type MachineLearningRegistryResource struct{}

var _ sdk.ResourceWithUpdate = (*MachineLearningRegistryResource)(nil)

func (m MachineLearningRegistryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-z0-9][a-zA-z0-9_-]{2,32}$`), "Registry name must be between 3 and 33 characters long. Its first character has to be alphanumeric, and the rest may contain hyphens and underscores. No whitespace is allowed."),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptionalForceNew(),

		"location": commonschema.Location(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},

		"storage_account_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateStorageAccountID,
			ConflictsWith: []string{"storage_account_type", "storage_account_hns_enabled"},
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard_LRS",
				"Standard_GRS",
				"Standard_RAGRS",
				"Standard_ZRS",
				"Standard_GZRS",
				"Standard_RAGZRS",
				"Premium_LRS",
				"Premium_ZRS",
			}, false),
			ConflictsWith: []string{"storage_account_id"},
		},

		"storage_account_hns_enabled": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"storage_account_id"},
		},

		"replication_regions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// This might only be conditionally ForceNew? Please investigate the API to find out
			// ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_registry_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: registries.ValidateRegistryID,
					},

					"location": commonschema.LocationWithoutForceNew(),

					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
						// ConflictsWith: []string{
						// 	"replication_regions.0.system_created_storage_account_type",
						// 	"replication_regions.0.storage_account_hns_enabled",
						// },
					},

					"storage_account_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Standard_LRS",
							"Standard_GRS",
							"Standard_RAGRS",
							"Standard_ZRS",
							"Standard_GZRS",
							"Standard_RAGZRS",
							"Premium_LRS",
							"Premium_ZRS",
						}, false),
						// ConflictsWith: []string{"replication_regions.0.storage_account_id"},
					},

					"storage_account_hns_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
						// ConflictsWith: []string{"replication_regions.0.storage_account_id"},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (m MachineLearningRegistryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"intellectual_property_publisher": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ml_flow_registry_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// API will return a generated mrg even provided in request
		"managed_resource_group": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m MachineLearningRegistryResource) ModelObject() interface{} {
	return &MachineLearningRegistryModel{}
}

func (m MachineLearningRegistryResource) ResourceType() string {
	return "azurerm_machine_learning_registry"
}

func (m MachineLearningRegistryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.MachineLearning.RegistryManagementClient

			var model MachineLearningRegistryModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := registrymanagement.NewRegistryID(subscriptionID, model.ResourceGroupName, model.Name)
			existing, err := client.RegistriesGet(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := registrymanagement.RegistryTrackedResource{
				Location: model.Location,
				Name:     &model.Name,
				Tags:     &model.Tags,
			}

			if req.Identity, err = identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity); err != nil {
				return fmt.Errorf("expand identity: %+v", err)
			}

			req.Properties = registrymanagement.Registry{
				PublicNetworkAccess: publicNetworkAccessFromBool(model.PublicNetworkAccess),
				RegionDetails:       m.expandRegionDetails(&model),
			}

			if err = client.RegistriesCreateOrUpdateThenPoll(ctx, id, req); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m MachineLearningRegistryResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.MachineLearning.RegistryManagementClient
			id, err := registrymanagement.ParseRegistryID(meta.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing registry ID: %+V", err)
			}

			existing, err := client.RegistriesGet(ctx, *id)
			if response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("registry %s not found to update", *id)
			}
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			var model MachineLearningRegistryModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			if meta.ResourceData.HasChange("replication_regions") {
				o, n := meta.ResourceData.GetChange("replication_regions")
				remove, add := m.regionsDiff(o, n)
				// remove first and then update
				req := registrymanagement.RegistryTrackedResource{
					Location: model.Location,
					Name:     &model.Name,
					Tags:     &model.Tags,
				}
				if req.Identity, err = identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity); err != nil {
					return fmt.Errorf("expand identity: %+v", err)
				}
				req.Properties = registrymanagement.Registry{
					PublicNetworkAccess: publicNetworkAccessFromBool(model.PublicNetworkAccess),
				}

				if len(remove) > 0 {
					var regions []registrymanagement.RegistryRegionArmDetails
					for _, region := range remove {
						regions = append(regions, m.buildRegionDetail(ReplicationRegion{
							Location: region.(map[string]interface{})["location"].(string),
						}))
					}
					req.Properties.RegionDetails = &regions
					if err = client.RegistriesRemoveRegionsThenPoll(ctx, *id, req); err != nil {
						return fmt.Errorf("remove region: %+v", err)
					}
				}
				if len(add) > 0 {
					req.Properties.RegionDetails = m.expandRegionDetails(&model)
					if err = client.RegistriesCreateOrUpdateThenPoll(ctx, *id, req); err != nil {
						return fmt.Errorf("remove region: %+v", err)
					}
				}
				// region list can be delay, need a wait state
				stateConf := pluginsdk.StateChangeConf{
					Delay:   0,
					Pending: []string{"Pending"},
					Refresh: func() (result interface{}, state string, err error) {
						if resp, err := client.RegistriesGet(ctx, *id); err != nil {
							return resp.Model, "Done", err
						} else {
							if getModel := resp.Model; getModel != nil && getModel.Properties.RegionDetails != nil {
								if len(*getModel.Properties.RegionDetails) == len(model.ReplicationRegions)+1 {
									return resp.Model, "Done", nil
								}
							}
						}
						return nil, "Pending", nil
					},
					Target:     []string{"Done"},
					Timeout:    time.Second * 10, // usually it maybe 5-seconds delay under my tests
					MinTimeout: time.Second * 3,
				}
				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("wait update replication regions: %+v", err)
				}
			}

			if meta.ResourceData.HasChanges("tags") {
				req := registrymanagement.PartialRegistryPartialTrackedResource{}
				req.Tags = pointer.To(model.Tags)

				if _, err = client.RegistriesUpdate(ctx, *id, req); err != nil {
					return fmt.Errorf("updating %q: %v", *id, err)
				}

				// remove wait state when issue fixed: https://github.com/Azure/azure-rest-api-specs/issues/25200
				// with this stateWait, it can still read old value if from different server cluster.
				deadline, ok := ctx.Deadline()
				if !ok {
					return fmt.Errorf("internal error: no deadline for ctx")
				}
				stateConf := &pluginsdk.StateChangeConf{
					Pending: []string{"Pending"},
					Target:  []string{"Finish"},
					Refresh: func() (result interface{}, state string, err error) {
						got, err := client.RegistriesGet(ctx, *id)
						if err != nil {
							return nil, "", fmt.Errorf("get updated registry: %+v", err)
						}
						state = "Pending"
						result = got
						if got.Model != nil {
							tag := got.Model.Tags
							if len(model.Tags) == 0 {
								if tag == nil || len(*tag) == 0 {
									state = "Finish"
								}
							} else if tag != nil && reflect.DeepEqual(model.Tags, *tag) {
								state = "Finish"
							}
						}
						return
					},
					MinTimeout: 10 * time.Second,
					Timeout:    time.Until(deadline),
				}
				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("wait tags updated: %+v", err)
				}
			}

			return nil
		},
	}
}

func (m MachineLearningRegistryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := registrymanagement.ParseRegistryID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.MachineLearning.RegistryManagementClient
			result, err := client.RegistriesGet(ctx, *id)
			if err != nil {
				return err
			}
			if result.Model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			model := result.Model
			prop := model.Properties
			output := MachineLearningRegistryModel{
				ResourceGroupName:             id.ResourceGroupName,
				Name:                          id.RegistryName,
				Tags:                          pointer.From(model.Tags),
				DiscoveryUrl:                  pointer.From(prop.DiscoveryUrl),
				MlFlowRegistryUri:             pointer.From(prop.MlFlowRegistryUri),
				PublicNetworkAccess:           pointer.From(prop.PublicNetworkAccess) == PublicNetworkAccessEnabled,
				IntellectualPropertyPublisher: pointer.From(prop.IntellectualPropertyPublisher),
				Location:                      model.Location,
			}

			output.Identity, err = identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("faltten legacy identity: %+v", err)
			}

			output.ManagedResourceGroup = m.flattenArmResourceID(prop.ManagedResourceGroup)
			primary, replicas := m.flattenRegionDetails(model.Location, prop.RegionDetails)
			output.StorageAccountID = primary.StorageAccountID
			output.StorageAccountType = primary.StorageAccountType
			output.StorageAccountHNSEnabled = primary.StorageAccountHNSEnabled
			output.ContainerRegistryID = primary.ContainerRegistryID

			output.ReplicationRegions = replicas

			return meta.Encode(&output)
		},
	}
}

func (m MachineLearningRegistryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := registrymanagement.ParseRegistryID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.MachineLearning.RegistryManagementClient
			if err = client.RegistriesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m MachineLearningRegistryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return registrymanagement.ValidateRegistryID
}

func (m MachineLearningRegistryResource) expandArmResourceID(id ArmResourceId) *registrymanagement.ArmResourceId {
	if id == "" {
		return nil
	}

	return &registrymanagement.ArmResourceId{
		ResourceId: pointer.To(string(id)),
	}
}

func (m MachineLearningRegistryResource) buildRegionDetail(detail ReplicationRegion) registrymanagement.RegistryRegionArmDetails {
	item := registrymanagement.RegistryRegionArmDetails{
		Location: pointer.To(detail.Location),
		// AcrDetails:            m.expandAcrDetails(detail.AcrDetails),
	}

	if detail.StorageAccountID != "" {
		item.StorageAccountDetails = &[]registrymanagement.StorageAccountDetails{
			{
				UserCreatedStorageAccount: &registrymanagement.UserCreatedStorageAccount{
					ArmResourceId: m.expandArmResourceID(detail.StorageAccountID),
				},
			},
		}
	} else {
		item.StorageAccountDetails = &[]registrymanagement.StorageAccountDetails{
			{
				SystemCreatedStorageAccount: &registrymanagement.SystemCreatedStorageAccount{
					StorageAccountHnsEnabled: pointer.To(detail.StorageAccountHNSEnabled),
					StorageAccountType:       pointer.To(detail.StorageAccountType),
				},
			},
		}
	}

	if detail.ContainerRegistryID != "" {
		item.AcrDetails = &[]registrymanagement.AcrDetails{
			{
				UserCreatedAcrAccount: &registrymanagement.UserCreatedAcrAccount{
					ArmResourceId: m.expandArmResourceID(detail.ContainerRegistryID),
				},
			},
		}
	} else {
		item.AcrDetails = &[]registrymanagement.AcrDetails{
			{
				SystemCreatedAcrAccount: &registrymanagement.SystemCreatedAcrAccount{
					AcrAccountSku: pointer.To(ACRAccountSKUPremium),
				},
			},
		}
	}

	return item
}

func (m MachineLearningRegistryResource) expandRegionDetails(model *MachineLearningRegistryModel) *[]registrymanagement.RegistryRegionArmDetails {
	var output []registrymanagement.RegistryRegionArmDetails
	// add primary location first
	output = append(output, m.buildRegionDetail(ReplicationRegion{
		Location:                 model.Location,
		ContainerRegistryID:      model.ContainerRegistryID,
		StorageAccountID:         model.StorageAccountID,
		StorageAccountType:       model.StorageAccountType,
		StorageAccountHNSEnabled: model.StorageAccountHNSEnabled,
	}))

	for _, detail := range model.ReplicationRegions {
		output = append(output, m.buildRegionDetail(detail))
	}

	return &output
}

func (m MachineLearningRegistryResource) flattenArmResourceID(id *registrymanagement.ArmResourceId) ArmResourceId {
	if nil == (id) {
		return ""
	}

	return ArmResourceId(pointer.From(id.ResourceId))
}

func (m MachineLearningRegistryResource) flattenRegion(detail registrymanagement.RegistryRegionArmDetails) ReplicationRegion {
	var res ReplicationRegion
	res.Location = pointer.From(detail.Location)
	if detail.AcrDetails != nil && len(*detail.AcrDetails) > 0 {
		acr := (*detail.AcrDetails)[0]
		if acr.UserCreatedAcrAccount != nil {
			res.ContainerRegistryID = m.flattenArmResourceID(acr.UserCreatedAcrAccount.ArmResourceId)
		} else if acr.SystemCreatedAcrAccount != nil {
		}
	}

	if detail.StorageAccountDetails != nil && len(*detail.StorageAccountDetails) > 0 {
		storage := (*detail.StorageAccountDetails)[0]
		if storage.UserCreatedStorageAccount != nil {
			res.StorageAccountID = m.flattenArmResourceID(storage.UserCreatedStorageAccount.ArmResourceId)
		} else if storage.SystemCreatedStorageAccount != nil {
			res.StorageAccountType = pointer.From(storage.SystemCreatedStorageAccount.StorageAccountType)
			res.StorageAccountHNSEnabled = pointer.From(storage.SystemCreatedStorageAccount.StorageAccountHnsEnabled)
		}
	}

	return res
}

func (m MachineLearningRegistryResource) flattenRegionDetails(primaryLocation string, details *[]registrymanagement.RegistryRegionArmDetails) (
	primary ReplicationRegion, replicas []ReplicationRegion) {
	if details == nil || len(*details) == 0 {
		return ReplicationRegion{}, nil
	}

	for _, detail := range *details {
		if pointer.From(detail.Location) == primaryLocation {
			primary = m.flattenRegion(detail)
		} else {
			replicas = append(replicas, m.flattenRegion(detail))
		}
	}
	return primary, replicas
}

func (m MachineLearningRegistryResource) regionsDiff(old, new interface{}) (remove, add []interface{}) {

	oldRegions := m.regionsToMap(old)
	newRegions := m.regionsToMap(new)

	// remove if exists in old, but not exists in new
	for location, v := range oldRegions {
		if _, ok := newRegions[location]; !ok {
			remove = append(remove, v)
		}
	}

	for location, v := range newRegions {
		if _, ok := oldRegions[location]; !ok {
			add = append(add, v)
		}
	}
	return
}

func (m MachineLearningRegistryResource) regionsToMap(obj interface{}) map[string]interface{} {
	v, ok := obj.([]interface{})
	if !ok {
		return nil
	}

	if len(v) == 0 {
		return nil
	}
	res := map[string]interface{}{}

	for _, value := range v {
		if value == nil {
			continue
		}
		if region, ok := value.(map[string]interface{}); ok {
			res[region["location"].(string)] = value
		}
	}
	return res
}
