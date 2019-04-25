package logger

func merge(parts ...map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for _, part := range parts {
		for k, v := range part {
			m[k] = v
		}
	}

	return m
}
