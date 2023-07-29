package iotcentral

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	iotcentralDataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

func resourceIotCentralServicePrincipalUser() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotCentralServicePrincipalUserCreate,
		Read:   resourceIotCentralServicePrincipalUserRead,
		Update: resourceIotCentralServicePrincipalUserUpdate,
		Delete: resourceIotCentralServicePrincipalUserDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ParseNestedItemID(id)
			return err
		}),

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
				ForceNew: true,
			},
			"object_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": schemaRole(),
		},
	}
}

func resourceIotCentralServicePrincipalUserCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := uuid.New().String()

	userClient, err := client.UsersClient(ctx, d.Get("sub_domain").(string))
	if err != nil {
		return fmt.Errorf("creating users client: %+v", err)
	}

	objectId := d.Get("object_id").(string)
	tenantId := d.Get("tenant_id").(string)

	roleData, ok := d.GetOk("roles")
	if !ok {
		return fmt.Errorf("'roles' is not specified")
	}

	roleAssignments := convertToRoleAssignments(roleData.([]interface{}))

	model := iotcentralDataplane.ServicePrincipalUser{
		ObjectID: &objectId,
		TenantID: &tenantId,
		Roles:    &roleAssignments,
	}

	resp, err := userClient.Create(ctx, id, model)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	baseUrl, err := parse.ParseBaseUrl(d.Get("sub_domain").(string))
	if err != nil {
		return err
	}

	user, isServicePrincipalUser := resp.Value.AsServicePrincipalUser()
	if !isServicePrincipalUser {
		return fmt.Errorf("expected %s to be an service principal user", id)
	}

	userId, err := parse.NewNestedItemID(baseUrl, parse.NestedItemTypeOrganization, *user.ID)
	if err != nil {
		return err
	}

	d.SetId(userId.ID())
	return resourceIotCentralServicePrincipalUserRead(d, meta)
}

func resourceIotCentralServicePrincipalUserRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	userClient, err := client.UsersClient(ctx, id.SubDomain)
	if err != nil {
		return fmt.Errorf("creating users client: %+v", err)
	}

	resp, err := userClient.Get(ctx, id.Id)
	if err != nil {
		return err
	}

	user, isServicePrincipalUser := resp.Value.AsServicePrincipalUser()
	if !isServicePrincipalUser {
		return fmt.Errorf("expected %s to be an service principal user", id)
	}

	d.Set("sub_domain", id.SubDomain)
	d.Set("object_id", user.ObjectID)
	d.Set("tenant_id", user.TenantID)
	d.Set("roles", convertFromRoleAssignments(*user.Roles))

	return nil
}

func resourceIotCentralServicePrincipalUserUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	userClient, err := client.UsersClient(ctx, id.SubDomain)
	if err != nil {
		return fmt.Errorf("creating users client: %+v", err)
	}

	resp, err := userClient.Get(ctx, id.Id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	existing, isServicePrincipalUser := resp.Value.AsServicePrincipalUser()
	if !isServicePrincipalUser {
		return fmt.Errorf("expected %s to be an service principal user", id)
	}

	objectId := d.Get("object_id").(string)
	if d.HasChange("object_id") {
		existing.ObjectID = &objectId
	}

	tenantId := d.Get("tenant_id").(string)
	if d.HasChange("tenant_id") {
		existing.TenantID = &tenantId
	}

	roleAssignments := convertToRoleAssignments(d.Get("roles").([]interface{}))
	if d.HasChange("roles") {
		existing.Roles = &roleAssignments
	}

	_, err = userClient.Update(ctx, *existing.ID, existing, "*")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceIotCentralServicePrincipalUserRead(d, meta)
}

func resourceIotCentralServicePrincipalUserDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTCentral
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ParseNestedItemID(d.Id())
	if err != nil {
		return err
	}

	userClient, err := client.UsersClient(ctx, id.SubDomain)
	if err != nil {
		return fmt.Errorf("creating users client: %+v", err)
	}

	_, err = userClient.Remove(ctx, id.Id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
