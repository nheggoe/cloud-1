package fp

func FoldLeft[A, B any](as []A, acc B, f func(A, B) B) B {
	if len(as) == 0 {
		return acc
	}
	return FoldLeft[A, B](as[1:], f(as[0], acc), f)
}

func ForAll[A any](as []A, predicate func(A) bool) bool {
	return FoldLeft(as, true, func(a A, b bool) bool { return predicate(a) && b })
}
