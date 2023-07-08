package iotcentral

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	iotcentralDataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

func resourceIotCentralOrganization() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotCentralOrganizationCreate,
		Read:   resourceIotCentralOrganizationRead,
		Update: resourceIotCentralOrganizationUpdate,
		Delete: resourceIotCentralOrganizationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"sub_domain": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"parent": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIotCentralOrganizationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := uuid.New().String()

	orgClient, err := client.OrganizationsClient(ctx, d.Get("sub_domain").(string))
	if err != nil {
		return fmt.Errorf("creating organization client: %+v", err)
	}

	displayName := d.Get("display_name").(string)
	parent := d.Get("parent").(string)

	model := iotcentralDataplane.Organization{
		DisplayName: &displayName,
	}

	if parent != "" {
		model.Parent = &parent
	}

	org, err := orgClient.Create(ctx, id, model)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(*org.ID)
	return resourceIotCentralOrganizationRead(d, meta)
}

func resourceIotCentralOrganizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Get("id").(string)

	orgClient, err := client.OrganizationsClient(ctx, d.Get("sub_domain").(string))
	if err != nil {
		return fmt.Errorf("creating organization client: %+v", err)
	}

	org, err := orgClient.Get(ctx, id)
	if err != nil {
		if org.ID == nil || *org.ID == "" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *org.ID, err)
	}

	d.SetId(*org.ID)
	d.Set("display_name", org.DisplayName)
	d.Set("parent", org.Parent)

	return nil
}

func resourceIotCentralOrganizationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Get("id").(string)

	orgClient, err := client.OrganizationsClient(ctx, d.Get("sub_domain").(string))
	if err != nil {
		return fmt.Errorf("creating organization client: %+v", err)
	}

	existing, err := orgClient.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	displayName := d.Get("display_name").(string)
	if d.HasChange("display_name") {
		existing.DisplayName = &displayName
	}

	parent := d.Get("parent").(string)
	if d.HasChange("parent") {
		existing.Parent = &parent
	}

	_, err = orgClient.Update(ctx, *existing.ID, existing, "*")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceIotCentralOrganizationRead(d, meta)
}

func resourceIotCentralOrganizationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Get("id").(string)

	orgClient, err := client.OrganizationsClient(ctx, d.Get("sub_domain").(string))
	if err != nil {
		return fmt.Errorf("creating organization client: %+v", err)
	}

	_, err = orgClient.Remove(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
