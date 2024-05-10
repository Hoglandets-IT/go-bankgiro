package tools

func SliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func SliceEqual(slices ...[]string) bool {
	for i, slice := range slices {
		for j, otherSlice := range slices {
			if i != j {
				if len(slice) != len(otherSlice) {
					return false
				}
				for k, value := range slice {
					if value != otherSlice[k] {
						return false
					}
				}
			}
		}
	}

	return true
}

func BytesliceEqual(slices ...[]byte) bool {
	for i, slice := range slices {
		for j, otherSlice := range slices {
			if i != j {
				if len(slice) != len(otherSlice) {
					return false
				}
				for k, value := range slice {
					if value != otherSlice[k] {
						return false
					}
				}
			}
		}
	}

	return true
}
