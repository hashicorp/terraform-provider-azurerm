package workloads

import (
	"fmt"
	"strings"

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
				"image_reference": {
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

							"ssh_key_pair": {
								Type:     pluginsdk.TypeList,
								Required: true,
								ForceNew: true,
								MaxItems: 1,
								Elem: &pluginsdk.Resource{
									Schema: map[string]*pluginsdk.Schema{
										"private_key": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ForceNew:     true,
											Sensitive:    true,
											ValidateFunc: validation.StringIsNotEmpty,
										},

										"public_key": {
											Type:         pluginsdk.TypeString,
											Required:     true,
											ForceNew:     true,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},
								},
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

				"size_gb": {
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
		AdminUsername:   utils.String(osProfile.AdminUsername),
		OsConfiguration: expandLinuxConfiguration(osProfile.SshKeyPair),
	}

	return result
}

func expandLinuxConfiguration(input []SshKeyPair) *sapvirtualinstances.LinuxConfiguration {
	if len(input) == 0 {
		return nil
	}

	sshKeyPair := &input[0]

	return &sapvirtualinstances.LinuxConfiguration{
		DisablePasswordAuthentication: utils.Bool(true),
		SshKeyPair: &sapvirtualinstances.SshKeyPair{
			PrivateKey: utils.String(sshKeyPair.PrivateKey),
			PublicKey:  utils.String(sshKeyPair.PublicKey),
		},
	}
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
		VirtualMachine: &sapvirtualinstances.VirtualMachineResourceNames{},
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

	if v := virtualMachineFullResourceNames.NetworkInterfaceNames; v != nil {
		result.VirtualMachine.NetworkInterfaces = expandNetworkInterfaceNames(v)
	}

	if v := virtualMachineFullResourceNames.DataDiskNames; v != nil {
		result.VirtualMachine.DataDiskNames = expandDataDiskNames(v)
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

func flattenDiskVolumeConfigurations(input *map[string]sapvirtualinstances.DiskVolumeConfiguration) []DiskVolumeConfiguration {
	if input == nil {
		return nil
	}

	result := make([]DiskVolumeConfiguration, 0)

	for k, v := range *input {
		diskVolumeConfiguration := DiskVolumeConfiguration{
			Count:      *v.Count,
			SizeGb:     *v.SizeGB,
			SkuName:    string(*v.Sku.Name),
			VolumeName: k,
		}

		result = append(result, diskVolumeConfiguration)
	}

	return result
}

func flattenVirtualMachineFullResourceNames(input sapvirtualinstances.SingleServerFullResourceNames) []VirtualMachineFullResourceNames {
	result := VirtualMachineFullResourceNames{}

	if vm := input.VirtualMachine; vm != nil {
		if v := vm.HostName; v != nil {
			result.HostName = *v
		}

		if v := vm.OsDiskName; v != nil {
			result.OSDiskName = *v
		}

		if v := vm.VirtualMachineName; v != nil {
			result.VMName = *v
		}

		if v := vm.NetworkInterfaces; v != nil {
			result.NetworkInterfaceNames = flattenNetworkInterfaceResourceNames(v)
		}

		if v := vm.DataDiskNames; v != nil {
			result.DataDiskNames = flattenDataDiskNames(v)
		}
	}

	return []VirtualMachineFullResourceNames{
		result,
	}
}

func flattenNetworkInterfaceResourceNames(input *[]sapvirtualinstances.NetworkInterfaceResourceNames) []string {
	if input == nil {
		return nil
	}

	result := make([]string, 0)

	for _, v := range *input {
		result = append(result, *v.NetworkInterfaceName)
	}

	return result
}

func flattenDataDiskNames(input *map[string][]string) map[string]interface{} {
	if input == nil {
		return nil
	}

	results := make(map[string]interface{})

	for k, v := range *input {
		results[k] = strings.Join(v, ",")
	}

	return results
}

func flattenVirtualMachineConfiguration(input sapvirtualinstances.VirtualMachineConfiguration, d *pluginsdk.ResourceData, basePath string) []VirtualMachineConfiguration {
	result := VirtualMachineConfiguration{
		ImageReference: flattenImageReference(input.ImageReference),
		OSProfile:      flattenOSProfile(input.OsProfile, d, fmt.Sprintf("%s.0.virtual_machine_configuration", basePath)),
		VmSize:         input.VMSize,
	}

	return []VirtualMachineConfiguration{
		result,
	}
}

func flattenImageReference(input sapvirtualinstances.ImageReference) []ImageReference {
	result := ImageReference{
		Offer:     *input.Offer,
		Publisher: *input.Publisher,
		Sku:       *input.Sku,
		Version:   *input.Version,
	}

	return []ImageReference{
		result,
	}
}

func flattenOSProfile(input sapvirtualinstances.OSProfile, d *pluginsdk.ResourceData, basePath string) []OSProfile {
	result := OSProfile{
		AdminUsername: *input.AdminUsername,
	}

	if osConfiguration := input.OsConfiguration; osConfiguration != nil {
		if v, ok := osConfiguration.(sapvirtualinstances.LinuxConfiguration); ok {
			result.SshKeyPair = flattenSshKeyPair(v.SshKeyPair, d, fmt.Sprintf("%s.0.os_profile", basePath))
		}
	}

	return []OSProfile{
		result,
	}
}

func flattenSshKeyPair(input *sapvirtualinstances.SshKeyPair, d *pluginsdk.ResourceData, basePath string) []SshKeyPair {
	if input == nil {
		return nil
	}

	privateKeyPath := fmt.Sprintf("%s.0.ssh_key_pair.0.private_key", basePath)
	result := SshKeyPair{
		PrivateKey: d.Get(privateKeyPath).(string),
	}

	if v := input.PublicKey; v != nil {
		result.PublicKey = *v
	}

	return []SshKeyPair{
		result,
	}
}

func expandApplicationServer(input []ApplicationServerConfiguration) sapvirtualinstances.ApplicationServerConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.ApplicationServerConfiguration{}
	}

	applicationServer := &input[0]

	result := sapvirtualinstances.ApplicationServerConfiguration{
		InstanceCount: applicationServer.InstanceCount,
		SubnetId:      applicationServer.SubnetId,
	}

	if v := applicationServer.VirtualMachineConfiguration; v != nil {
		result.VirtualMachineConfiguration = expandVirtualMachineConfiguration(v)
	}

	return result
}

func expandCentralServer(input []CentralServerConfiguration) sapvirtualinstances.CentralServerConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.CentralServerConfiguration{}
	}

	centralServer := &input[0]

	result := sapvirtualinstances.CentralServerConfiguration{
		InstanceCount: centralServer.InstanceCount,
		SubnetId:      centralServer.SubnetId,
	}

	if v := centralServer.VirtualMachineConfiguration; v != nil {
		result.VirtualMachineConfiguration = expandVirtualMachineConfiguration(v)
	}

	return result
}

func expandDatabaseServer(input []DatabaseServerConfiguration) sapvirtualinstances.DatabaseConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.DatabaseConfiguration{}
	}

	databaseServer := &input[0]

	result := sapvirtualinstances.DatabaseConfiguration{
		InstanceCount:               databaseServer.InstanceCount,
		SubnetId:                    databaseServer.SubnetId,
		VirtualMachineConfiguration: expandVirtualMachineConfiguration(databaseServer.VirtualMachineConfiguration),
	}

	if v := databaseServer.DatabaseType; v != "" {
		dbType := sapvirtualinstances.SAPDatabaseType(v)
		result.DatabaseType = &dbType
	}

	if v := databaseServer.DiskVolumeConfigurations; v != nil {
		result.DiskConfiguration = expandDiskVolumeConfigurations(v)
	}

	return result
}

func expandStorageConfiguration(input *ThreeTierConfiguration) *sapvirtualinstances.StorageConfiguration {
	if len(input.TransportCreateAndMount) == 0 && len(input.TransportMount) == 0 {
		return &sapvirtualinstances.StorageConfiguration{
			sapvirtualinstances.SkipFileShareConfiguration{},
		}
	}

	result := sapvirtualinstances.StorageConfiguration{}

	if len(input.TransportCreateAndMount) != 0 {
		result.TransportFileShareConfiguration = expandTransportCreateAndMount(input.TransportCreateAndMount)
	}

	if len(input.TransportMount) != 0 {
		result.TransportFileShareConfiguration = expandTransportMount(input.TransportMount)
	}

	return &result
}

func expandTransportCreateAndMount(input []TransportCreateAndMount) sapvirtualinstances.CreateAndMountFileShareConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.CreateAndMountFileShareConfiguration{}
	}

	transportCreateAndMount := &input[0]

	result := sapvirtualinstances.CreateAndMountFileShareConfiguration{}

	if v := transportCreateAndMount.ResourceGroupName; v != "" {
		result.ResourceGroup = utils.String(v)
	}

	if v := transportCreateAndMount.StorageAccountName; v != "" {
		result.StorageAccountName = utils.String(v)
	}

	return result
}

func expandTransportMount(input []TransportMount) sapvirtualinstances.MountFileShareConfiguration {
	if len(input) == 0 {
		return sapvirtualinstances.MountFileShareConfiguration{}
	}

	transportMount := &input[0]

	result := sapvirtualinstances.MountFileShareConfiguration{
		Id:                transportMount.ShareFileId,
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
		vmResourceNames := sapvirtualinstances.VirtualMachineResourceNames{}

		if v := item.HostName; v != "" {
			vmResourceNames.HostName = utils.String(v)
		}

		if v := item.OSDiskName; v != "" {
			vmResourceNames.OsDiskName = utils.String(v)
		}

		if v := item.VMName; v != "" {
			vmResourceNames.VirtualMachineName = utils.String(v)
		}

		if v := item.NetworkInterfaceNames; v != nil {
			vmResourceNames.NetworkInterfaces = expandNetworkInterfaceNames(v)
		}

		if v := item.DataDiskNames; v != nil {
			vmResourceNames.DataDiskNames = expandDataDiskNames(v)
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
	result := ApplicationServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.application_server_configuration", basePath)),
	}

	return []ApplicationServerConfiguration{
		result,
	}
}

func flattenCentralServer(input sapvirtualinstances.CentralServerConfiguration, d *pluginsdk.ResourceData, basePath string) []CentralServerConfiguration {
	result := CentralServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.central_server_configuration", basePath)),
	}

	return []CentralServerConfiguration{
		result,
	}
}

func flattenDatabaseServer(input sapvirtualinstances.DatabaseConfiguration, d *pluginsdk.ResourceData, basePath string) []DatabaseServerConfiguration {
	result := DatabaseServerConfiguration{
		InstanceCount:               input.InstanceCount,
		SubnetId:                    input.SubnetId,
		VirtualMachineConfiguration: flattenVirtualMachineConfiguration(input.VirtualMachineConfiguration, d, fmt.Sprintf("%s.0.database_server_configuration", basePath)),
	}

	if v := input.DatabaseType; v != nil {
		result.DatabaseType = string(*v)
	}

	if v := input.DiskConfiguration; v != nil && v.DiskVolumeConfigurations != nil {
		result.DiskVolumeConfigurations = flattenDiskVolumeConfigurations(v.DiskVolumeConfigurations)
	}

	return []DatabaseServerConfiguration{
		result,
	}
}

func flattenFullResourceNames(input sapvirtualinstances.ThreeTierFullResourceNames) []FullResourceNames {
	result := FullResourceNames{}

	if v := input.ApplicationServer; v != nil {
		result.ApplicationServer = flattenApplicationServerFullResourceNames(v)
	}

	if v := input.CentralServer; v != nil {
		result.CentralServer = flattenCentralServerFullResourceNames(v)
	}

	if v := input.DatabaseServer; v != nil {
		result.DatabaseServer = flattenDatabaseServerFullResourceNames(v)
	}

	if v := input.SharedStorage; v != nil {
		result.SharedStorage = flattenSharedStorage(v)
	}

	return []FullResourceNames{
		result,
	}
}

func flattenApplicationServerFullResourceNames(input *sapvirtualinstances.ApplicationServerFullResourceNames) []ApplicationServerFullResourceNames {
	if input == nil {
		return nil
	}

	result := ApplicationServerFullResourceNames{}

	if v := input.AvailabilitySetName; v != nil {
		result.AvailabilitySetName = *v
	}

	if v := input.VirtualMachines; v != nil {
		result.VirtualMachines = flattenVirtualMachinesFullResourceNames(v)
	}

	return []ApplicationServerFullResourceNames{
		result,
	}
}

func flattenVirtualMachinesFullResourceNames(input *[]sapvirtualinstances.VirtualMachineResourceNames) []VirtualMachineFullResourceNames {
	if input == nil {
		return nil
	}

	result := make([]VirtualMachineFullResourceNames, 0)

	for _, item := range *input {
		virtualMachineFullResourceNames := VirtualMachineFullResourceNames{}

		if v := item.HostName; v != nil {
			virtualMachineFullResourceNames.HostName = *v
		}

		if v := item.OsDiskName; v != nil {
			virtualMachineFullResourceNames.OSDiskName = *v
		}

		if v := item.VirtualMachineName; v != nil {
			virtualMachineFullResourceNames.VMName = *v
		}

		if v := item.NetworkInterfaces; v != nil {
			virtualMachineFullResourceNames.NetworkInterfaceNames = flattenNetworkInterfaceResourceNames(v)
		}

		if v := item.DataDiskNames; v != nil {
			virtualMachineFullResourceNames.DataDiskNames = flattenDataDiskNames(v)
		}

		result = append(result, virtualMachineFullResourceNames)
	}

	return result
}

func flattenCentralServerFullResourceNames(input *sapvirtualinstances.CentralServerFullResourceNames) []CentralServerFullResourceNames {
	if input == nil {
		return nil
	}

	result := CentralServerFullResourceNames{}

	if v := input.AvailabilitySetName; v != nil {
		result.AvailabilitySetName = *v
	}

	if v := input.LoadBalancer; v != nil {
		result.LoadBalancer = flattenLoadBalancerFullResourceNames(v)
	}

	if v := input.VirtualMachines; v != nil {
		result.VirtualMachines = flattenVirtualMachinesFullResourceNames(v)
	}

	return []CentralServerFullResourceNames{
		result,
	}
}

func flattenLoadBalancerFullResourceNames(input *sapvirtualinstances.LoadBalancerResourceNames) []LoadBalancer {
	if input == nil {
		return nil
	}

	result := LoadBalancer{}

	if v := input.LoadBalancerName; v != nil {
		result.Name = *v
	}

	if v := input.BackendPoolNames; v != nil {
		result.BackendPoolNames = *v
	}

	if v := input.FrontendIPConfigurationNames; v != nil {
		result.FrontendIpConfigurationNames = *v
	}

	if v := input.HealthProbeNames; v != nil {
		result.HealthProbeNames = *v
	}

	return []LoadBalancer{
		result,
	}
}

func flattenDatabaseServerFullResourceNames(input *sapvirtualinstances.DatabaseServerFullResourceNames) []DatabaseServerFullResourceNames {
	if input == nil {
		return nil
	}

	result := DatabaseServerFullResourceNames{}

	if v := input.AvailabilitySetName; v != nil {
		result.AvailabilitySetName = *v
	}

	if v := input.LoadBalancer; v != nil {
		result.LoadBalancer = flattenLoadBalancerFullResourceNames(v)
	}

	if v := input.VirtualMachines; v != nil {
		result.VirtualMachines = flattenVirtualMachinesFullResourceNames(v)
	}

	return []DatabaseServerFullResourceNames{
		result,
	}
}

func flattenSharedStorage(input *sapvirtualinstances.SharedStorageResourceNames) []SharedStorage {
	if input == nil {
		return nil
	}

	result := SharedStorage{}

	if v := input.SharedStorageAccountName; v != nil {
		result.AccountName = *v
	}

	if v := input.SharedStorageAccountPrivateEndPointName; v != nil {
		result.PrivateEndpointName = *v
	}

	return []SharedStorage{
		result,
	}
}
