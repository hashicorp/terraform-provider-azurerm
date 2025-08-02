package databaselink

func ForceUnlinkItems(oldItemList []interface{}, newItemList []interface{}) (bool, []string) {
	newItems := make(map[string]bool)
	forceUnlinkList := make([]string, 0)
	for _, newItem := range newItemList {
		newItems[newItem.(string)] = true
	}

	for _, oldItem := range oldItemList {
		if !newItems[oldItem.(string)] {
			forceUnlinkList = append(forceUnlinkList, oldItem.(string))
		}
	}
	if len(forceUnlinkList) > 0 {
		return true, forceUnlinkList
	}
	return false, nil
}
