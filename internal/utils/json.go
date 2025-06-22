package utils

func ArrayToJSON(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}
	
	json := `["`
	for i, item := range arr {
		if i > 0 {
			json += `","`
		}
		json += item
	}
	json += `"]`
	return json
}