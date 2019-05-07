package array

/** It returns the elements in the first slice that are not present in the second one.
*	ex) Difference([A,B,C],[A]) will return [B,C]
 */
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
