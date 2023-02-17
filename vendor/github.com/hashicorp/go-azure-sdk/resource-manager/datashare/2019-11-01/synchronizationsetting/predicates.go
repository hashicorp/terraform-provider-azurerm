package synchronizationsetting

type SynchronizationSettingOperationPredicate struct {
}

func (p SynchronizationSettingOperationPredicate) Matches(input SynchronizationSetting) bool {

	return true
}
