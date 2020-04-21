package blueprints

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceArmBlueprintAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBlueprintAssignmentCreateUpdate,
		Update: nil,
		Read:   resourceArmBlueprintAssignmentRead,
		Delete: resourceArmBlueprintAssignmentDelete,

		Importer: nil,

		Timeouts: nil,

		Schema: nil,
	}
}

func resourceArmBlueprintAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceArmBlueprintAssignmentRead(d, meta)
}

func resourceArmBlueprintAssignmentRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceArmBlueprintAssignmentDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
