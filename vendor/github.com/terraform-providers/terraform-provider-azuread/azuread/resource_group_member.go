package azuread

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

const groupMemberResourceName = "azuread_group_member"

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupMemberCreate,
		Read:   resourceGroupMemberRead,
		Delete: resourceGroupMemberDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"member_object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},
		},
	}
}

func resourceGroupMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	groupID := d.Get("group_object_id").(string)
	memberID := d.Get("member_object_id").(string)

	tf.LockByName(groupMemberResourceName, groupID)
	defer tf.UnlockByName(groupMemberResourceName, groupID)

	if err := graph.GroupAddMember(client, ctx, groupID, memberID); err != nil {
		return err
	}

	d.SetId(graph.GroupMemberIdFrom(groupID, memberID).String())
	return resourceGroupMemberRead(d, meta)
}

func resourceGroupMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := graph.ParseGroupMemberId(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to parse ID: %v", err)
	}

	members, err := graph.GroupAllMembers(client, ctx, id.GroupId)
	if err != nil {
		return fmt.Errorf("Error retrieving Azure AD Group members (groupObjectId: %q): %+v", id.GroupId, err)
	}

	var memberObjectID string
	for _, objectID := range members {
		if objectID == id.MemberId {
			memberObjectID = objectID
		}
	}

	if memberObjectID == "" {
		d.SetId("")
		return nil
	}

	d.Set("group_object_id", id.GroupId)
	d.Set("member_object_id", memberObjectID)

	return nil
}

func resourceGroupMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).groupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := graph.ParseGroupMemberId(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to parse ID: %v", err)
	}

	tf.LockByName(groupMemberResourceName, id.GroupId)
	defer tf.UnlockByName(groupMemberResourceName, id.GroupId)

	resp, err := client.RemoveMember(ctx, id.GroupId, id.MemberId)
	if err != nil {
		if !ar.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error removing Member (memberObjectId: %q) from Azure AD Group (groupObjectId: %q): %+v", id.MemberId, id.GroupId, err)
		}
	}

	return nil
}
