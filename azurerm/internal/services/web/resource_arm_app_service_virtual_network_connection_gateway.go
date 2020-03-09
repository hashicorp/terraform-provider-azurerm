package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/Azure/go-autorest/autorest"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceVirtualNetworkConnectionGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceVirtualNetworkConnectionGatewayCreateUpdate,
		Read:   resourceArmAppServiceVirtualNetworkConnectionGatewayRead,
		Update: resourceArmAppServiceVirtualNetworkConnectionGatewayCreateUpdate,
		Delete: resourceArmAppServiceVirtualNetworkConnectionGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAppServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_network_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"virtual_network_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"certificate_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resync_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"start_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"end_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmAppServiceVirtualNetworkConnectionGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	vnetGatewayClient := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM App Service Virtual Network Connection creation.")

	resGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)
	vnetId := d.Get("virtual_network_id").(string)

	id, err := azure.ParseAzureResourceID(vnetId)
	if err != nil {
		return err
	}
	vnetResGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetVnetConnection(ctx, resGroup, appServiceName, vnetName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing App Service Virtual Network Connection for app %q vnet %q (Resource Group %q): %s", appServiceName, vnetName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_virtual_network_connection_gateway", *existing.ID)
		}
	}

	virtualNetworkGatewayId := d.Get("virtual_network_gateway_id").(string)
	gatewayId, err := azure.ParseAzureResourceID(virtualNetworkGatewayId)
	if err != nil {
		return err
	}
	gatewayResGroup := gatewayId.ResourceGroup
	gatewayName := gatewayId.Path["virtualNetworkGateways"]
	virtualNetworkGateway, err := vnetGatewayClient.Get(ctx, gatewayResGroup, gatewayName)
	if err != nil {
		return fmt.Errorf("error making Read request on AzureRM Virtual Network Gateway %q (Resource Group %q): %+v", gatewayName, gatewayResGroup, err)
	}
	if virtualNetworkGateway.VpnClientConfiguration == nil || virtualNetworkGateway.VpnClientConfiguration.VpnClientAddressPool == nil {
		return fmt.Errorf("this gateways %q under vnet %q (Resource Group %q) does not have a Point-to-site Address Range. Please specify one in CIDR notation, e.g. 10.0.0.0/8", gatewayName, vnetName, vnetResGroup)
	}

	// there are two parameters in the schema: virtual_network_id and virtual_network_gateway_id
	// we should check the virtual network gateway is within the virtual network
	isRelated, err := checkGatewayInVirtualNetwork(virtualNetworkGateway, vnetId)
	if err != nil {
		return fmt.Errorf("the virtual network gateway %q is not related with vnet %q: %+v", virtualNetworkGatewayId, vnetName, err)
	}
	if !isRelated {
		return fmt.Errorf("the virtual network gateway %q is not related with vnet %q", virtualNetworkGatewayId, vnetName)
	}

	// the create functions contains four steps:
	// 1. CreateOrUpdateVnetConnection
	// 2. result of step 1 contains cert infomation, we should set the cert to virtual network gateway (check duplicate)
	// 3. generate vpn package uri
	// 4. CreateOrUpdateVnetConnectionGateway using step 3's result

	connectionEnvelope := web.VnetInfo{
		VnetInfoProperties: &web.VnetInfoProperties{
			VnetResourceID: &vnetId,
		},
	}
	vnetInfo, err := client.CreateOrUpdateVnetConnection(ctx, resGroup, appServiceName, vnetName, connectionEnvelope)
	if err != nil {
		return fmt.Errorf("error creating/updating App Service Virtual Network Connection for app %q vnet %q (Resource Group %q): %+v", appServiceName, vnetName, resGroup, err)
	}

	// add certificate if not exists in the gateway
	if err := addCertificateIfNotExistsInGateway(vnetGatewayClient, ctx, &virtualNetworkGateway, &vnetInfo, &vnetResGroup, &gatewayName); err != nil {
		return fmt.Errorf("error add certificate for gateway %q: %+v", gatewayName, err)
	}

	// generate vpn package uri
	packageUri, err := retrieveVPNPackageUri(vnetGatewayClient, ctx, &vnetResGroup, &gatewayName)
	if err != nil {
		return fmt.Errorf("error get vpn package uri of gateway %q: %+v", gatewayName, err)
	}

	vnetGateway := web.VnetGateway{
		VnetGatewayProperties: &web.VnetGatewayProperties{
			VnetName:      &vnetName,
			VpnPackageURI: &packageUri,
		},
	}
	if _, err := client.CreateOrUpdateVnetConnectionGateway(ctx, resGroup, appServiceName, vnetName, "primary", vnetGateway); err != nil {
		return fmt.Errorf("error creating/updating App Service Virtual Network Connection gateway for app %q vnet %q (Resource Group %q): %+v", appServiceName, vnetName, resGroup, err)
	}

	d.SetId(*vnetInfo.ID)

	return resourceArmAppServiceVirtualNetworkConnectionGatewayRead(d, meta)
}

func resourceArmAppServiceVirtualNetworkConnectionGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Azure App Service Virtual Network Connection %s", id)

	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	vnetName := id.Path["virtualNetworkConnections"]
	resp, err := client.GetVnetConnection(ctx, resourceGroup, appServiceName, vnetName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Virtual Network Connection for app %q vnet %q was not found in Resource Group %q - removnig from state!", appServiceName, vnetName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error making Read request on App Service Virtual Network Connection for app %q vnet %q (Resource Group %q): %+v", appServiceName, vnetName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("name", resp.Name)
	d.Set("app_service_name", appServiceName)
	if props := resp.VnetInfoProperties; props != nil {
		d.Set("virtual_network_id", props.VnetResourceID)
		d.Set("certificate_thumbprint", props.CertThumbprint)
		d.Set("certificate_blob", props.CertBlob)
		d.Set("resync_required", props.ResyncRequired)
		if props.DNSServers != nil {
			d.Set("dns_servers", strings.Split(*props.DNSServers, ","))
		} else {
			d.Set("dns_servers", []string{})
		}
		if err := d.Set("routes", flattenAppServiceVirtualNetworkConnectionPropertiesRoutes(props.Routes)); err != nil {
			return fmt.Errorf("Error setting `routes`: %+v", err)
		}
	}

	return nil
}

func resourceArmAppServiceVirtualNetworkConnectionGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	appServiceName := id.Path["sites"]
	vnetName := id.Path["virtualNetworkConnections"]

	log.Printf("[DEBUG] Deleting App Service Virtual Network Connection for app %q vnet %q (Resource Group %q)", appServiceName, vnetName, resourceGroup)

	resp, err := client.DeleteVnetConnection(ctx, resourceGroup, appServiceName, vnetName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("error deleting App Service Virtual Network Connection for app %q vnet %q (Resource Group %q): %+v", appServiceName, vnetName, resourceGroup, err)
		}
	}

	return nil
}

func checkGatewayInVirtualNetwork(virtualNetworkGateway network.VirtualNetworkGateway, virtualNetworkId string) (bool, error) {
	if virtualNetworkGateway.IPConfigurations == nil || len(*virtualNetworkGateway.IPConfigurations) == 0 {
		return false, fmt.Errorf("error get IPConfigurations of virtual network gateway %q", *virtualNetworkGateway.ID)
	}
	subnetID := (*virtualNetworkGateway.IPConfigurations)[0].Subnet.ID
	if subnetID == nil {
		return false, fmt.Errorf("error get subnetID of virtual network gateway %q", *virtualNetworkGateway.ID)
	}
	return strings.HasPrefix(*subnetID, virtualNetworkId) && strings.HasSuffix(*subnetID, "GatewaySubnet"), nil
}

func addCertificateIfNotExistsInGateway(vnetGatewayClient *network.VirtualNetworkGatewaysClient, ctx context.Context, virtualNetworkGateway *network.VirtualNetworkGateway, vnetInfo *web.VnetInfo, resGroup *string, gatewayName *string) error {
	for _, certificate := range *virtualNetworkGateway.VpnClientConfiguration.VpnClientRootCertificates {
		if *certificate.PublicCertData == *vnetInfo.VnetInfoProperties.CertBlob {
			return nil
		}
	}
	log.Printf("[INFO] Adding certificate for virtual network gateway.")

	certName := fmt.Sprintf("AppServiceCertificate_%d.cer", tf.AccRandTimeInt())
	vpnClientRootCertToAdd := network.VpnClientRootCertificate{
		Name: &certName,
		VpnClientRootCertificatePropertiesFormat: &network.VpnClientRootCertificatePropertiesFormat{
			PublicCertData: vnetInfo.VnetInfoProperties.CertBlob,
		},
	}
	*virtualNetworkGateway.VpnClientConfiguration.VpnClientRootCertificates = append(*virtualNetworkGateway.VpnClientConfiguration.VpnClientRootCertificates, vpnClientRootCertToAdd)

	virtualNetworkGatewaysCreateOrUpdateFuture, err := vnetGatewayClient.CreateOrUpdate(ctx, *resGroup, *gatewayName, *virtualNetworkGateway)
	if err != nil {
		return fmt.Errorf("error adding cerfiticate for gateway %q (Resource Group %q): %+v", *gatewayName, *resGroup, err)
	}
	if err = virtualNetworkGatewaysCreateOrUpdateFuture.WaitForCompletionRef(ctx, vnetGatewayClient.Client); err != nil {
		return fmt.Errorf("error adding cerfiticate for gateway %q (Resource Group %q): %+v", *gatewayName, *resGroup, err)
	}
	return nil
}

func retrieveVPNPackageUri(vnetGatewayClient *network.VirtualNetworkGatewaysClient, ctx context.Context, resGroup *string, gatewayName *string) (packageUri string, err error) {
	log.Printf("[INFO] Retrieving VPN Package and supplying to App.")
	vpnClientParameters := network.VpnClientParameters{
		ProcessorArchitecture: network.Amd64,
	}
	virtualNetworkGatewaysGeneratevpnclientpackageFuture, err := vnetGatewayClient.Generatevpnclientpackage(ctx, *resGroup, *gatewayName, vpnClientParameters)
	if err != nil {
		err = fmt.Errorf("error generating vpn client package for Virtual Network Gateway %q (Resource Group %q): vpnClientParameters %+v %+v", *gatewayName, *resGroup, vpnClientParameters, err)
		return
	}
	if err = virtualNetworkGatewaysGeneratevpnclientpackageFuture.WaitForCompletionRef(ctx, vnetGatewayClient.Client); err != nil {
		err = fmt.Errorf("error waiting the result of generating vpn client package for Virtual Network Gateway %q (Resource Group %q): vpnClientParameters %+v %+v", *gatewayName, *resGroup, vpnClientParameters, err)
		return
	}
	packageUri, err = getResult(*vnetGatewayClient, virtualNetworkGatewaysGeneratevpnclientpackageFuture)
	if err != nil {
		err = fmt.Errorf("error getting result vpn client package for Virtual Network Gateway %q (Resource Group %q): vpnClientParameters %+v %+v", *gatewayName, *resGroup, vpnClientParameters, err)
		return
	}
	if len(packageUri) > 0 && packageUri[0] == '"' && packageUri[len(packageUri)-1] == '"' {
		packageUri = packageUri[1 : len(packageUri)-1]
	}
	return
}

func getResult(client network.VirtualNetworkGatewaysClient, future network.VirtualNetworkGatewaysGeneratevpnclientpackageFuture) (str string, err error) {
	var done bool
	done, err = future.DoneWithContext(context.Background(), client)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewaysGeneratevpnclientpackageFuture", "Result", future.Response(), "Polling failure")
		return
	}
	if !done {
		err = autorestAzure.NewAsyncOpIncompleteError("network.VirtualNetworkGatewaysGeneratevpnclientpackageFuture")
		return
	}
	sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
	var s network.String
	if s.Response.Response, err = future.GetResult(sender); err == nil && s.Response.Response.StatusCode != http.StatusNoContent {
		str, err = generatevpnclientpackageResponder(client, s.Response.Response)
		if err != nil {
			err = autorest.NewErrorWithError(err, "network.VirtualNetworkGatewaysGeneratevpnclientpackageFuture", "Result", s.Response.Response, "Failure responding to request")
		}
	}
	return
}

func generatevpnclientpackageResponder(client network.VirtualNetworkGatewaysClient, resp *http.Response) (result string, err error) {
	byteArr := make([]byte, 1024)
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		autorestAzure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted),
		autorest.ByUnmarshallingBytes(&byteArr),
		autorest.ByClosing())
	result = string(byteArr)
	return
}

func flattenAppServiceVirtualNetworkConnectionPropertiesRoutes(input *[]web.VnetRoute) []interface{} {
	if input == nil {
		return nil
	}

	routes := make([]interface{}, 0)
	for _, route := range *input {
		attr := make(map[string]interface{})
		if route.Name != nil {
			attr["name"] = *route.Name
		}
		if props := route.VnetRouteProperties; props != nil {
			if props.StartAddress != nil {
				attr["start_address"] = *props.StartAddress
			}
			if props.EndAddress != nil {
				attr["end_address"] = *props.EndAddress
			}
			attr["route_type"] = string(props.RouteType)
		}

		routes = append(routes, attr)
	}
	return routes
}
