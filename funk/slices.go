package funk

func Map[A, B any](slice []A, f func(A) B) []B {
	res := make([]B, len(slice))
	for i, a := range slice {
		res[i] = f(a)
	}

	return res
}
