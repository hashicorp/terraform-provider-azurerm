package graph

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
)

type GroupMemberId struct {
	ObjectSubResourceId
	GroupId  string
	MemberId string
}

func GroupMemberIdFrom(groupId, memberId string) GroupMemberId {
	return GroupMemberId{
		ObjectSubResourceId: ObjectSubResourceIdFrom(groupId, "member", memberId),
		GroupId:             groupId,
		MemberId:            memberId,
	}
}

func ParseGroupMemberId(idString string) (GroupMemberId, error) {
	id, err := ParseObjectSubResourceId(idString, "member")
	if err != nil {
		return GroupMemberId{}, fmt.Errorf("Unable to parse Member ID: %v", err)
	}

	return GroupMemberId{
		ObjectSubResourceId: id,
		GroupId:             id.objectId,
		MemberId:            id.subId,
	}, nil
}

func GroupGetByDisplayName(client *graphrbac.GroupsClient, ctx context.Context, displayName string) (*graphrbac.ADGroup, error) {

	filter := fmt.Sprintf("displayName eq '%s'", displayName)

	resp, err := client.ListComplete(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Error listing Azure AD Groups for filter %q: %+v", filter, err)
	}

	values := resp.Response().Value
	if values == nil {
		return nil, fmt.Errorf("nil values for AD Groups matching %q", filter)
	}
	if len(*values) == 0 {
		return nil, fmt.Errorf("Found no AD Groups matching %q", filter)
	}
	if len(*values) > 2 {
		return nil, fmt.Errorf("Found multiple AD Groups matching %q", filter)
	}

	group := (*values)[0]
	if group.DisplayName == nil {
		return nil, fmt.Errorf("nil DisplayName for AD Groups matching %q", filter)
	}
	if *group.DisplayName != displayName {
		return nil, fmt.Errorf("displayname for AD Groups matching %q does is does not match(%q!=%q)", filter, *group.DisplayName, displayName)
	}

	return &group, nil
}

func DirectoryObjectListToIDs(objects graphrbac.DirectoryObjectListResultIterator, ctx context.Context) ([]string, error) {
	ids := make([]string, 0)
	for objects.NotDone() {
		v := objects.Value()

		// possible members are users, groups or service principals
		// we try to 'cast' each result as the corresponding type and diff
		// if we found the object we're looking for
		user, _ := v.AsUser()
		if user != nil {
			ids = append(ids, *user.ObjectID)
		}

		group, _ := v.AsADGroup()
		if group != nil {
			ids = append(ids, *group.ObjectID)
		}

		servicePrincipal, _ := v.AsServicePrincipal()
		if servicePrincipal != nil {
			ids = append(ids, *servicePrincipal.ObjectID)
		}

		if err := objects.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("Error during pagination of directory objects: %+v", err)
		}
	}

	return ids, nil
}

func GroupAllMembers(client graphrbac.GroupsClient, ctx context.Context, groupId string) ([]string, error) {
	members, err := client.GetGroupMembersComplete(ctx, groupId)

	if err != nil {
		return nil, fmt.Errorf("Error listing existing group members from Azure AD Group with ID %q: %+v", groupId, err)
	}

	existingMembers, err := DirectoryObjectListToIDs(members, ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting objects IDs of group members for Azure AD Group with ID %q: %+v", groupId, err)
	}

	log.Printf("[DEBUG] %d members in Azure AD group with ID: %q", len(existingMembers), groupId)

	return existingMembers, nil
}

func GroupAddMember(client graphrbac.GroupsClient, ctx context.Context, groupId string, member string) error {
	memberGraphURL := fmt.Sprintf("https://graph.windows.net/%s/directoryObjects/%s", client.TenantID, member)

	properties := graphrbac.GroupAddMemberParameters{
		URL: &memberGraphURL,
	}

	log.Printf("[DEBUG] Adding member with id %q to Azure AD group with id %q", member, groupId)
	if _, err := client.AddMember(ctx, groupId, properties); err != nil {
		return fmt.Errorf("Error adding group member %q to Azure AD Group with ID %q: %+v", member, groupId, err)
	}

	return nil
}

func GroupAddMembers(client graphrbac.GroupsClient, ctx context.Context, groupId string, members []string) error {
	for _, memberUuid := range members {
		err := GroupAddMember(client, ctx, groupId, memberUuid)

		if err != nil {
			return fmt.Errorf("Error while adding members to Azure AD Group with ID %q: %+v", groupId, err)
		}
	}

	return nil
}

func GroupAllOwners(client graphrbac.GroupsClient, ctx context.Context, groupId string) ([]string, error) {
	owners, err := client.ListOwnersComplete(ctx, groupId)

	if err != nil {
		return nil, fmt.Errorf("Error listing existing group owners from Azure AD Group with ID %q: %+v", groupId, err)
	}

	existingMembers, err := DirectoryObjectListToIDs(owners, ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting objects IDs of group owners for Azure AD Group with ID %q: %+v", groupId, err)
	}

	log.Printf("[DEBUG] %d members in Azure AD group with ID: %q", len(existingMembers), groupId)
	return existingMembers, nil
}

func GroupAddOwner(client graphrbac.GroupsClient, ctx context.Context, groupId string, owner string) error {
	ownerGraphURL := fmt.Sprintf("https://graph.windows.net/%s/directoryObjects/%s", client.TenantID, owner)

	properties := graphrbac.AddOwnerParameters{
		URL: &ownerGraphURL,
	}

	log.Printf("[DEBUG] Adding owner with id %q to Azure AD group with id %q", owner, groupId)
	if _, err := client.AddOwner(ctx, groupId, properties); err != nil {
		return fmt.Errorf("Error adding group owner %q to Azure AD Group with ID %q: %+v", owner, groupId, err)
	}

	return nil
}

func GroupAddOwners(client graphrbac.GroupsClient, ctx context.Context, groupId string, owner []string) error {
	for _, ownerUuid := range owner {
		err := GroupAddOwner(client, ctx, groupId, ownerUuid)

		if err != nil {
			return fmt.Errorf("Error while adding owners to Azure AD Group with ID %q: %+v", groupId, err)
		}
	}

	return nil
}
