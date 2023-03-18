package di

import (
	"reflect"
)

var container = newContainer()

func GetContainer() *Container {
	return container
}

func CreateContainer() *Container {
	return newContainer()
}

func AddService(descriptor ServiceDescriptor) *Container {
	container.Add(descriptor)
	return container
}

func TryAddService(descriptor ServiceDescriptor) *Container {
	container.TryAdd(descriptor)
	return container
}

func Configure[T any](f func(container *Container, v *T)) {
	container.Configure(getServiceType[T](), func(container *Container, instance any) {
		f(container, instance.(*T))
	})
}

func Add[T any](scope ServiceScope, creator func(*Container) *T) {
	AddService(ServiceDescriptor{
		ServiceType: getServiceType[T](),
		Creator:     convertCreator(creator),
		Scope:       scope,
	})
}

func TryAdd[T any](scope ServiceScope, creator func(*Container) *T) {
	TryAddService(ServiceDescriptor{
		ServiceType: getServiceType[T](),
		Creator:     convertCreator(creator),
		Scope:       scope,
	})
}

func AddTransient[T any](creator func(*Container) *T) {
	Add(Transient, creator)
}

// 添加Transient无参构造对象
func AddTransientDefault[T any]() {
	Add(Transient, func(c *Container) *T {
		return new(T)
	})
}

func TryAddTransient[T any](creator func(*Container) *T) {
	TryAdd(Transient, creator)
}

// 添加Transient无参构造对象
func TryAddTransientDefault[T any]() {
	TryAdd(Transient, func(c *Container) *T {
		return new(T)
	})
}

func AddSingleton[T any](creator func(*Container) *T) {
	Add(Singleton, creator)
}

// 添加 Singleton 无参构造对象
func AddSingletonDefault[T any]() {
	Add(Singleton, func(c *Container) *T {
		return new(T)
	})
}

func TryAddSingleton[T any](creator func(*Container) *T) {
	TryAdd(Singleton, creator)
}

// 添加 Singleton 无参构造对象
func TryAddSingletonDefault[T any]() {
	TryAdd(Singleton, func(c *Container) *T {
		return new(T)
	})
}

func AddValue(value any) {
	AddService(ServiceDescriptor{
		ServiceType: getInterfaceType(value),
		Value:       value,
		Scope:       Singleton,
	})
}

func TryAddValue(value any) {
	TryAddService(ServiceDescriptor{
		ServiceType: getInterfaceType(value),
		Value:       value,
		Scope:       Singleton,
	})
}

func AddByType[T any](serviceType reflect.Type, scope ServiceScope, creator func(*Container) *T) {
	AddService(ServiceDescriptor{
		ServiceType: serviceType,
		Creator:     func(c *Container) any { return creator(c) },
		Scope:       scope,
	})
}

func TryAddByType[T any](serviceType reflect.Type, scope ServiceScope, creator func(*Container) *T) {
	TryAddService(ServiceDescriptor{
		ServiceType: serviceType,
		Creator:     func(c *Container) any { return creator(c) },
		Scope:       scope,
	})
}

func Get[T any](p *Provider) *T {
	return p.Get(getServiceType[T]()).(*T)
}

func GetInterface[T any](p *Provider) T {
	return *p.Get(getServiceType[T]()).(*T)
}

func GetOrNil[T any](p *Provider) (*T, bool) {
	v, ok := p.GetOrNil(getServiceType[T]())
	if !ok {
		return nil, ok
	}
	return v.(*T), true
}

func GetByType[T any](p *Provider, serviceType reflect.Type) *T {
	return p.Get(serviceType).(*T)
}

func GetArray[T any](p *Provider) []T {
	items := p.GetArray(getServiceType[T]())
	typeItems := make([]T, len(items))
	for i := 0; i < len(items); i++ {
		typeItems[i] = items[i].(T)
	}
	return typeItems
}

func getServiceType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func getInterfaceType(value any) reflect.Type {
	return reflect.TypeOf(value).Elem()
}

func convertCreator[T any](creator func(*Container) *T) ServiceCreator {
	if creator == nil {
		return nil
	}
	return func(c *Container) any { return creator(c) }
}
