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
	ToExpandedConfig() ExpandedConfig
	FromExpandedConfig(ExpandedConfig)
}

var _ ExpandedConfigCaster = &SystemAssignedIdentity{}

type SystemAssignedIdentity struct {
	Type        Type    `json:"type,omitempty"`
	TenantId    *string `json:"tenantId,omitempty"`
	PrincipalId *string `json:"principalId,omitempty"`
}

func (s *SystemAssignedIdentity) ToExpandedConfig() ExpandedConfig {
	if s == nil {
		return ExpandedConfig{}
	}
	principalId := ""
	if s.PrincipalId != nil {
		principalId = *s.PrincipalId
	}
	tenantId := ""
	if s.TenantId != nil {
		tenantId = *s.TenantId
	}
	return ExpandedConfig{
		Type:        s.Type,
		PrincipalId: principalId,
		TenantId:    tenantId,
	}
}

func (s *SystemAssignedIdentity) FromExpandedConfig(config ExpandedConfig) {
	if s == nil {
		return
	}
	*s = SystemAssignedIdentity{
		Type:        config.Type,
		TenantId:    &config.TenantId,
		PrincipalId: &config.PrincipalId,
	}
}

var _ ExpandedConfigCaster = &UserAssignedIdentityList{}

type UserAssignedIdentityList struct {
	Type                   Type                    `json:"type,omitempty"`
	UserAssignedIdentities *[]userAssignedIdentity `json:"userAssignedIdentities"`
}

func (u *UserAssignedIdentityList) ToExpandedConfig() ExpandedConfig {
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
	out.UserAssignedIdentityIds = identities

	return out
}

func (u *UserAssignedIdentityList) FromExpandedConfig(config ExpandedConfig) {
	if u == nil {
		return
	}
	*u = UserAssignedIdentityList{
		Type: config.Type,
	}

	if len(config.UserAssignedIdentityIds) == 0 {
		return
	}

	var identities []userAssignedIdentity
	for _, id := range config.UserAssignedIdentityIds {
		identities = append(identities, userAssignedIdentity{
			ResourceId: &id,
		})
	}
	u.UserAssignedIdentities = &identities
}

var _ ExpandedConfigCaster = &UserAssignedIdentityMap{}

type UserAssignedIdentityMap struct {
	Type                   Type                                 `json:"type,omitempty"`
	UserAssignedIdentities map[string]*userAssignedIdentityInfo `json:"userAssignedIdentities"`
}

func (u *UserAssignedIdentityMap) ToExpandedConfig() ExpandedConfig {
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

	out.UserAssignedIdentityIds = identities

	return out
}

func (u *UserAssignedIdentityMap) FromExpandedConfig(config ExpandedConfig) {
	if u == nil {
		return
	}

	*u = UserAssignedIdentityMap{
		Type: config.Type,
	}

	if len(config.UserAssignedIdentityIds) == 0 {
		return
	}

	u.UserAssignedIdentities = map[string]*userAssignedIdentityInfo{}
	for _, id := range config.UserAssignedIdentityIds {
		// The user assigned identity information is not used by the provider. So simply assign the value to nil.
		u.UserAssignedIdentities[id] = &userAssignedIdentityInfo{ClientId: nil, PrincipalId: nil}
	}
}

var _ ExpandedConfigCaster = &SystemUserAssignedIdentityList{}

type SystemUserAssignedIdentityList struct {
	Type                   Type                    `json:"type,omitempty"`
	TenantId               *string                 `json:"tenantId,omitempty"`
	PrincipalId            *string                 `json:"principalId,omitempty"`
	UserAssignedIdentities *[]userAssignedIdentity `json:"userAssignedIdentities"`
}

func (s *SystemUserAssignedIdentityList) ToExpandedConfig() ExpandedConfig {
	if s == nil {
		return ExpandedConfig{}
	}
	principalId := ""
	if s.PrincipalId != nil {
		principalId = *s.PrincipalId
	}
	tenantId := ""
	if s.TenantId != nil {
		tenantId = *s.TenantId
	}

	out := ExpandedConfig{
		Type:        s.Type,
		PrincipalId: principalId,
		TenantId:    tenantId,
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
	out.UserAssignedIdentityIds = identities

	return out
}

func (s *SystemUserAssignedIdentityList) FromExpandedConfig(config ExpandedConfig) {
	if s == nil {
		return
	}

	*s = SystemUserAssignedIdentityList{
		Type:        config.Type,
		TenantId:    &config.TenantId,
		PrincipalId: &config.PrincipalId,
	}

	if len(config.UserAssignedIdentityIds) == 0 {
		return
	}

	var identities []userAssignedIdentity
	for _, id := range config.UserAssignedIdentityIds {
		identities = append(identities, userAssignedIdentity{
			ResourceId: &id,
		})
	}
	s.UserAssignedIdentities = &identities
}

var _ ExpandedConfigCaster = &SystemUserAssignedIdentityMap{}

type SystemUserAssignedIdentityMap struct {
	Type                   Type                                 `json:"type,omitempty"`
	TenantId               *string                              `json:"tenantId,omitempty"`
	PrincipalId            *string                              `json:"principalId,omitempty"`
	UserAssignedIdentities map[string]*userAssignedIdentityInfo `json:"userAssignedIdentities"`
}

func (s *SystemUserAssignedIdentityMap) ToExpandedConfig() ExpandedConfig {
	if s == nil {
		return ExpandedConfig{}
	}
	principalId := ""
	if s.PrincipalId != nil {
		principalId = *s.PrincipalId
	}
	tenantId := ""
	if s.TenantId != nil {
		tenantId = *s.TenantId
	}

	out := ExpandedConfig{
		Type:        s.Type,
		PrincipalId: principalId,
		TenantId:    tenantId,
	}

	var identities []string
	for k := range s.UserAssignedIdentities {
		identities = append(identities, k)
	}
	out.UserAssignedIdentityIds = identities

	return out
}

func (s *SystemUserAssignedIdentityMap) FromExpandedConfig(config ExpandedConfig) {
	if s == nil {
		return
	}

	*s = SystemUserAssignedIdentityMap{
		Type:        config.Type,
		TenantId:    &config.TenantId,
		PrincipalId: &config.PrincipalId,
	}

	if len(config.UserAssignedIdentityIds) == 0 {
		return
	}

	s.UserAssignedIdentities = map[string]*userAssignedIdentityInfo{}
	for _, id := range config.UserAssignedIdentityIds {
		// The user assigned identity information is not used by the provider. So simply assign the value to nil.
		s.UserAssignedIdentities[id] = &userAssignedIdentityInfo{ClientId: nil, PrincipalId: nil}
	}
}
