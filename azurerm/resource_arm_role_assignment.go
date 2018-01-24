package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRoleAssignmentCreate,
		Read:   resourceArmRoleAssignmentRead,
		Delete: resourceArmRoleAssignmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"role_definition_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)
	roleDefinitionId := d.Get("role_definition_id").(string)
	principalId := d.Get("principal_id").(string)

	if name == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Role Assignment: %+v", err)
		}

		name = uuid
	}

	properties := authorization.RoleAssignmentCreateParameters{
		Properties: &authorization.RoleAssignmentProperties{
			RoleDefinitionID: utils.String(roleDefinitionId),
			PrincipalID:      utils.String(principalId),
		},
	}

	_, err := client.Create(ctx, scope, name, properties)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, scope, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Role Assignment ID for %q (Scope %q)", name, scope)
	}

	d.SetId(*read.ID)
	return resourceArmRoleAssignmentRead(d, meta)
}

func resourceArmRoleAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.GetByID(ctx, d.Id())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Role Assignment ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading Role Assignment %q: %+v", d.Id(), err)
	}

	d.Set("name", resp.Name)

	if props := resp.Properties; props != nil {
		d.Set("scope", props.Scope)
		d.Set("role_definition_id", props.RoleDefinitionID)
		d.Set("principal_id", props.PrincipalID)
	}

	return nil
}

func resourceArmRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	scope := d.Get("scope").(string)
	name := d.Get("name").(string)

	resp, err := client.Delete(ctx, scope, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return err
		}
	}

	return nil
}
