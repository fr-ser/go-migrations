package internal

// StrSliceEqual compares two string slices for equality
func StrSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx++ {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true
}
