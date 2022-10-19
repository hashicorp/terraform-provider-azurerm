package contact

type AvailableContactsOperationPredicate struct {
	GroundStationName *string
}

func (p AvailableContactsOperationPredicate) Matches(input AvailableContacts) bool {

	if p.GroundStationName != nil && (input.GroundStationName == nil && *p.GroundStationName != *input.GroundStationName) {
		return false
	}

	return true
}

type ContactOperationPredicate struct {
	Etag *string
	Id   *string
	Name *string
	Type *string
}

func (p ContactOperationPredicate) Matches(input Contact) bool {

	if p.Etag != nil && (input.Etag == nil && *p.Etag != *input.Etag) {
		return false
	}

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
