package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVPNServerConfigurationPolicyGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVPNServerConfigurationPolicyGroupCreateUpdate,
		Read:   resourceVPNServerConfigurationPolicyGroupRead,
		Update: resourceVPNServerConfigurationPolicyGroupCreateUpdate,
		Delete: resourceVPNServerConfigurationPolicyGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VpnServerConfigurationPolicyGroupID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vpn_server_configuration_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VpnServerConfigurationID,
			},

			"policy": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.VpnPolicyMemberAttributeTypeAADGroupID),
								string(network.VpnPolicyMemberAttributeTypeCertificateGroupID),
								string(network.VpnPolicyMemberAttributeTypeRadiusAzureGroupID),
							}, false),
						},

						"value": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"is_default": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"priority": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceVPNServerConfigurationPolicyGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.ConfigurationPolicyGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	vpnServerConfigurationId, err := parse.VpnServerConfigurationID(d.Get("vpn_server_configuration_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewVpnServerConfigurationPolicyGroupID(subscriptionId, vpnServerConfigurationId.ResourceGroup, vpnServerConfigurationId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_vpn_server_configuration_policy_group", id.ID())
		}
	}

	props := network.VpnServerConfigurationPolicyGroup{
		Name: utils.String(d.Get("name").(string)),
		VpnServerConfigurationPolicyGroupProperties: &network.VpnServerConfigurationPolicyGroupProperties{
			IsDefault:     utils.Bool(d.Get("is_default").(bool)),
			PolicyMembers: expandVPNServerConfigurationPolicyGroupPolicyMembers(d.Get("policy").(*pluginsdk.Set).List()),
			Priority:      utils.Int32(int32(d.Get("priority").(int))),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName, props)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVPNServerConfigurationPolicyGroupRead(d, meta)
}

func resourceVPNServerConfigurationPolicyGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConfigurationPolicyGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnServerConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConfigurationPolicyGroupName)

	vpnServerConfigurationId := parse.NewVpnServerConfigurationID(id.SubscriptionId, id.ResourceGroup, id.VpnServerConfigurationName)
	d.Set("vpn_server_configuration_id", vpnServerConfigurationId.ID())

	if props := resp.VpnServerConfigurationPolicyGroupProperties; props != nil {
		d.Set("is_default", props.IsDefault)
		d.Set("priority", props.Priority)

		if err := d.Set("policy", flattenVPNServerConfigurationPolicyGroupPolicyMembers(props.PolicyMembers)); err != nil {
			return fmt.Errorf("setting `policy`: %+v", err)
		}
	}

	return nil
}

func resourceVPNServerConfigurationPolicyGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ConfigurationPolicyGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnServerConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the %s: %+v", id, err)
	}

	return nil
}

func expandVPNServerConfigurationPolicyGroupPolicyMembers(input []interface{}) *[]network.VpnServerConfigurationPolicyGroupMember {
	results := make([]network.VpnServerConfigurationPolicyGroupMember, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, network.VpnServerConfigurationPolicyGroupMember{
			Name:           utils.String(v["name"].(string)),
			AttributeType:  network.VpnPolicyMemberAttributeType(v["type"].(string)),
			AttributeValue: utils.String(v["value"].(string)),
		})
	}

	return &results
}

func flattenVPNServerConfigurationPolicyGroupPolicyMembers(input *[]network.VpnServerConfigurationPolicyGroupMember) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var attributeType network.VpnPolicyMemberAttributeType
		if item.AttributeType != "" {
			attributeType = item.AttributeType
		}

		var attributeValue string
		if item.AttributeValue != nil {
			attributeValue = *item.AttributeValue
		}

		results = append(results, map[string]interface{}{
			"name":  name,
			"type":  attributeType,
			"value": attributeValue,
		})
	}

	return results
}
