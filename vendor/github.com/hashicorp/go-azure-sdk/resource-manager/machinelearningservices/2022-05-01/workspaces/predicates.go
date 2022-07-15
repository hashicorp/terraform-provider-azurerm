package workspaces

type AmlUserFeatureOperationPredicate struct {
	Description *string
	DisplayName *string
	Id          *string
}

func (p AmlUserFeatureOperationPredicate) Matches(input AmlUserFeature) bool {

	if p.Description != nil && (input.Description == nil && *p.Description != *input.Description) {
		return false
	}

	if p.DisplayName != nil && (input.DisplayName == nil && *p.DisplayName != *input.DisplayName) {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	return true
}

type WorkspaceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p WorkspaceOperationPredicate) Matches(input Workspace) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil && *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}
