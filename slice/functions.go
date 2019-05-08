package slice

// Difference between two slices.
// This function compares elements of slice `a` against one or more other arrays and returns the
// values in `a` that are not present in any of the other arrays.
// ex) Difference([A,B,C],[A]) will return [B,C]
func Difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}
