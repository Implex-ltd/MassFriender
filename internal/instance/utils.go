package instance

func containsTask(slice []string, task string) bool {
	for _, t := range slice {
		if t == task {
			return true
		}
	}
	return false
}
