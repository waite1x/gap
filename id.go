package gap

type IdGenerator[T comparable] interface {
	Create() T
}
