package mapslice

func Map2Slice[T comparable, V any](m map[T]V) ([]T, []V) {
	sk := make([]T, 0, len(m))
	sv := make([]V, 0, len(m))
	for k, v := range m {
		sk = append(sk, k)
		sv = append(sv, v)
	}
	return sk, sv
}

func Slice2Map[T comparable](s []T) map[T]struct{} {
	m := make(map[T]struct{}, len(s))
	for _, t := range s {
		m[t] = struct{}{}
	}
	return m
}
func MapKey2Slice[T comparable, V any](m map[T]V) []T {
	s := make([]T, 0, len(m))
	for key := range m {
		s = append(s, key)
	}
	return s
}

func MapValue2Slice[T comparable, V any](m map[T]V) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
