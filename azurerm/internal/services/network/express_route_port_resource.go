package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var expressRoutePortSchema = &schema.Schema{
	Type: schema.TypeList,
	// Service will always create a pair of links automatically. Users can't add or remove link, but only manipulate existing ones.
	// This is because the link is actually a map to the physical pair of ports on the MS edge device.
	Optional: true,
	Computed: true,
	MinItems: 1,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"admin_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"macsec_cipher": {
				Type:     schema.TypeString,
				Optional: true,

				// TODO: The following hardcode can be replaced by SDK types once following is merged:
				// 	https://github.com/Azure/azure-rest-api-specs/pull/12329
				Default: "GcmAes128",
				// Default: string(network.GcmAes128),

				// TODO: The following hardcode can be replaced by SDK types once following is merged:
				// 	https://github.com/Azure/azure-rest-api-specs/pull/12329
				ValidateFunc: validation.StringInSlice([]string{
					"GcmAes128",
					"GcmAes256",
					// string(network.GcmAes128),
					// string(network.GcmAes256),
				}, false),
			},
			"macsec_ckn_keyvault_secret_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"macsec_cak_keyvault_secret_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"router_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"patch_panel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rack_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connector_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	},
}

func resourceArmExpressRoutePort() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmExpressRoutePortCreateUpdate,
		Read:   resourceArmExpressRoutePortRead,
		Update: resourceArmExpressRoutePortCreateUpdate,
		Delete: resourceArmExpressRoutePortDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRoutePortID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRoutePortName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"peering_location": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"bandwidth_in_gbps": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"encapsulation": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Dot1Q),
					string(network.QinQ),
				}, false),
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"identity_ids": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateFunc:     validation.NoZeroValues,
								DiffSuppressFunc: suppress.CaseDifference,
							},
						},

						"type": {
							Type:     schema.TypeString,
							Required: true,
							// TODO: The "ignoreCase" and diff suppression function can be removed once
							// following issue get resolved:
							// https://github.com/Azure/azure-rest-api-specs/issues/12330
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ResourceIdentityTypeUserAssigned),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"link1": expressRoutePortSchema,

			"link2": expressRoutePortSchema,

			"ethertype": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mtu": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmExpressRoutePortCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	id := parse.NewExpressRoutePortID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Express Route Port %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_express_route_port", id.ID())
		}
	}

	param := network.ExpressRoutePort{
		Name:     &name,
		Location: &location,
		ExpressRoutePortPropertiesFormat: &network.ExpressRoutePortPropertiesFormat{
			PeeringLocation: utils.String(d.Get("peering_location").(string)),
			BandwidthInGbps: utils.Int32(int32(d.Get("bandwidth_in_gbps").(int))),
			Encapsulation:   network.ExpressRoutePortsEncapsulation(d.Get("encapsulation").(string)),
		},
		Identity: expandExpressRoutePortIdentity(d.Get("identity").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	// The link properties can't be specified in first creation. It will result into either error (e.g. setting `adminState`) or being ignored (e.g. setting MACSec)
	// Hence, if this is a new creation we will do a create-then-update here.
	if d.IsNewResource() {
		future, err := client.CreateOrUpdate(ctx, resourceGroup, name, param)
		if err != nil {
			return fmt.Errorf("creating Express Route Port %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for creation of Express Route Port %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	param.ExpressRoutePortPropertiesFormat.Links = expandExpressRoutePortLinks(d.Get("link1").([]interface{}), d.Get("link2").([]interface{}))

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, param)
	if err != nil {
		return fmt.Errorf("creating Express Route Port %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Express Route Port %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceArmExpressRoutePortRead(d, meta)
}

func resourceArmExpressRoutePortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRoutePortID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Express Route Port %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Express Route Port %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if err := d.Set("identity", flattenExpressRoutePortIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("error setting `identity`: %v", err)
	}
	if prop := resp.ExpressRoutePortPropertiesFormat; prop != nil {
		d.Set("peering_location", prop.PeeringLocation)
		d.Set("bandwidth_in_gbps", prop.BandwidthInGbps)
		d.Set("encapsulation", prop.Encapsulation)
		link1, link2, err := flattenExpressRoutePortLinks(resp.Links)
		if err != nil {
			return fmt.Errorf("error flattening links: %v", err)
		}
		if err := d.Set("link1", link1); err != nil {
			return fmt.Errorf("error setting `link1`: %v", err)
		}
		if err := d.Set("link2", link2); err != nil {
			return fmt.Errorf("error setting `link2`: %v", err)
		}
		d.Set("ethertype", prop.EtherType)
		d.Set("guid", prop.ResourceGUID)
		d.Set("mtu", prop.Mtu)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmExpressRoutePortDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRoutePortsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRoutePortID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Express Route Port %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("waiting for deletion of Express Route Port %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandExpressRoutePortIdentity(input []interface{}) *network.ManagedServiceIdentity {
	if len(input) == 0 {
		return nil
	}
	identity := input[0].(map[string]interface{})
	identityType := network.ResourceIdentityType(identity["type"].(string))

	managedServiceIdentity := network.ManagedServiceIdentity{
		Type: identityType,
	}

	identityIds := make(map[string]*network.ManagedServiceIdentityUserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &network.ManagedServiceIdentityUserAssignedIdentitiesValue{}
	}

	// TODO: once following issue get resolved, can directly equal test:
	// https://github.com/Azure/azure-rest-api-specs/issues/12330
	if strings.EqualFold(string(managedServiceIdentity.Type), string(network.ResourceIdentityTypeUserAssigned)) ||
		strings.EqualFold(string(managedServiceIdentity.Type), string(network.ResourceIdentityTypeSystemAssignedUserAssigned)) {
		managedServiceIdentity.UserAssignedIdentities = identityIds
	}

	return &managedServiceIdentity
}

func flattenExpressRoutePortIdentity(identity *network.ManagedServiceIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	identityIds := make([]string, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			identityIds = append(identityIds, key)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"identity_ids": identityIds,
			"type":         string(identity.Type),
		},
	}
}

func expandExpressRoutePortLinks(link1, link2 []interface{}) *[]network.ExpressRouteLink {
	var out []network.ExpressRouteLink
	if link := expandExpressRoutePortLink(1, link1); link != nil {
		out = append(out, *link)
	}
	if link := expandExpressRoutePortLink(2, link2); link != nil {
		out = append(out, *link)
	}
	if len(out) == 0 {
		return nil
	}
	return &out
}

func expandExpressRoutePortLink(idx int, input []interface{}) *network.ExpressRouteLink {
	if len(input) == 0 {
		return nil
	}

	b := input[0].(map[string]interface{})
	adminState := network.ExpressRouteLinkAdminStateDisabled
	if b["admin_enabled"].(bool) {
		adminState = network.ExpressRouteLinkAdminStateEnabled
	}

	link := network.ExpressRouteLink{
		// The link name is fixed
		Name: utils.String(fmt.Sprintf("link%d", idx)),
		ExpressRouteLinkPropertiesFormat: &network.ExpressRouteLinkPropertiesFormat{
			AdminState: adminState,
			MacSecConfig: &network.ExpressRouteLinkMacSecConfig{
				Cipher: network.ExpressRouteLinkMacSecCipher(b["macsec_cipher"].(string)),
			},
		},
	}

	if cknSecretId := b["macsec_ckn_keyvault_secret_id"].(string); cknSecretId != "" {
		link.ExpressRouteLinkPropertiesFormat.MacSecConfig.CknSecretIdentifier = &cknSecretId
	}
	if cakSecretId := b["macsec_cak_keyvault_secret_id"].(string); cakSecretId != "" {
		link.ExpressRouteLinkPropertiesFormat.MacSecConfig.CakSecretIdentifier = &cakSecretId
	}
	return &link
}

func flattenExpressRoutePortLinks(links *[]network.ExpressRouteLink) ([]interface{}, []interface{}, error) {
	if links == nil {
		return nil, nil, nil
	}
	length := len(*links)
	if length != 2 {
		return nil, nil, fmt.Errorf("expected two links, but got %d", length)
	}

	return flattenExpressRoutePortLink((*links)[0]), flattenExpressRoutePortLink((*links)[1]), nil
}

func flattenExpressRoutePortLink(link network.ExpressRouteLink) []interface{} {
	var id string
	if link.ID != nil {
		id = *link.ID
	}

	var (
		routerName    string
		interfaceName string
		patchPanelId  string
		rackId        string
		connectorType string
		adminState    bool
		cknSecretId   string
		cakSecretId   string
		cipher        string
	)

	if prop := link.ExpressRouteLinkPropertiesFormat; prop != nil {
		if prop.RouterName != nil {
			routerName = *prop.RouterName
		}
		if prop.InterfaceName != nil {
			interfaceName = *prop.InterfaceName
		}
		if prop.PatchPanelID != nil {
			patchPanelId = *prop.PatchPanelID
		}
		if prop.RackID != nil {
			rackId = *prop.RackID
		}
		connectorType = string(prop.ConnectorType)
		adminState = prop.AdminState == network.ExpressRouteLinkAdminStateEnabled
		if cfg := prop.MacSecConfig; cfg != nil {
			if cfg.CknSecretIdentifier != nil {
				cknSecretId = *cfg.CknSecretIdentifier
			}
			if cfg.CakSecretIdentifier != nil {
				cakSecretId = *cfg.CakSecretIdentifier
			}
			cipher = string(cfg.Cipher)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"id":                            id,
			"router_name":                   routerName,
			"interface_name":                interfaceName,
			"patch_panel_id":                patchPanelId,
			"rack_id":                       rackId,
			"connector_type":                connectorType,
			"admin_enabled":                 adminState,
			"macsec_ckn_keyvault_secret_id": cknSecretId,
			"macsec_cak_keyvault_secret_id": cakSecretId,
			"macsec_cipher":                 cipher,
		},
	}
}
