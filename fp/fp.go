package fp

func FoldLeft[A any, B any](as []A, acc B, f func(A, B) B) B {
	if as == nil || len(as) == 0 {
		return acc
	}
	return FoldLeft(as[1:], f(as[0], acc), f)
}

func Map[A any, B any](as []A, f func(A) B) []B {
	return FoldLeft[A, []B](as, make([]B, 0), func(a A, acc []B) []B { return append(acc, f(a)) })
}

func Reduce[A any](as []A, f func(A, A) A) A {
	if len(as) == 0 {
		panic("Reduce: empty slice")
	}
	return FoldLeft[A, A](as[1:], as[0], func(a A, acc A) A { return f(acc, a) })
}
