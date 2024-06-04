// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVPNServerConfigurationPolicyGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVPNServerConfigurationPolicyGroupCreate,
		Read:   resourceVPNServerConfigurationPolicyGroupRead,
		Update: resourceVPNServerConfigurationPolicyGroupUpdate,
		Delete: resourceVPNServerConfigurationPolicyGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseConfigurationPolicyGroupID(id)
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
				ValidateFunc: virtualwans.ValidateVpnServerConfigurationID,
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
								string(virtualwans.VpnPolicyMemberAttributeTypeAADGroupId),
								string(virtualwans.VpnPolicyMemberAttributeTypeCertificateGroupId),
								string(virtualwans.VpnPolicyMemberAttributeTypeRadiusAzureGroupId),
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

func resourceVPNServerConfigurationPolicyGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vpnServerConfigurationId, err := virtualwans.ParseVpnServerConfigurationID(d.Get("vpn_server_configuration_id").(string))
	if err != nil {
		return err
	}

	locks.ByID(vpnServerConfigurationId.ID())
	defer locks.UnlockByID(vpnServerConfigurationId.ID())

	id := virtualwans.NewConfigurationPolicyGroupID(subscriptionId, vpnServerConfigurationId.ResourceGroupName, vpnServerConfigurationId.VpnServerConfigurationName, d.Get("name").(string))

	existing, err := client.ConfigurationPolicyGroupsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_vpn_server_configuration_policy_group", id.ID())
	}

	payload := virtualwans.VpnServerConfigurationPolicyGroup{
		Properties: &virtualwans.VpnServerConfigurationPolicyGroupProperties{
			IsDefault:     utils.Bool(d.Get("is_default").(bool)),
			PolicyMembers: expandVPNServerConfigurationPolicyGroupPolicyMembers(d.Get("policy").(*pluginsdk.Set).List()),
			Priority:      pointer.To(int64(d.Get("priority").(int))),
		},
	}

	if err := client.ConfigurationPolicyGroupsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVPNServerConfigurationPolicyGroupRead(d, meta)
}

func resourceVPNServerConfigurationPolicyGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ConfigurationPolicyGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConfigurationPolicyGroupName)

	vpnServerConfigurationId := virtualwans.NewVpnServerConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName)
	d.Set("vpn_server_configuration_id", vpnServerConfigurationId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("is_default", props.IsDefault)
			d.Set("priority", props.Priority)

			if err := d.Set("policy", flattenVPNServerConfigurationPolicyGroupPolicyMembers(props.PolicyMembers)); err != nil {
				return fmt.Errorf("setting `policy`: %+v", err)
			}
		}
	}

	return nil
}

func resourceVPNServerConfigurationPolicyGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	vpnServerConfigurationId := virtualwans.NewVpnServerConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName)

	locks.ByID(vpnServerConfigurationId.ID())
	defer locks.UnlockByID(vpnServerConfigurationId.ID())

	existing, err := client.ConfigurationPolicyGroupsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("policy") {
		payload.Properties.PolicyMembers = expandVPNServerConfigurationPolicyGroupPolicyMembers(d.Get("policy").(*pluginsdk.Set).List())
	}

	if d.HasChange("priority") {
		payload.Properties.Priority = pointer.To(int64(d.Get("priority").(int)))
	}

	if err := client.ConfigurationPolicyGroupsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVPNServerConfigurationPolicyGroupRead(d, meta)
}

func resourceVPNServerConfigurationPolicyGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseConfigurationPolicyGroupID(d.Id())
	if err != nil {
		return err
	}

	vpnServerConfigurationId := virtualwans.NewVpnServerConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName)

	locks.ByID(vpnServerConfigurationId.ID())
	defer locks.UnlockByID(vpnServerConfigurationId.ID())

	if err := client.ConfigurationPolicyGroupsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandVPNServerConfigurationPolicyGroupPolicyMembers(input []interface{}) *[]virtualwans.VpnServerConfigurationPolicyGroupMember {
	results := make([]virtualwans.VpnServerConfigurationPolicyGroupMember, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, virtualwans.VpnServerConfigurationPolicyGroupMember{
			Name:           utils.String(v["name"].(string)),
			AttributeType:  pointer.To(virtualwans.VpnPolicyMemberAttributeType(v["type"].(string))),
			AttributeValue: utils.String(v["value"].(string)),
		})
	}

	return &results
}

func flattenVPNServerConfigurationPolicyGroupPolicyMembers(input *[]virtualwans.VpnServerConfigurationPolicyGroupMember) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		attributeType := ""
		if item.AttributeType != nil {
			attributeType = string(*item.AttributeType)
		}

		results = append(results, map[string]interface{}{
			"name":  pointer.From(item.Name),
			"type":  attributeType,
			"value": pointer.From(item.AttributeValue),
		})
	}

	return results
}
