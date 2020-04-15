package network

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/peering/mgmt/2020-01-01-preview/peering"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPeerAsn() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPeerAsnCreateUpdate,
		Read:   resourceArmPeerAsnRead,
		Update: resourceArmPeerAsnCreateUpdate,
		Delete: resourceArmPeerAsnDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PeerAsnID(id)
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
				ValidateFunc: validate.PeerAsnName(),
			},

			"asn": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, math.MaxUint16),
			},

			"contact": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(peering.RoleNoc),
								string(peering.RoleOther),
								string(peering.RolePolicy),
								string(peering.RoleService),
								string(peering.RoleTechnical),
							}, false),
						},
						"email": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},

			"peer_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.PeerName(),
			},
		},
	}
}

func resourceArmPeerAsnCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PeerAsnsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to check for existing Peer Asn %q: %+v", name, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_peer_asn", *resp.ID)
		}
	}

	param := peering.PeerAsn{
		PeerAsnProperties: &peering.PeerAsnProperties{
			PeerAsn:           utils.Int32(int32(d.Get("asn").(int))),
			PeerContactDetail: expandPeerContact(d.Get("contact").(*schema.Set).List()),
			PeerName:          utils.String(d.Get("peer_name").(string)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, name, param); err != nil {
		return fmt.Errorf("failed to create Peer Asn %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to retrieving Peer Asn %q: %+v", name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Peer Asn %q ID", name)
	}
	d.SetId(*resp.ID)

	return resourceArmPeerAsnRead(d, meta)
}

func resourceArmPeerAsnRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PeerAsnsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PeerAsnID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Peer Asn %q was not found - removing from state!", id.Name)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to retrieve Peer Asn %q: %+v", id.Name, err)
	}

	d.Set("name", resp.Name)
	if props := resp.PeerAsnProperties; props != nil {
		d.Set("asn", props.PeerAsn)
		if err := d.Set("contact", flattenPeerContact(props.PeerContactDetail)); err != nil {
			return fmt.Errorf("failed to set `contact`: %+v", err)
		}
		d.Set("peer_name", props.PeerName)
	}

	return nil
}

func resourceArmPeerAsnDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PeerAsnsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PeerAsnID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.Name); err != nil {
		return fmt.Errorf("failed to delete Peer Asn %q: %+v", id.Name, err)
	}

	return nil
}

func expandPeerContact(input []interface{}) *[]peering.ContactDetail {
	result := make([]peering.ContactDetail, 0)

	for _, e := range input {
		b := e.(map[string]interface{})
		result = append(result, peering.ContactDetail{
			Role:  peering.Role(b["role"].(string)),
			Email: utils.String(b["email"].(string)),
			Phone: utils.String(b["phone"].(string)),
		})
	}

	return &result
}

func flattenPeerContact(input *[]peering.ContactDetail) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		email := ""
		if e.Email != nil {
			email = *e.Email
		}
		phone := ""
		if e.Phone != nil {
			phone = *e.Phone
		}
		output = append(output, map[string]interface{}{
			"role":  string(e.Role),
			"email": email,
			"phone": phone,
		})
	}

	return output
}
