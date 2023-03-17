package compute

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2022-08-01/compute"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type OrchestratedVirtualMachineScaleSetDataSource struct{}

var _ sdk.DataSource = OrchestratedVirtualMachineScaleSetDataSource{}

type OrchestratedVirtualMachineScaleSetDataSourceModel struct {
	Name                    string                                        `tfschema:"name"`
	ResourceGroup           string                                        `tfschema:"resource_group_name"`
	Location                string                                        `tfschema:"location"`
	NetworkInterface        []VirtualMachineScaleSetNetworkInterface      `tfschema:"network_interface"`
	OSDisk                  VirtualMachineScaleSetOSDisk                  `tfschema:"os_disk"`
	SkuName                 string                                        `tfschema:"sku_name"`
	OSProfile               VirtualMachineScaleSetOSProfile               `tfschema:"os_profile"`
	AutomaticInstanceRepair VirtualMachineScaleSetAutomaticInstanceRepair `tfschema:"automatic_instance_repair"`
}

type VirtualMachineScaleSetNetworkInterface struct {
	Name                         string `tfschema:"name"`
	IPConfiguration              []VirtualMachineScaleSetNetworkInterfaceIPConfiguration
	DNSServers                   []string `tfschema:"dns_servers"`
	AcceleratedNetworkingEnabled bool     `tfschema:"accelerated_networking_enabled"`
	IPForwardingEnabled          bool     `tfschema:"ip_forwarding_enabled"`
	NetworkSecurityGroupId       string   `tfschema:"network_security_group_id"`
	Primary                      bool     `tfschema:"primary"`
}

type VirtualMachineScaleSetNetworkInterfaceIPConfiguration struct {
	Name                                    string                                                                 `tfschema:"name"`
	ApplicationGatewayBackendAddressPoolIds []string                                                               `tfschema:"application_gateway_backend_address_pool_ids"`
	ApplicationSecurityGroupIds             []string                                                               `tfschema:"application_security_group_ids"`
	LoadBalancerBackendAddressPoolIds       []string                                                               `tfschema:"load_balancer_backend_address_pool_ids"`
	Primary                                 bool                                                                   `tfschema:"primary"`
	PublicIPAddress                         []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress `tfschema:"public_ip_address"`
	SubnetId                                string                                                                 `tfschema:"subnet_id"`
	Version                                 string                                                                 `tfschema:"version"`
}

type VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress struct {
	Name                 string                                                                      `tfschema:"name"`
	DomainNameLabel      string                                                                      `tfschema:"domain_name_label"`
	IdleTimeoutInMinutes int                                                                         `tfschema:"idle_timeout_in_minutes"`
	IPTag                []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag `tfschema:"ip_tag"`
	PublicIpPrefixId     string                                                                      `tfschema:"public_ip_prefix_id"`
	SkuName              string                                                                      `tfschema:"sku_name"`
	Version              string                                                                      `tfschema:"version"`
}

type VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag struct {
	Tag  string `tfschema:"tag"`
	Type string `tfschema:"type"`
}

type VirtualMachineScaleSetOSDisk struct {
	Caching                 string                                       `tfschema:"caching"`
	StorageAccountType      string                                       `tfschema:"storage_account_type"`
	DiffDiskSettings        VirtualMachineScaleSetOSDiskDiffDiskSettings `tfschema:"diff_disk_settings"`
	DiskEncryptionSetId     string                                       `tfschema:"disk_encryption_set_id"`
	DiskSizeGB              string                                       `tfschema:"disk_size_gb"`
	WriteAcceleratorEnabled bool                                         `tfschema:"write_accelerator_enabled"`
}

type VirtualMachineScaleSetOSDiskDiffDiskSettings struct {
	Option    string `tfschema:"option"`
	Placement string `tfschema:"placement"`
}

type VirtualMachineScaleSetOSProfile struct {
	CustomData           string                                              `tfschema:"custom_data"`
	WindowsConfiguration VirtualMachineScaleSetOSProfileWindowsConfiguration `tfschema:"windows_configuration"`
	LinuxConfiguration   VirtualMachineScaleSetOSProfileLinuxConfiguration   `tfschema:"linux_configuration"`
}

type VirtualMachineScaleSetOSProfileWindowsConfiguration struct {
	AdminUsername           string                                                             `tfschema:"admin_username"`
	AdminPassword           string                                                             `tfschema:"admin_password"`
	ComputerNamePrefix      string                                                             `tfschema:"computer_name_prefix"`
	AutomaticUpdatesEnabled bool                                                               `tfschema:"automatic_updates_enabled"`
	HotpatchingEnabled      bool                                                               `tfschema:"hotpatching_enabled"`
	ProvisionVMAgent        bool                                                               `tfschema:"provision_vm_agent"`
	PatchAssessmentMode     string                                                             `tfschema:"patch_assessment_mode"`
	PatchMode               string                                                             `tfschema:"patch_mode"`
	Secret                  []VirtualMachineScaleSetOSProfileWindowsConfigurationSecret        `tfschema:"secret"`
	Timezone                string                                                             `tfschema:"timezone"`
	WinRMListener           []VirtualMachineScaleSetOSProfileWindowsConfigurationWinRMListener `tfschema:"winrm_listener"`
}

type VirtualMachineScaleSetOSProfileWindowsConfigurationSecret struct {
	KeyVaultId  string                                                                `tfschema:"key_vault_id"`
	Certificate []VirtualMachineScaleSetOSProfileWindowsConfigurationSecretCertficate `tfschema:"certificate"`
}

type VirtualMachineScaleSetOSProfileWindowsConfigurationSecretCertficate struct {
	Store string `tfschema:"store"`
	Url   string `tfschema:"url"`
}

type VirtualMachineScaleSetOSProfileWindowsConfigurationWinRMListener struct {
	Protocol       string `tfschema:"protocol"`
	CertificateUrl string `tfschema:"certificate_url"`
}

type VirtualMachineScaleSetOSProfileLinuxConfiguration struct {
	AdminUsername                 string                                                         `tfschema:"admin_username"`
	AdminPassword                 string                                                         `tfschema:"admin_password"`
	AdminSSHKey                   []VirtualMachineScaleSetOSProfileLinuxConfigurationAdminSSHKey `tfschema:"admin_ssh_key"`
	ComputerNamePrefix            string                                                         `tfschema:"computer_name_prefix"`
	DisablePasswordAuthentication bool                                                           `tfschema:"disable_password_authentication"`
	ProvisionVMAgent              bool                                                           `tfschema:"provision_vm_agent"`
	PatchAssessmentMode           string                                                         `tfschema:"patch_assessment_mode"`
	PatchMode                     string                                                         `tfschema:"patch_mode"`
	Secret                        []VirtualMachineScaleSetOSProfileLinuxConfigurationSecret      `tfschema:"secret"`
}

type VirtualMachineScaleSetOSProfileLinuxConfigurationAdminSSHKey struct {
	PublicKey string `tfschema:"public_key"`
	Username  string `tfschema:"username"`
}

type VirtualMachineScaleSetOSProfileLinuxConfigurationSecret struct {
	KeyVaultId  string                                                              `tfschema:"key_vault_id"`
	Certificate []VirtualMachineScaleSetOSProfileLinuxConfigurationSecretCertficate `tfschema:"certificate"`
}

type VirtualMachineScaleSetOSProfileLinuxConfigurationSecretCertficate struct {
	Url string `tfschema:"url"`
}

func (r OrchestratedVirtualMachineScaleSetDataSource) ModelObject() interface{} {
	return &OrchestratedVirtualMachineScaleSetDataSourceModel{}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) ResourceType() string {
	return "azurerm_orchestrated_virtual_machine_scale_set"
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: computeValidate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),
	}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VMScaleSetClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var orchestratedVMSS OrchestratedVirtualMachineScaleSetDataSourceModel
			if err := metadata.Decode(&orchestratedVMSS); err != nil {
				return err
			}

			id := parse.NewVirtualMachineScaleSetID(subscriptionId, orchestratedVMSS.ResourceGroup, orchestratedVMSS.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.Name, compute.ExpandTypesForGetVMScaleSetsUserData)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			orchestratedVMSS.Location = location.NormalizeNilable(existing.Location)

			metadata.SetID(id)

			return metadata.Encode(&orchestratedVMSS)
		},
	}
}
