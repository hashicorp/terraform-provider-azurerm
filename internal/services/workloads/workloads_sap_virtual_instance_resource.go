package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPVirtualInstanceModel struct {
	Name                          string                          `tfschema:"name"`
	ResourceGroupName             string                          `tfschema:"resource_group_name"`
	Location                      string                          `tfschema:"location"`
	DeploymentWithOSConfiguration []DeploymentWithOSConfiguration `tfschema:"deployment_with_os_configuration"`
	DiscoveryConfiguration        []DiscoveryConfiguration        `tfschema:"discovery_configuration"`
	Environment                   string                          `tfschema:"environment"`
	Identity                      []identity.ModelUserAssigned    `tfschema:"identity"`
	ManagedResourceGroupName      string                          `tfschema:"managed_resource_group_name"`
	SapProduct                    string                          `tfschema:"sap_product"`
	Tags                          map[string]string               `tfschema:"tags"`
}

type DiscoveryConfiguration struct {
	CentralServerVmId         string `tfschema:"central_server_vm_id"`
	ManagedStorageAccountName string `tfschema:"managed_storage_account_name"`
}

type DeploymentWithOSConfiguration struct {
	AppLocation               string                      `tfschema:"app_location"`
	OsSapConfiguration        []OsSapConfiguration        `tfschema:"os_sap_configuration"`
	SingleServerConfiguration []SingleServerConfiguration `tfschema:"single_server_configuration"`
}

type OsSapConfiguration struct {
	DeployerVmPackages []DeployerVmPackages `tfschema:"deployer_vm_packages"`
	SapFqdn            string               `tfschema:"sap_fqdn"`
}

type DeployerVmPackages struct {
	StorageAccountId string `tfschema:"storage_account_id"`
	Url              string `tfschema:"url"`
}

type SingleServerConfiguration struct {
	AppResourceGroupName            string                            `tfschema:"app_resource_group_name"`
	DatabaseType                    string                            `tfschema:"database_type"`
	DiskVolumeConfigurations        []DiskVolumeConfiguration         `tfschema:"disk_volume_configuration"`
	IsSecondaryIpEnabled            bool                              `tfschema:"is_secondary_ip_enabled"`
	SubnetId                        string                            `tfschema:"subnet_id"`
	VirtualMachineConfiguration     []VirtualMachineConfiguration     `tfschema:"virtual_machine_configuration"`
	VirtualMachineFullResourceNames []VirtualMachineFullResourceNames `tfschema:"virtual_machine_full_resource_names"`
}

type DiskVolumeConfiguration struct {
	VolumeName string `tfschema:"volume_name"`
	Count      int64  `tfschema:"count"`
	SizeGb     int64  `tfschema:"size_gb"`
	SkuName    string `tfschema:"sku_name"`
}

type VirtualMachineConfiguration struct {
	ImageReference []ImageReference `tfschema:"image_reference"`
	OSProfile      []OSProfile      `tfschema:"os_profile"`
	VmSize         string           `tfschema:"vm_size"`
}

type ImageReference struct {
	Offer     string `tfschema:"offer"`
	Publisher string `tfschema:"publisher"`
	Sku       string `tfschema:"sku"`
	Version   string `tfschema:"version"`
}

type OSProfile struct {
	AdminPassword      string               `tfschema:"admin_password"`
	AdminUsername      string               `tfschema:"admin_username"`
	LinuxConfiguration []LinuxConfiguration `tfschema:"linux_configuration"`
}

type LinuxConfiguration struct {
	PasswordAuthenticationEnabled bool         `tfschema:"password_authentication_enabled"`
	SshKeyPair                    []SshKeyPair `tfschema:"ssh_key_pair"`
	SshPublicKeyDataList          []string     `tfschema:"ssh_public_key_data_list"`
}

type SshKeyPair struct {
	PrivateKey string `tfschema:"private_key"`
	PublicKey  string `tfschema:"public_key"`
}

type VirtualMachineFullResourceNames struct {
	DataDiskNames         map[string]interface{} `tfschema:"data_disk_names"`
	HostName              string                 `tfschema:"host_name"`
	NetworkInterfaceNames []string               `tfschema:"network_interface_names"`
	OSDiskName            string                 `tfschema:"os_disk_name"`
	VMName                string                 `tfschema:"vm_name"`
}

type WorkloadsSAPVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPVirtualInstanceResource{}

func (r WorkloadsSAPVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_virtual_instance"
}

func (r WorkloadsSAPVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPVirtualInstanceModel{}
}

func (r WorkloadsSAPVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}

func (r WorkloadsSAPVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"deployment_with_os_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_location": commonschema.LocationOptional(),

					"os_sap_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"deployer_vm_packages": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"storage_account_id": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validate.StorageAccountID,
											},

											"url": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ForceNew:     true,
												ValidateFunc: validation.IsURLWithHTTPorHTTPS,
											},
										},
									},
								},

								"sap_fqdn": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"single_server_configuration": SchemaForSAPVirtualInstanceSingleServerConfiguration(),
				},
			},
			ConflictsWith: []string{"discovery_configuration"},
		},

		"discovery_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"central_server_vm_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"managed_storage_account_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ConflictsWith: []string{"deployment_with_os_configuration"},
		},

		"environment": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPEnvironmentTypeNonProd),
				string(sapvirtualinstances.SAPEnvironmentTypeProd),
			}, false),
		},

		"sap_product": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPProductTypeECC),
				string(sapvirtualinstances.SAPProductTypeOther),
				string(sapvirtualinstances.SAPProductTypeSFourHANA),
			}, false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSAPVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPVirtualInstances
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapvirtualinstances.NewSapVirtualInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identity, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identity,
				Location: location.Normalize(model.Location),
				Properties: sapvirtualinstances.SAPVirtualInstanceProperties{
					Environment: sapvirtualinstances.SAPEnvironmentType(model.Environment),
					SapProduct:  sapvirtualinstances.SAPProductType(model.SapProduct),
				},
				Tags: &model.Tags,
			}

			if v := model.DeploymentWithOSConfiguration; v != nil {
				parameters.Properties.Configuration = expandDeploymentWithOSConfiguration(v)
			}

			if v := model.DiscoveryConfiguration; v != nil {
				parameters.Properties.Configuration = expandDiscoveryConfiguration(v)
			}

			if v := model.ManagedResourceGroupName; v != "" {
				parameters.Properties.ManagedResourceGroupConfiguration = &sapvirtualinstances.ManagedRGConfiguration{
					Name: utils.String(v),
				}
			}

			if err := client.CreateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &sapvirtualinstances.UpdateSAPVirtualInstanceRequest{}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			if _, err := client.Update(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := WorkloadsSAPVirtualInstanceModel{
				Name:              id.SapVirtualInstanceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identity, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			properties := &model.Properties
			state.Environment = string(properties.Environment)
			state.SapProduct = string(properties.SapProduct)

			if properties.Configuration != nil {
				if v, ok := properties.Configuration.(sapvirtualinstances.DeploymentWithOSConfiguration); ok {
					state.DeploymentWithOSConfiguration = flattenDeploymentWithOSConfiguration(&v)
				}

				if v, ok := properties.Configuration.(sapvirtualinstances.DiscoveryConfiguration); ok {
					state.DiscoveryConfiguration = flattenDiscoveryConfiguration(&v)
				}
			}

			if v := properties.ManagedResourceGroupConfiguration; v != nil && v.Name != nil {
				state.ManagedResourceGroupName = *v.Name
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstances

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDiscoveryConfiguration(input []DiscoveryConfiguration) *sapvirtualinstances.DiscoveryConfiguration {
	if len(input) == 0 {
		return nil
	}

	configuration := &input[0]

	result := sapvirtualinstances.DiscoveryConfiguration{}

	if v := configuration.CentralServerVmId; v != "" {
		result.CentralServerVMId = utils.String(v)
	}

	if v := configuration.ManagedStorageAccountName; v != "" {
		result.ManagedRgStorageAccountName = utils.String(v)
	}

	return &result
}

func flattenDiscoveryConfiguration(input *sapvirtualinstances.DiscoveryConfiguration) []DiscoveryConfiguration {
	if input == nil {
		return nil
	}

	result := DiscoveryConfiguration{}

	if v := input.CentralServerVMId; v != nil {
		result.CentralServerVmId = *v
	}

	if v := input.ManagedRgStorageAccountName; v != nil {
		result.ManagedStorageAccountName = *v
	}

	return []DiscoveryConfiguration{
		result,
	}
}

func expandDeploymentWithOSConfiguration(input []DeploymentWithOSConfiguration) *sapvirtualinstances.DeploymentWithOSConfiguration {
	if len(input) == 0 {
		return nil
	}

	configuration := &input[0]

	result := sapvirtualinstances.DeploymentWithOSConfiguration{}

	if v := configuration.AppLocation; v != "" {
		result.AppLocation = utils.String(v)
	}

	if v := configuration.OsSapConfiguration; v != nil {
		result.OsSapConfiguration = expandOsSapConfiguration(v)
	}

	if v := configuration.SingleServerConfiguration; v != nil {
		result.InfrastructureConfiguration = expandSingleServerConfiguration(v)
	}

	return &result
}

func expandOsSapConfiguration(input []OsSapConfiguration) *sapvirtualinstances.OsSapConfiguration {
	if len(input) == 0 {
		return nil
	}

	osSapConfiguration := &input[0]

	result := sapvirtualinstances.OsSapConfiguration{}

	if v := osSapConfiguration.DeployerVmPackages; v != nil {
		result.DeployerVMPackages = expandDeployerVmPackages(v)
	}

	if v := osSapConfiguration.SapFqdn; v != "" {
		result.SapFqdn = utils.String(v)
	}

	return &result
}

func expandDeployerVmPackages(input []DeployerVmPackages) *sapvirtualinstances.DeployerVMPackages {
	if len(input) == 0 {
		return nil
	}

	deployerVmPackages := &input[0]

	result := sapvirtualinstances.DeployerVMPackages{}

	if v := deployerVmPackages.StorageAccountId; v != "" {
		result.StorageAccountId = utils.String(v)
	}

	if v := deployerVmPackages.Url; v != "" {
		result.Url = utils.String(v)
	}

	return &result
}

func flattenDeploymentWithOSConfiguration(input *sapvirtualinstances.DeploymentWithOSConfiguration) []DeploymentWithOSConfiguration {
	if input == nil {
		return nil
	}

	result := DeploymentWithOSConfiguration{}

	if v := input.AppLocation; v != nil {
		result.AppLocation = *v
	}

	if v := input.OsSapConfiguration; v != nil {
		result.OsSapConfiguration = flattenOsSapConfiguration(v)
	}

	if configuration := input.InfrastructureConfiguration; configuration != nil {
		if v, ok := configuration.(sapvirtualinstances.SingleServerConfiguration); ok {
			result.SingleServerConfiguration = flattenSingleServerConfiguration(v)
		}
	}

	return []DeploymentWithOSConfiguration{
		result,
	}
}

func flattenOsSapConfiguration(input *sapvirtualinstances.OsSapConfiguration) []OsSapConfiguration {
	if input == nil {
		return nil
	}

	result := OsSapConfiguration{}

	if v := input.DeployerVMPackages; v != nil {
		result.DeployerVmPackages = flattenDeployerVMPackages(v)
	}

	if v := input.SapFqdn; v != nil {
		result.SapFqdn = *v
	}

	return []OsSapConfiguration{
		result,
	}
}

func flattenDeployerVMPackages(input *sapvirtualinstances.DeployerVMPackages) []DeployerVmPackages {
	if input == nil {
		return nil
	}

	result := DeployerVmPackages{}

	if v := input.StorageAccountId; v != nil {
		result.StorageAccountId = *v
	}

	if v := input.Url; v != nil {
		result.Url = *v
	}

	return []DeployerVmPackages{
		result,
	}
}
