package databaselink

func ForceLinkNeeded(oldItemList []interface{}, newItemList []interface{}) bool {
	oldItems := make(map[string]bool)
	for _, oldItem := range oldItemList {
		oldItems[oldItem.(string)] = true
	}
	for _, newItem := range newItemList {
		if !oldItems[newItem.(string)] {
			return true
		}
	}
	return false
}
