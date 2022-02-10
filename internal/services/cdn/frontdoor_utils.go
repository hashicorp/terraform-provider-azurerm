package cdn

func ConvertFrontdoorProfileTags(tagMap *map[string]string) map[string]*string {
	t := make(map[string]*string)
	if tagMap != nil {
		for k, v := range *tagMap {
			tagKey := k
			tagValue := v
			t[tagKey] = &tagValue
		}
	}

	return t
}
