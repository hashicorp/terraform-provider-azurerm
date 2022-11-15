package configurationstores

type ApiKeyOperationPredicate struct {
	ConnectionString *string
	Id               *string
	LastModified     *string
	Name             *string
	ReadOnly         *bool
	Value            *string
}

func (p ApiKeyOperationPredicate) Matches(input ApiKey) bool {

	if p.ConnectionString != nil && (input.ConnectionString == nil && *p.ConnectionString != *input.ConnectionString) {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.LastModified != nil && (input.LastModified == nil && *p.LastModified != *input.LastModified) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.ReadOnly != nil && (input.ReadOnly == nil && *p.ReadOnly != *input.ReadOnly) {
		return false
	}

	if p.Value != nil && (input.Value == nil && *p.Value != *input.Value) {
		return false
	}

	return true
}

type ConfigurationStoreOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ConfigurationStoreOperationPredicate) Matches(input ConfigurationStore) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
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
