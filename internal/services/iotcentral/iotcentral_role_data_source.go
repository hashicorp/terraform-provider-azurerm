package iotcentral

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	iotcentralDataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

func dataSourceIotCentralRole() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceIotCentralRoleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"sub_domain": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIotCentralRoleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	if client == nil {
		return fmt.Errorf("client is nil")
	}
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rolesClient, err := client.RolesClient(ctx, d.Get("sub_domain").(string))
	if err != nil {
		return fmt.Errorf("creating roles client: %+v", err)
	}

	resp, err := rolesClient.ListComplete(ctx)
	if err != nil {
		return fmt.Errorf("retrieving roles: %+v", err)
	}

	cleanedName := cleanNameInconsistency(d.Get("display_name").(string))

	role := iotcentralDataplane.Role{}
	for resp.NotDone() {
		role = resp.Value()
		if *role.DisplayName == cleanedName {
			break
		}

		// Advance to the next value in the iterator
		err := resp.NextWithContext(ctx)
		if err != nil {
			// handle error
			fmt.Println("Error advancing iterator:", err)
			break
		}
	}

	if role.ID == nil || *role.ID == "" || *role.DisplayName != cleanedName {
		return fmt.Errorf("role not found for name: %s", d.Get("display_name"))
	}

	d.SetId(*role.ID)
	d.Set("id", role.ID)

	return nil
}

func cleanNameInconsistency(name string) string {
	// clean up inconsistencies in the role names, as they are displayed differently in the portal
	switch name {
	case "App Administrator", "App Builder", "App Operator":
		name = strings.TrimPrefix(name, "App ")
	case "Org Administrator":
		name = "Org Admin"
	}

	return name
}
