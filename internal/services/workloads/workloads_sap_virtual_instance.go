package workloads

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SchemaForSAPVirtualInstanceVirtualMachineConfiguration() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"image": {
					Type:     pluginsdk.TypeList,
					Required: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"offer": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"publisher": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									"RedHat",
									"SUSE",
								}, false),
							},

							"sku": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"version": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"os_profile": {
					Type:     pluginsdk.TypeList,
					Required: true,
					ForceNew: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"admin_username": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validate.AdminUsername,
							},

							"ssh_private_key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								Sensitive:    true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"ssh_public_key": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"vm_size": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaForSAPVirtualInstanceDiskVolumeConfiguration() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"volume_name": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						"backup",
						"hana/data",
						"hana/log",
						"hana/shared",
						"os",
						"usr/sap",
					}, false),
				},

				"count": {
					Type:     pluginsdk.TypeInt,
					Required: true,
					ForceNew: true,
				},

				"size_in_gb": {
					Type:     pluginsdk.TypeInt,
					Required: true,
					ForceNew: true,
				},

				"sku_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringInSlice(sapvirtualinstances.PossibleValuesForDiskSkuName(), false),
				},
			},
		},
	}
}

func SchemaForSAPVirtualInstanceVirtualMachineFullResourceNames() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"data_disk_names": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validate.DiskName,
					},
				},

				"host_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validate.HostName,
				},

				"network_interface_names": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: networkValidate.NetworkInterfaceName,
					},
				},

				"os_disk_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validate.DiskName,
				},

				"vm_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: computeValidate.VirtualMachineName,
				},
			},
		},
	}
}

func SchemaForSAPVirtualInstanceLoadBalancerFullResourceNames() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"backend_pool_names": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"frontend_ip_configuration_names": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"health_probe_names": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func expandVirtualMachineConfiguration(input []VirtualMachineConfiguration) sapvirtualinstances.VirtualMachineConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.VirtualMachineConfiguration{}
	}

	virtualMachineConfiguration := &input[0]

	result := sapvirtualinstances.VirtualMachineConfiguration{
		ImageReference: expandImageReference(virtualMachineConfiguration.ImageReference),
		OsProfile:      expandOsProfile(virtualMachineConfiguration.OSProfile),
		VMSize:         virtualMachineConfiguration.VmSize,
	}

	return result
}

func expandImageReference(input []ImageReference) sapvirtualinstances.ImageReference {
	if len(input) == 0 {
		return sapvirtualinstances.ImageReference{}
	}

	imageReference := &input[0]

	result := sapvirtualinstances.ImageReference{
		Offer:     utils.String(imageReference.Offer),
		Publisher: utils.String(imageReference.Publisher),
		Sku:       utils.String(imageReference.Sku),
		Version:   utils.String(imageReference.Version),
	}

	return result
}

func expandOsProfile(input []OSProfile) sapvirtualinstances.OSProfile {
	if len(input) == 0 {
		return sapvirtualinstances.OSProfile{}
	}

	osProfile := &input[0]

	result := sapvirtualinstances.OSProfile{
		AdminUsername: utils.String(osProfile.AdminUsername),
		OsConfiguration: &sapvirtualinstances.LinuxConfiguration{
			DisablePasswordAuthentication: utils.Bool(true),
			SshKeyPair: &sapvirtualinstances.SshKeyPair{
				PrivateKey: utils.String(osProfile.SshPrivateKey),
				PublicKey:  utils.String(osProfile.SshPublicKey),
			},
		},
	}

	return result
}

func expandVirtualMachineFullResourceNames(input []VirtualMachineFullResourceNames) (sapvirtualinstances.SingleServerFullResourceNames, error) {
	if len(input) == 0 {
		return sapvirtualinstances.SingleServerFullResourceNames{}, nil
	}

	if len(input) > 1 {
		return sapvirtualinstances.SingleServerFullResourceNames{}, fmt.Errorf("`virtual_machine_full_resource_names` only supports 1 item")
	}

	virtualMachineFullResourceNames := &input[0]

	result := sapvirtualinstances.SingleServerFullResourceNames{
		VirtualMachine: &sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandDataDiskNames(virtualMachineFullResourceNames.DataDiskNames),
			NetworkInterfaces: expandNetworkInterfaceNames(virtualMachineFullResourceNames.NetworkInterfaceNames),
		},
	}

	if v := virtualMachineFullResourceNames.HostName; v != "" {
		result.VirtualMachine.HostName = utils.String(v)
	}

	if v := virtualMachineFullResourceNames.OSDiskName; v != "" {
		result.VirtualMachine.OsDiskName = utils.String(v)
	}

	if v := virtualMachineFullResourceNames.VMName; v != "" {
		result.VirtualMachine.VirtualMachineName = utils.String(v)
	}

	return result, nil
}

func expandNetworkInterfaceNames(input []string) *[]sapvirtualinstances.NetworkInterfaceResourceNames {
	if len(input) == 0 {
		return nil
	}

	result := make([]sapvirtualinstances.NetworkInterfaceResourceNames, 0)

	for _, v := range input {
		networkInterfaceName := sapvirtualinstances.NetworkInterfaceResourceNames{
			NetworkInterfaceName: utils.String(v),
		}

		result = append(result, networkInterfaceName)
	}

	return &result
}

func expandDataDiskNames(input map[string]interface{}) *map[string][]string {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string][]string)

	for k, v := range input {
		result[k] = strings.Split(v.(string), ",")
	}

	return &result
}

func expandDiskVolumeConfigurations(input []DiskVolumeConfiguration) *sapvirtualinstances.DiskConfiguration {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]sapvirtualinstances.DiskVolumeConfiguration, 0)

	for _, v := range input {
		skuName := sapvirtualinstances.DiskSkuName(v.SkuName)

		result[v.VolumeName] = sapvirtualinstances.DiskVolumeConfiguration{
			Count:  utils.Int64(v.Count),
			SizeGB: utils.Int64(v.SizeGb),
			Sku: &sapvirtualinstances.DiskSku{
				Name: &skuName,
			},
		}
	}

	return &sapvirtualinstances.DiskConfiguration{
		DiskVolumeConfigurations: &result,
	}
}

func flattenDiskVolumeConfigurations(input *sapvirtualinstances.DiskConfiguration) []DiskVolumeConfiguration {
	result := make([]DiskVolumeConfiguration, 0)
	if input == nil || input.DiskVolumeConfigurations == nil {
		return result
	}

	for k, v := range *input.DiskVolumeConfigurations {
		diskVolumeConfiguration := DiskVolumeConfiguration{
			Count:      *v.Count,
			SizeGb:     *v.SizeGB,
			VolumeName: k,
		}

		if sku := v.Sku; sku != nil && sku.Name != nil {
			diskVolumeConfiguration.SkuName = string(*sku.Name)
		}

		result = append(result, diskVolumeConfiguration)
	}

	return result
}

func flattenVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []VirtualMachineFullResourceNames {
	result := make([]VirtualMachineFullResourceNames, 0)
	vmFullResourceNames := VirtualMachineFullResourceNames{}

	if vm := input.VirtualMachine; vm != nil {
		vmFullResourceNames.HostName = pointer.From(vm.HostName)
		vmFullResourceNames.OSDiskName = pointer.From(vm.OsDiskName)
		vmFullResourceNames.VMName = pointer.From(vm.VirtualMachineName)
		vmFullResourceNames.NetworkInterfaceNames = flattenNetworkInterfaceResourceNames(vm.NetworkInterfaces)
		vmFullResourceNames.DataDiskNames = flattenDataDiskNames(vm.DataDiskNames)
	}

	return append(result, vmFullResourceNames)
}

func flattenNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, *v.NetworkInterfaceName)
	}

	return result
}

func flattenDataDiskNames(input *map[string][]string) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	for k, v := range *input {
		results[k] = strings.Join(v, ",")
	}

	return results
}

func flattenVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData, basePath string) []VirtualMachineConfiguration {
	result := make([]VirtualMachineConfiguration, 0)

	return append(result, VirtualMachineConfiguration{
		ImageReference: flattenImageReference(input.ImageReference),
		OSProfile:      flattenOSProfile(input.OsProfile, d, fmt.Sprintf("%s.0.virtual_machine_configuration", basePath)),
		VmSize:         input.VMSize,
	})
}

func flattenImageReference(input sapvirtualinstances.ImageReference) []ImageReference {
	result := make([]ImageReference, 0)

	imageReference := ImageReference{
		Offer:     pointer.From(input.Offer),
		Publisher: pointer.From(input.Publisher),
		Sku:       pointer.From(input.Sku),
		Version:   pointer.From(input.Version),
	}

	return append(result, imageReference)
}

func flattenOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData, basePath string) []OSProfile {
	result := make([]OSProfile, 0)

	osProfile := OSProfile{
		AdminUsername: *input.AdminUsername,
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			if sshKeyPair := v.SshKeyPair; sshKeyPair != nil {
				osProfile.SshPrivateKey = d.Get(fmt.Sprintf("%s.0.os_profile.0.ssh_private_key", basePath)).(string)
				osProfile.SshPublicKey = pointer.From(sshKeyPair.PublicKey)
			}
		}
	}

	return append(result, osProfile)
}

func expandApplicationServer(input []ApplicationServerConfiguration) sapvirtualinstances.ApplicationServerConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.ApplicationServerConfiguration{}
	}

	applicationServer := &input[0]

	result := sapvirtualinstances.ApplicationServerConfiguration{
		InstanceCount:               applicationServer.InstanceCount,
		SubnetId:                    applicationServer.SubnetId,
		VirtualMachineConfiguration: expandVirtualMachineConfiguration(applicationServer.VirtualMachineConfiguration),
	}

	return result
}

func expandCentralServer(input []CentralServerConfiguration) sapvirtualinstances.CentralServerConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.CentralServerConfiguration{}
	}

	centralServer := &input[0]

	result := sapvirtualinstances.CentralServerConfiguration{
		InstanceCount:               centralServer.InstanceCount,
		SubnetId:                    centralServer.SubnetId,
		VirtualMachineConfiguration: expandVirtualMachineConfiguration(centralServer.VirtualMachineConfiguration),
	}

	return result
}

func expandDatabaseServer(input []DatabaseServerConfiguration) sapvirtualinstances.DatabaseConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.DatabaseConfiguration{}
	}

	databaseServer := &input[0]

	result := sapvirtualinstances.DatabaseConfiguration{
		DiskConfiguration:           expandDiskVolumeConfigurations(databaseServer.DiskVolumeConfigurations),
		InstanceCount:               databaseServer.InstanceCount,
		SubnetId:                    databaseServer.SubnetId,
		VirtualMachineConfiguration: expandVirtualMachineConfiguration(databaseServer.VirtualMachineConfiguration),
	}

	if v := databaseServer.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	return result
}

func expandStorageConfiguration(input *ThreeTierConfiguration) (*sapvirtualinstances.StorageConfiguration, error) {
	if len(input.TransportCreateAndMount) == 0 && len(input.TransportMount) == 0 {
		return &sapvirtualinstances.StorageConfiguration{
			TransportFileShareConfiguration: sapvirtualinstances.SkipFileShareConfiguration{},
		}, nil
	}

	result := sapvirtualinstances.StorageConfiguration{}

	if len(input.TransportCreateAndMount) != 0 {
		transportCreateAndMount, err := expandTransportCreateAndMount(input.TransportCreateAndMount)
		if err != nil {
			return nil, err
		}
		result.TransportFileShareConfiguration = transportCreateAndMount
	}

	if len(input.TransportMount) != 0 {
		result.TransportFileShareConfiguration = expandTransportMount(input.TransportMount)
	}

	return &result, nil
}

func expandTransportCreateAndMount(input []TransportCreateAndMount) (sapvirtualinstances.CreateAndMountFileShareConfiguration, error) {
	if len(input) == 0 {
		return sapvirtualinstances.CreateAndMountFileShareConfiguration{}, nil
	}

	transportCreateAndMount := &input[0]

	result := sapvirtualinstances.CreateAndMountFileShareConfiguration{}

	if v := transportCreateAndMount.ResourceGroupId; v != "" {
		resourceGroupId, err := commonids.ParseResourceGroupID(v)
		if err != nil {
			return sapvirtualinstances.CreateAndMountFileShareConfiguration{}, err
		}
		result.ResourceGroup = utils.String(resourceGroupId.ResourceGroupName)
	}

	if v := transportCreateAndMount.StorageAccountName; v != "" {
		result.StorageAccountName = utils.String(v)
	}

	return result, nil
}

func expandTransportMount(input []TransportMount) sapvirtualinstances.MountFileShareConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.MountFileShareConfiguration{}
	}

	transportMount := &input[0]

	// Currently, the last segment of the Storage File Share resource manager ID in Swagger is defined as `/shares/` but it's unexpected.
	// The last segment of the Storage File Share resource manager ID should be `/fileshares/` not `/shares/` in Swagger since the backend service is using `/fileshares/`.
	// See more details from https://github.com/Azure/azure-rest-api-specs/issues/25209
	result := sapvirtualinstances.MountFileShareConfiguration{
		Id:                strings.Replace(transportMount.FileShareId, "/fileshares/", "/shares/", 1),
		PrivateEndpointId: transportMount.PrivateEndpointId,
	}

	return result
}

func expandFullResourceNames(input []FullResourceNames) sapvirtualinstances.ThreeTierFullResourceNames {
	if len(input) == 0 {
		return sapvirtualinstances.ThreeTierFullResourceNames{}
	}

	fullResourceNames := &input[0]

	result := sapvirtualinstances.ThreeTierFullResourceNames{
		ApplicationServer: expandApplicationServerFullResourceNames(fullResourceNames.ApplicationServer),
		CentralServer:     expandCentralServerFullResourceNames(fullResourceNames.CentralServer),
		DatabaseServer:    expandDatabaseServerFullResourceNames(fullResourceNames.DatabaseServer),
		SharedStorage:     expandSharedStorage(fullResourceNames.SharedStorage),
	}

	return result
}

func expandApplicationServerFullResourceNames(input []ApplicationServerFullResourceNames) *sapvirtualinstances.ApplicationServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	applicationServerFullResourceNames := &input[0]

	result := sapvirtualinstances.ApplicationServerFullResourceNames{
		VirtualMachines: expandVirtualMachinesFullResourceNames(applicationServerFullResourceNames.VirtualMachines),
	}

	if v := applicationServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return &result
}

func expandVirtualMachinesFullResourceNames(input []VirtualMachineFullResourceNames) *[]sapvirtualinstances.VirtualMachineResourceNames {
	if len(input) == 0 {
		return nil
	}

	result := make([]sapvirtualinstances.VirtualMachineResourceNames, 0)

	for _, item := range input {
		vmResourceNames := sapvirtualinstances.VirtualMachineResourceNames{
			DataDiskNames:     expandDataDiskNames(item.DataDiskNames),
			NetworkInterfaces: expandNetworkInterfaceNames(item.NetworkInterfaceNames),
		}

		if v := item.HostName; v != "" {
			vmResourceNames.HostName = utils.String(v)
		}

		if v := item.OSDiskName; v != "" {
			vmResourceNames.OsDiskName = utils.String(v)
		}

		if v := item.VMName; v != "" {
			vmResourceNames.VirtualMachineName = utils.String(v)
		}

		result = append(result, vmResourceNames)
	}

	return &result
}

func expandCentralServerFullResourceNames(input []CentralServerFullResourceNames) *sapvirtualinstances.CentralServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	centralServerFullResourceNames := &input[0]

	result := sapvirtualinstances.CentralServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerFullResourceNames(centralServerFullResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesFullResourceNames(centralServerFullResourceNames.VirtualMachines),
	}

	if v := centralServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return &result
}

func expandLoadBalancerFullResourceNames(input []LoadBalancer) *sapvirtualinstances.LoadBalancerResourceNames {
	if len(input) == 0 {
		return nil
	}

	loadBalancerFullResourceNames := &input[0]

	result := sapvirtualinstances.LoadBalancerResourceNames{}

	if v := loadBalancerFullResourceNames.Name; v != "" {
		result.LoadBalancerName = utils.String(v)
	}

	if v := loadBalancerFullResourceNames.BackendPoolNames; v != nil {
		result.BackendPoolNames = &v
	}

	if v := loadBalancerFullResourceNames.FrontendIpConfigurationNames; v != nil {
		result.FrontendIPConfigurationNames = &v
	}

	if v := loadBalancerFullResourceNames.HealthProbeNames; v != nil {
		result.HealthProbeNames = &v
	}

	return &result
}

func expandDatabaseServerFullResourceNames(input []DatabaseServerFullResourceNames) *sapvirtualinstances.DatabaseServerFullResourceNames {
	if len(input) == 0 {
		return nil
	}

	databaseServerFullResourceNames := &input[0]

	result := sapvirtualinstances.DatabaseServerFullResourceNames{
		LoadBalancer:    expandLoadBalancerFullResourceNames(databaseServerFullResourceNames.LoadBalancer),
		VirtualMachines: expandVirtualMachinesFullResourceNames(databaseServerFullResourceNames.VirtualMachines),
	}

	if v := databaseServerFullResourceNames.AvailabilitySetName; v != "" {
		result.AvailabilitySetName = utils.String(v)
	}

	return &result
}

func expandSharedStorage(input []SharedStorage) *sapvirtualinstances.SharedStorageResourceNames {
	if len(input) == 0 {
		return nil
	}

	sharedStorage := &input[0]

	result := sapvirtualinstances.SharedStorageResourceNames{}

	if v := sharedStorage.AccountName; v != "" {
		result.SharedStorageAccountName = utils.String(v)
	}

	if v := sharedStorage.PrivateEndpointName; v != "" {
		result.SharedStorageAccountPrivateEndPointName = utils.String(v)
	}

	return &result
}

func flattenApplicationServer(input sapvirtualinstances.ApplicationServerConfiguration, d *pluginsdk.ResourceData, basePath string) []ApplicationServerConfiguration {
	result := make([]ApplicationServerConfiguration, 0)

	applicationServerConfig := ApplicationServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.application_server_configuration", basePath)),
	}

	return append(result, applicationServerConfig)
}

func flattenCentralServer(input sapvirtualinstances.CentralServerConfiguration, d *pluginsdk.ResourceData, basePath string) []CentralServerConfiguration {
	result := make([]CentralServerConfiguration, 0)

	centralServerConfig := CentralServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.central_server_configuration", basePath)),
	}

	return append(result, centralServerConfig)
}

func flattenDatabaseServer(input sapvirtualinstances.DatabaseConfiguration, d *pluginsdk.ResourceData, basePath string) []DatabaseServerConfiguration {
	result := make([]DatabaseServerConfiguration, 0)

	databaseServerConfig := DatabaseServerConfiguration{
		DiskVolumeConfigurations:    flattenDiskVolumeConfigurations(input.DiskConfiguration),
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.database_server_configuration", basePath)),
	}

	if v := input.DatabaseType; v != nil {
		databaseServerConfig.DatabaseType = string(*v)
	}

	return append(result, databaseServerConfig)
}

func flattenFullResourceNames(input sapvirtualinstances.ThreeTierFullResourceNames) []FullResourceNames {
	result := make([]FullResourceNames, 0)

	fullResourceNames := FullResourceNames{
		ApplicationServer: flattenApplicationServerFullResourceNames(input.ApplicationServer),
		CentralServer:     flattenCentralServerFullResourceNames(input.CentralServer),
		DatabaseServer:    flattenDatabaseServerFullResourceNames(input.DatabaseServer),
		SharedStorage:     flattenSharedStorage(input.SharedStorage),
	}

	return append(result, fullResourceNames)
}

func flattenApplicationServerFullResourceNames(input *sapvirtualinstances.ApplicationServerFullResourceNames) []ApplicationServerFullResourceNames {
	result := make([]ApplicationServerFullResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, ApplicationServerFullResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		VirtualMachines:     flattenVirtualMachinesFullResourceNames(input.VirtualMachines),
	})
}

func flattenVirtualMachinesFullResourceNames(input *[]sapvirtualinstances.VirtualMachineResourceNames) []VirtualMachineFullResourceNames {
	result := make([]VirtualMachineFullResourceNames, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, VirtualMachineFullResourceNames{
			HostName:              pointer.From(item.HostName),
			OSDiskName:            pointer.From(item.OsDiskName),
			VMName:                pointer.From(item.VirtualMachineName),
			DataDiskNames:         flattenDataDiskNames(item.DataDiskNames),
			NetworkInterfaceNames: flattenNetworkInterfaceResourceNames(item.NetworkInterfaces),
		})
	}

	return result
}

func flattenCentralServerFullResourceNames(input *sapvirtualinstances.CentralServerFullResourceNames) []CentralServerFullResourceNames {
	result := make([]CentralServerFullResourceNames, 0)
	if input == nil {
		return result
	}

	centralServerFullResourceNames := CentralServerFullResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerFullResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesFullResourceNames(input.VirtualMachines),
	}

	return append(result, centralServerFullResourceNames)
}

func flattenLoadBalancerFullResourceNames(input *sapvirtualinstances.LoadBalancerResourceNames) []LoadBalancer {
	result := make([]LoadBalancer, 0)
	if input == nil {
		return result
	}

	return append(result, LoadBalancer{
		Name:                         pointer.From(input.LoadBalancerName),
		BackendPoolNames:             pointer.From(input.BackendPoolNames),
		FrontendIpConfigurationNames: pointer.From(input.FrontendIPConfigurationNames),
		HealthProbeNames:             pointer.From(input.HealthProbeNames),
	})
}

func flattenDatabaseServerFullResourceNames(input *sapvirtualinstances.DatabaseServerFullResourceNames) []DatabaseServerFullResourceNames {
	result := make([]DatabaseServerFullResourceNames, 0)
	if input == nil {
		return result
	}

	return append(result, DatabaseServerFullResourceNames{
		AvailabilitySetName: pointer.From(input.AvailabilitySetName),
		LoadBalancer:        flattenLoadBalancerFullResourceNames(input.LoadBalancer),
		VirtualMachines:     flattenVirtualMachinesFullResourceNames(input.VirtualMachines),
	})
}

func flattenSharedStorage(input *sapvirtualinstances.SharedStorageResourceNames) []SharedStorage {
	result := make([]SharedStorage, 0)
	if input == nil {
		return result
	}

	return append(result, SharedStorage{
		AccountName:         pointer.From(input.SharedStorageAccountName),
		PrivateEndpointName: pointer.From(input.SharedStorageAccountPrivateEndPointName),
	})
}
