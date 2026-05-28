package maputil

func ToLoMap[T, R any](conv func(T) R) func(T, int) R {
	return func(item T, _ int) R { return conv(item) }
}
