package identity

type userAssignedIdentity struct {
	ResourceId *string `json:"resourceId,omitempty"`
	userAssignedIdentityInfo
}

type userAssignedIdentityInfo struct {
	ClientId    *string `json:"clientId,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
}

type ExpandedConfigCaster interface {
	CastToExpandedConfig() ExpandedConfig
	CastFromExpandedConfig(ExpandedConfig)
}

var _ ExpandedConfigCaster = &SystemAssignedIdentity{}

type SystemAssignedIdentity struct {
	Type        string  `json:"type,omitempty"`
	TenantID    *string `json:"tenantId,omitempty"`
	PrincipalID *string `json:"principalId,omitempty"`
}

func (s *SystemAssignedIdentity) CastToExpandedConfig() ExpandedConfig {
	if s == nil {
		return ExpandedConfig{}
	}
	return ExpandedConfig{
		Type:        s.Type,
		PrincipalId: s.PrincipalID,
		TenantId:    s.TenantID,
	}
}

func (s *SystemAssignedIdentity) CastFromExpandedConfig(config ExpandedConfig) {
	if s == nil {
		return
	}
	*s = SystemAssignedIdentity{
		Type:        config.Type,
		TenantID:    config.TenantId,
		PrincipalID: config.PrincipalId,
	}

	return
}

var _ ExpandedConfigCaster = &UserAssignedIdentityList{}

type UserAssignedIdentityList struct {
	Type                   string                  `json:"type,omitempty"`
	UserAssignedIdentities *[]userAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}

func (u *UserAssignedIdentityList) CastToExpandedConfig() ExpandedConfig {
	if u == nil {
		return ExpandedConfig{}
	}
	out := ExpandedConfig{
		Type: u.Type,
	}

	if u.UserAssignedIdentities == nil {
		return out
	}

	var identities []string
	for _, id := range *u.UserAssignedIdentities {
		if id.ResourceId == nil {
			continue
		}
		identities = append(identities, *id.PrincipalId)
	}
	out.UserAssignedIdentityIds = &identities

	return out
}

func (u *UserAssignedIdentityList) CastFromExpandedConfig(config ExpandedConfig) {
	if u == nil {
		return
	}
	*u = UserAssignedIdentityList{
		Type: config.Type,
	}

	if config.UserAssignedIdentityIds == nil {
		return
	}

	var identities []userAssignedIdentity
	for _, id := range *config.UserAssignedIdentityIds {
		identities = append(identities, userAssignedIdentity{
			ResourceId: &id,
		})
	}
	u.UserAssignedIdentities = &identities

	return
}

var _ ExpandedConfigCaster = &UserAssignedIdentityMap{}

type UserAssignedIdentityMap struct {
	Type                   string                               `json:"type,omitempty"`
	UserAssignedIdentities map[string]*userAssignedIdentityInfo `json:"userAssignedIdentities,omitempty"`
}

func (u *UserAssignedIdentityMap) CastToExpandedConfig() ExpandedConfig {
	if u == nil {
		return ExpandedConfig{}
	}
	out := ExpandedConfig{
		Type: u.Type,
	}

	var identities []string
	for k := range u.UserAssignedIdentities {
		identities = append(identities, k)
	}

	if len(identities) > 0 {
		out.UserAssignedIdentityIds = &identities
	}

	return out
}

func (u *UserAssignedIdentityMap) CastFromExpandedConfig(config ExpandedConfig) {
	if u == nil {
		return
	}

	*u = UserAssignedIdentityMap{
		Type: config.Type,
	}

	if config.UserAssignedIdentityIds == nil {
		return
	}

	u.UserAssignedIdentities = map[string]*userAssignedIdentityInfo{}
	for _, id := range *config.UserAssignedIdentityIds {
		// The user assigned identity information is not used by the provider. So simply assign the value to nil.
		u.UserAssignedIdentities[id] = nil
	}

	return
}

var _ ExpandedConfigCaster = &SystemUserAssignedIdentityList{}

type SystemUserAssignedIdentityList struct {
	Type                   string                  `json:"type,omitempty"`
	TenantId               *string                 `json:"tenantId,omitempty"`
	PrincipalId            *string                 `json:"principalId,omitempty"`
	UserAssignedIdentities *[]userAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}

func (s *SystemUserAssignedIdentityList) CastToExpandedConfig() ExpandedConfig {
	if s == nil {
		return ExpandedConfig{}
	}

	out := ExpandedConfig{
		Type:        s.Type,
		PrincipalId: s.PrincipalId,
		TenantId:    s.TenantId,
	}

	if s.UserAssignedIdentities == nil {
		return out
	}

	var identities []string
	for _, id := range *s.UserAssignedIdentities {
		if id.ResourceId == nil {
			continue
		}
		identities = append(identities, *id.ResourceId)
	}
	out.UserAssignedIdentityIds = &identities

	return out
}

func (s *SystemUserAssignedIdentityList) CastFromExpandedConfig(config ExpandedConfig) {
	if s == nil {
		return
	}

	*s = SystemUserAssignedIdentityList{
		Type:        config.Type,
		TenantId:    config.TenantId,
		PrincipalId: config.PrincipalId,
	}

	if config.UserAssignedIdentityIds == nil {
		return
	}

	var identities []userAssignedIdentity
	for _, id := range *config.UserAssignedIdentityIds {
		identities = append(identities, userAssignedIdentity{
			ResourceId: &id,
		})
	}
	s.UserAssignedIdentities = &identities

	return
}

var _ ExpandedConfigCaster = &SystemUserAssignedIdentityMap{}

type SystemUserAssignedIdentityMap struct {
	Type                   string                               `json:"type,omitempty"`
	TenantId               *string                              `json:"tenantId,omitempty"`
	PrincipalId            *string                              `json:"principalId,omitempty"`
	UserAssignedIdentities map[string]*userAssignedIdentityInfo `json:"userAssignedIdentities,omitempty"`
}

func (s *SystemUserAssignedIdentityMap) CastToExpandedConfig() ExpandedConfig {
	if s == nil {
		return ExpandedConfig{}
	}

	out := ExpandedConfig{
		Type:        s.Type,
		PrincipalId: s.PrincipalId,
		TenantId:    s.TenantId,
	}

	var identities []string
	for k := range s.UserAssignedIdentities {
		identities = append(identities, k)
	}
	if len(identities) > 0 {
		out.UserAssignedIdentityIds = &identities
	}

	return out
}

func (s *SystemUserAssignedIdentityMap) CastFromExpandedConfig(config ExpandedConfig) {
	if s == nil {
		return
	}

	*s = SystemUserAssignedIdentityMap{
		Type:        config.Type,
		TenantId:    config.TenantId,
		PrincipalId: config.PrincipalId,
	}

	if config.UserAssignedIdentityIds == nil {
		return
	}

	s.UserAssignedIdentities = map[string]*userAssignedIdentityInfo{}
	for _, id := range *config.UserAssignedIdentityIds {
		// The user assigned identity information is not used by the provider. So simply assign the value to nil.
		s.UserAssignedIdentities[id] = nil
	}

	return
}
