package nginx

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01/nginxdeployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FrontendPrivate struct {
	IpAddress        string `tfschema:"ip_address"`
	AllocationMethod string `tfschema:"allocation_method"`
	SubnetId         string `tfschema:"subnet_id"`
}

type FrontendPublic struct {
	IpAddress []string `tfschema:"ip_address"`
}

type LoggingStorageAccount struct {
	Name          string `tfschema:"name"`
	ContainerName string `tfschema:"container_name"`
}

type NetworkInterface struct {
	SubnetId string `tfschema:"subnet_id"`
}

type DeploymentModel struct {
	ResourceGroupName      string                                     `tfschema:"resource_group_name"`
	Name                   string                                     `tfschema:"name"`
	NginxVersion           string                                     `tfschema:"nginx_version"`
	Identity               []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Sku                    string                                     `tfschema:"sku"`
	ManagedResourceGroup   string                                     `tfschema:"managed_resource_group"`
	Location               string                                     `tfschema:"location"`
	DiagnoseSupportEnabled bool                                       `tfschema:"diagnose_support_enabled"`
	IpAddress              string                                     `tfschema:"ip_address"`
	LoggingStorageAccount  []LoggingStorageAccount                    `tfschema:"logging_storage_account"`
	FrontendPublic         []FrontendPublic                           `tfschema:"frontend_public"`
	FrontendPrivate        []FrontendPrivate                          `tfschema:"frontend_private"`
	NetworkInterface       []NetworkInterface                         `tfschema:"network_interface"`
	Tags                   map[string]string                          `tfschema:"tags"`
}

func (d *DeploymentModel) LoadSDK(model *nginxdeployment.NginxDeployment) {
	d.Name = pointer.ToString(model.Name)
	d.Location = pointer.ToString(model.Location)
	d.Tags = pointer.ToMapOfStringStrings(model.Tags)
	if model.Sku != nil {
		d.Sku = model.Sku.Name
	}

	prop := model.Properties
	d.IpAddress = pointer.ToString(prop.IPAddress)
	d.ManagedResourceGroup = pointer.ToString(prop.ManagedResourceGroup)
	d.NginxVersion = pointer.ToString(prop.NginxVersion)
	d.DiagnoseSupportEnabled = pointer.ToBool(prop.EnableDiagnosticsSupport)

	if prop.Logging != nil && prop.Logging.StorageAccount != nil {
		d.LoggingStorageAccount = []LoggingStorageAccount{
			{
				Name:          pointer.ToString(prop.Logging.StorageAccount.AccountName),
				ContainerName: pointer.ToString(prop.Logging.StorageAccount.ContainerName),
			},
		}
	}

	if prop.NetworkProfile == nil {
		return
	}

	if frontend := prop.NetworkProfile.FrontEndIPConfiguration; frontend != nil {
		if publicIps := frontend.PublicIPAddresses; publicIps != nil && len(*publicIps) > 0 {
			d.FrontendPublic = append(d.FrontendPublic, FrontendPublic{})
			for _, ip := range *publicIps {
				d.FrontendPublic[0].IpAddress = append(d.FrontendPublic[0].IpAddress, pointer.ToString(ip.Id))
			}
		}

		if privateIPs := frontend.PrivateIPAddresses; privateIPs != nil && len(*privateIPs) > 0 {
			for _, ip := range *privateIPs {
				d.FrontendPrivate = append(d.FrontendPrivate, FrontendPrivate{
					IpAddress:        pointer.ToString(ip.PrivateIPAddress),
					AllocationMethod: pointer.ToString((*string)(ip.PrivateIPAllocationMethod)),
					SubnetId:         pointer.ToString(ip.SubnetId),
				})
			}
		}
	}

	if netIf := prop.NetworkProfile.NetworkInterfaceConfiguration; netIf != nil {
		d.NetworkInterface = []NetworkInterface{
			{SubnetId: pointer.ToString(netIf.SubnetId)},
		}
	}
}

type DeploymentResource struct{}

var _ sdk.ResourceWithUpdate = (*DeploymentResource)(nil)

func (m DeploymentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku": {
			// move to enum once this issue fixed: <https://github.com/Azure/azure-rest-api-specs/issues/20848>
			// docs: <https://docs.nginx.com/nginx-for-azure/billing/overview/>
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"publicpreview_Monthly_gmz7xq9ge3py",
					"standard_Monthly",
				}, false),
		},

		// only UserIdentity supported, but api defined as SystemAndUserAssigned
		// issue link: https://github.com/Azure/azure-rest-api-specs/issues/20914
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"managed_resource_group": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"diagnose_support_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ValidateFunc: nil,
		},

		"logging_storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"container_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"frontend_public": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			ConflictsWith: []string{"frontend_private"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"frontend_private": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"frontend_public"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"allocation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(nginxdeployment.PossibleValuesForNginxPrivateIPAllocationMethod(), false),
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (m DeploymentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nginx_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m DeploymentResource) ModelObject() interface{} {
	return &DeploymentModel{}
}

func (m DeploymentResource) ResourceType() string {
	return "azurerm_nginx_deployment"
}

func (m DeploymentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxDeployment

			var model DeploymentModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := nginxdeployment.NewNginxDeploymentID(subscriptionID, model.ResourceGroupName, model.Name)
			existing, err := client.DeploymentsGet(ctx, id)

			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(m.ResourceType(), id)
			}

			req := nginxdeployment.NginxDeployment{}
			req.Name = pointer.FromString(model.Name)
			req.Location = pointer.FromString(model.Location)
			req.Tags = pointer.FromMapOfStringStrings(model.Tags)

			if model.Sku != "" {
				sku := nginxdeployment.ResourceSku{Name: model.Sku}
				req.Sku = &sku
			}

			prop := &nginxdeployment.NginxDeploymentProperties{}
			prop.ManagedResourceGroup = pointer.FromString(model.ManagedResourceGroup)

			if len(model.LoggingStorageAccount) > 0 {
				prop.Logging = &nginxdeployment.NginxLogging{
					StorageAccount: &nginxdeployment.NginxStorageAccount{
						AccountName:   pointer.FromString(model.LoggingStorageAccount[0].Name),
						ContainerName: pointer.FromString(model.LoggingStorageAccount[0].ContainerName),
					},
				}
			}

			prop.EnableDiagnosticsSupport = pointer.FromBool(model.DiagnoseSupportEnabled)
			prop.NetworkProfile = &nginxdeployment.NginxNetworkProfile{
				FrontEndIPConfiguration:       &nginxdeployment.NginxFrontendIPConfiguration{},
				NetworkInterfaceConfiguration: &nginxdeployment.NginxNetworkInterfaceConfiguration{},
			}

			if public := model.FrontendPublic; len(public) > 0 && len(public[0].IpAddress) > 0 {
				var publicIPs []nginxdeployment.NginxPublicIPAddress
				for _, ip := range public[0].IpAddress {
					publicIPs = append(publicIPs, nginxdeployment.NginxPublicIPAddress{
						Id: pointer.FromString(ip),
					})
				}
				prop.NetworkProfile.FrontEndIPConfiguration.PublicIPAddresses = &publicIPs
			}

			if private := model.FrontendPrivate; len(private) > 0 {
				var privateIPs []nginxdeployment.NginxPrivateIPAddress
				for _, ip := range private {
					alloc := nginxdeployment.NginxPrivateIPAllocationMethod(ip.AllocationMethod)
					privateIPs = append(privateIPs, nginxdeployment.NginxPrivateIPAddress{
						PrivateIPAddress:          pointer.FromString(ip.IpAddress),
						PrivateIPAllocationMethod: &alloc,
						SubnetId:                  pointer.FromString(ip.SubnetId),
					})
				}
				prop.NetworkProfile.FrontEndIPConfiguration.PrivateIPAddresses = &privateIPs
			}

			if len(model.NetworkInterface) > 0 {
				prop.NetworkProfile.NetworkInterfaceConfiguration.SubnetId = pointer.FromString(model.NetworkInterface[0].SubnetId)
			}

			req.Properties = prop

			req.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding user identities: %+v", err)
			}

			err = client.DeploymentsCreateOrUpdateThenPoll(ctx, id, req)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m DeploymentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxdeployment.ParseNginxDeploymentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			client := meta.Client.Nginx.NginxDeployment
			result, err := client.DeploymentsGet(ctx, *id)
			if err != nil {
				return err
			}

			model := result.Model
			if model == nil {
				return fmt.Errorf("retrieving %s got nil model", id)
			}

			var output DeploymentModel
			output.LoadSDK(model)
			output.ResourceGroupName = id.ResourceGroupName
			output.Name = id.DeploymentName

			if id, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity); err != nil {
				return fmt.Errorf("flattening `identity`: %v", err)
			} else {
				output.Identity = *id
			}

			return meta.Encode(&output)
		},
	}
}

func (m DeploymentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: time.Minute * 30,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxdeployment.ParseNginxDeploymentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			var model DeploymentModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("Decode NginxDeploymentModel %s: %v", id, err)
			}

			var req nginxdeployment.NginxDeploymentUpdateParameters
			if meta.ResourceData.HasChange("sku") {
				req.Sku = &nginxdeployment.ResourceSku{Name: model.Sku}
			}

			if meta.ResourceData.HasChange("tags") {
				req.Tags = pointer.FromMapOfStringStrings(model.Tags)
			}

			if meta.ResourceData.HasChange("identity") {
				if req.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity); err != nil {
					return fmt.Errorf("expanding user identities: %+v", err)
				}
			}

			req.Properties = &nginxdeployment.NginxDeploymentUpdateProperties{}
			if meta.ResourceData.HasChange("logging_storage_account") && len(model.LoggingStorageAccount) > 0 {
				req.Properties.Logging = &nginxdeployment.NginxLogging{
					StorageAccount: &nginxdeployment.NginxStorageAccount{
						AccountName:   pointer.FromString(model.LoggingStorageAccount[0].Name),
						ContainerName: pointer.FromString(model.LoggingStorageAccount[0].ContainerName),
					},
				}
			}

			if meta.ResourceData.HasChange("diagnose_support_enabled") {
				req.Properties.EnableDiagnosticsSupport = pointer.FromBool(model.DiagnoseSupportEnabled)
			}

			res, err := meta.Client.Nginx.NginxDeployment.DeploymentsUpdate(ctx, *id, req)
			if err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			if err = res.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("polling update %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m DeploymentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := nginxdeployment.ParseNginxDeploymentID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)

			client := meta.Client.Nginx.NginxDeployment
			result, err := client.DeploymentsDelete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			if err := result.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting delete %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m DeploymentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nginxdeployment.ValidateNginxDeploymentID
}
