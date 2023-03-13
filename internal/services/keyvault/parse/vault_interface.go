package parse

import "fmt"

type VaultType string

const (
	VaultTypeDefault VaultType = "vault"
	VaultTypeMHSM    VaultType = "managedhsm"
)

func IsValidValtType(s string) bool {
	return s == string(VaultTypeDefault) || s == string(VaultTypeMHSM)
}

type Vaulter interface {
	ID() string
	GetSubscriptionID() string
	GetResourceGroup() string
	GetName() string
	GetCacheKey() string
	Type() VaultType
}

var (
	_ Vaulter = VaultId{}
	_ Vaulter = ManagedHSMId{}
)

func IsMHSMVaulter(v Vaulter) bool {
	_, ok := v.(*ManagedHSMId)
	return ok
}

func NewVaulterFromString(input string) (Vaulter, error) {
	var e1, e2 error
	if vid, e1 := VaultID(input); e1 == nil {
		return vid, nil
	}
	if mid, e2 := ManagedHSMID(input); e2 == nil {
		return mid, nil
	}
	return nil, fmt.Errorf("parse vautler err: %+v or +%v", e1, e2)
}

func (id VaultId) GetSubscriptionID() string {
	return id.SubscriptionId
}

func (id VaultId) GetResourceGroup() string {
	return id.ResourceGroup
}

func (id VaultId) GetName() string {
	return id.Name
}
func (id VaultId) Type() VaultType {
	return VaultTypeDefault
}

func (id VaultId) GetCacheKey() string {
	return MakeCacheKey(id.Type(), id.GetName())
}

func (id ManagedHSMId) GetSubscriptionID() string {
	return id.SubscriptionId
}

func (id ManagedHSMId) GetResourceGroup() string {
	return id.ResourceGroup
}

func (id ManagedHSMId) GetName() string {
	return id.Name
}

func (id ManagedHSMId) Type() VaultType {
	return VaultTypeMHSM
}

func (id ManagedHSMId) GetCacheKey() string {
	return MakeCacheKey(id.Type(), id.GetName())
}

func MakeCacheKey(typ VaultType, name string) string {
	return fmt.Sprintf("%s:%s", typ, name)
}
