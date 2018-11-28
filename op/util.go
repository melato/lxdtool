package op

func StringSliceDiff(ar []string, exclude []string) []string {
	if exclude == nil {
		return ar
	}
	var xmap = make(map[string]bool)
	for _, s := range exclude {
		xmap[s] = true
	}
	var result []string
	for _, s := range ar {
		if !xmap[s] {
			result = append(result, s)
		}
	}
	return result
}
