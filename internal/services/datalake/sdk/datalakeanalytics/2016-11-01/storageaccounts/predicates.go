package storageaccounts

type SasTokenInformationPredicate struct {
	AccessToken *string
}

func (p SasTokenInformationPredicate) Matches(input SasTokenInformation) bool {

	if p.AccessToken != nil && (input.AccessToken == nil && *p.AccessToken != *input.AccessToken) {
		return false
	}

	return true
}

type StorageAccountInformationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p StorageAccountInformationPredicate) Matches(input StorageAccountInformation) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

type StorageContainerPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p StorageContainerPredicate) Matches(input StorageContainer) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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
