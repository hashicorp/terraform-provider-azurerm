// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ssh_public_key

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/sshpublickeys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	computeClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing SSH Public Key.",
		Read:        dataSourceSshPublicKeyRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of this SSH Public Key.",
				Type:        pluginsdk.TypeString,
				Required:    true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-a-zA-Z0-9(_).]{1,128}$`),
					"Public SSH Key name must be 1 - 128 characters long, can contain letters, numbers, underscores, dots and hyphens (but the first and last character must be a letter or number).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"public_key": {
				Description: "The SSH public key used to authenticate to a Virtual Machine through SSH.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceSshPublicKeyRead(d *pluginsdk.ResourceData, meta any) error {
	client := computeClient.From(meta).SSHPublicKeysClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := sshpublickeys.NewSshPublicKeyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("making Read request on %s: %s", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.SshPublicKeyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		var publicKey *string
		if props := model.Properties; props.PublicKey != nil {
			publicKey, err = suppress.NormalizeSSHKey(*props.PublicKey)
			if err != nil {
				return fmt.Errorf("normalising public key: %+v", err)
			}
		}
		d.Set("public_key", publicKey)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
