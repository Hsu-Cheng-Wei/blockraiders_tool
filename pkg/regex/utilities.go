package regex

func NamespaceRegex(base string, dirs []string) string {
	var result = base
	for _, s := range dirs {
		if len(s) == 0 {
			continue
		}
		result = result + "." + s
	}
	return result
}
