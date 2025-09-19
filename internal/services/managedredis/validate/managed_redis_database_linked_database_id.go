package validate

import "fmt"

func ValidateLinkedDatabaseIncludesSelf(linkedDatabaseIds []string, selfDbId string) error {
	if len(linkedDatabaseIds) == 0 {
		return nil
	}

	for _, id := range linkedDatabaseIds {
		if id == selfDbId {
			return nil
		}
	}

	return fmt.Errorf("linked_database_id must include this database ID: %s", selfDbId)
}
