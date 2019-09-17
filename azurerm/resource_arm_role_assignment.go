package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2018-09-01-preview/authorization"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/resource"
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
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"role_definition_name"},
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"role_definition_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"role_definition_id"},
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validateRoleDefinitionName,
			},

			"principal_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"principal_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"skip_service_principal_aad_check": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceArmRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	roleAssignmentsClient := meta.(*ArmClient).authorization.RoleAssignmentsClient
	roleDefinitionsClient := meta.(*ArmClient).authorization.RoleDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	var roleDefinitionId string
	if v, ok := d.GetOk("role_definition_id"); ok {
		roleDefinitionId = v.(string)
	} else if v, ok := d.GetOk("role_definition_name"); ok {
		roleName := v.(string)
		roleDefinitions, err := roleDefinitionsClient.List(ctx, scope, fmt.Sprintf("roleName eq '%s'", roleName))
		if err != nil {
			return fmt.Errorf("Error loading Role Definition List: %+v", err)
		}
		if len(roleDefinitions.Values()) != 1 {
			return fmt.Errorf("Error loading Role Definition List: could not find role '%s'", roleName)
		}
		roleDefinitionId = *roleDefinitions.Values()[0].ID
	} else {
		return fmt.Errorf("Error: either role_definition_id or role_definition_name needs to be set")
	}
	d.Set("role_definition_id", roleDefinitionId)

	principalId := d.Get("principal_id").(string)

	if name == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Role Assignment: %+v", err)
		}

		name = uuid
	}

	if features.ShouldResourcesBeImported() {
		existing, err := roleAssignmentsClient.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Role Assignment ID for %q (Scope %q): %+v", name, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_role_assignment", *existing.ID)
		}
	}

	properties := authorization.RoleAssignmentCreateParameters{
		RoleAssignmentProperties: &authorization.RoleAssignmentProperties{
			RoleDefinitionID: utils.String(roleDefinitionId),
			PrincipalID:      utils.String(principalId),
		},
	}

	skipPrincipalCheck := d.Get("skip_service_principal_aad_check").(bool)

	if skipPrincipalCheck {
		properties.RoleAssignmentProperties.PrincipalType = authorization.ServicePrincipal
	}

	if err := resource.Retry(300*time.Second, retryRoleAssignmentsClient(scope, name, properties, meta)); err != nil {
		return err
	}

	read, err := roleAssignmentsClient.Get(ctx, scope, name)
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
	client := meta.(*ArmClient).authorization.RoleAssignmentsClient
	roleDefinitionsClient := meta.(*ArmClient).authorization.RoleDefinitionsClient
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

	if props := resp.RoleAssignmentPropertiesWithScope; props != nil {
		d.Set("scope", props.Scope)
		d.Set("role_definition_id", props.RoleDefinitionID)
		d.Set("principal_id", props.PrincipalID)
		d.Set("principal_type", props.PrincipalType)

		//allows for import when role name is used (also if the role name changes a plan will show a diff)
		if roleId := props.RoleDefinitionID; roleId != nil {
			roleResp, err := roleDefinitionsClient.GetByID(ctx, *roleId)
			if err != nil {
				return fmt.Errorf("Error loading Role Definition %q: %+v", *roleId, err)
			}

			if roleProps := roleResp.RoleDefinitionProperties; roleProps != nil {
				d.Set("role_definition_name", roleProps.RoleName)
			}
		}
	}

	return nil
}

func resourceArmRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).authorization.RoleAssignmentsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseRoleAssignmentId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.scope, id.name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return err
		}
	}

	return nil
}

func validateRoleDefinitionName(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}

	if ok := strings.Contains(v, "(Preview)"); ok {
		return nil, []error{fmt.Errorf("Preview roles are not supported")}
	}
	return nil, nil
}

func retryRoleAssignmentsClient(scope string, name string, properties authorization.RoleAssignmentCreateParameters, meta interface{}) func() *resource.RetryError {
	return func() *resource.RetryError {
		roleAssignmentsClient := meta.(*ArmClient).authorization.RoleAssignmentsClient
		ctx := meta.(*ArmClient).StopContext

		resp, err := roleAssignmentsClient.Create(ctx, scope, name, properties)
		if err != nil {
			if utils.ResponseErrorIsRetryable(err) {
				return resource.RetryableError(err)
			} else if resp.Response.StatusCode == 400 && strings.Contains(err.Error(), "PrincipalNotFound") {
				// When waiting for service principal to become available
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	}
}

type roleAssignmentId struct {
	scope string
	name  string
}

func parseRoleAssignmentId(input string) (*roleAssignmentId, error) {
	segments := strings.Split(input, "/providers/Microsoft.Authorization/roleAssignments/")
	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected Role Assignment ID to be in the format `{scope}/providers/Microsoft.Authorization/roleAssignments/{name}` but got %q", input)
	}

	// /{scope}/providers/Microsoft.Authorization/roleAssignments/{roleAssignmentName}
	id := roleAssignmentId{
		scope: strings.TrimPrefix(segments[0], "/"),
		name:  segments[1],
	}
	return &id, nil
}
