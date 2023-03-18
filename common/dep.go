package common

import (
	"reflect"

	"github.com/waite1x/gapp/common/di"
)

var container = di.NewContainer()

func GetContainer() *di.Container {
	return container
}

func CreateContainer() *di.Container {
	return di.NewContainer()
}

func AddService(descriptor di.ServiceDescriptor) *di.Container {
	container.Add(descriptor)
	return container
}

func TryAddService(descriptor di.ServiceDescriptor) *di.Container {
	container.TryAdd(descriptor)
	return container
}

func Configure[T any](f func(container *di.Container, v *T)) {
	container.Configure(getServiceType[T](), func(container *di.Container, instance any) {
		f(container, instance.(*T))
	})
}

func Add[T any](scope di.ServiceScope, creator func(*di.Container) *T) {
	AddService(di.ServiceDescriptor{
		ServiceType: getServiceType[T](),
		Creator:     convertCreator(creator),
		Scope:       scope,
	})
}

func TryAdd[T any](scope di.ServiceScope, creator func(*di.Container) *T) {
	TryAddService(di.ServiceDescriptor{
		ServiceType: getServiceType[T](),
		Creator:     convertCreator(creator),
		Scope:       scope,
	})
}

func AddTransient[T any](creator func(*di.Container) *T) {
	Add(di.Transient, creator)
}

// 添加Transient无参构造对象
func AddTransientDefault[T any]() {
	Add(di.Transient, func(c *di.Container) *T {
		return new(T)
	})
}

func TryAddTransient[T any](creator func(*di.Container) *T) {
	TryAdd(di.Transient, creator)
}

// 添加Transient无参构造对象
func TryAddTransientDefault[T any]() {
	TryAdd(di.Transient, func(c *di.Container) *T {
		return new(T)
	})
}

func AddSingleton[T any](creator func(*di.Container) *T) {
	Add(di.Singleton, creator)
}

// 添加 Singleton 无参构造对象
func AddSingletonDefault[T any]() {
	Add(di.Singleton, func(c *di.Container) *T {
		return new(T)
	})
}

func TryAddSingleton[T any](creator func(*di.Container) *T) {
	TryAdd(di.Singleton, creator)
}

// 添加 Singleton 无参构造对象
func TryAddSingletonDefault[T any]() {
	TryAdd(di.Singleton, func(c *di.Container) *T {
		return new(T)
	})
}

func AddValue(value any) {
	AddService(di.ServiceDescriptor{
		ServiceType: getInterfaceType(value),
		Value:       value,
		Scope:       di.Singleton,
	})
}

func TryAddValue(value any) {
	TryAddService(di.ServiceDescriptor{
		ServiceType: getInterfaceType(value),
		Value:       value,
		Scope:       di.Singleton,
	})
}

func AddByType[T any](serviceType reflect.Type, scope di.ServiceScope, creator func(*di.Container) *T) {
	AddService(di.ServiceDescriptor{
		ServiceType: serviceType,
		Creator:     func(c *di.Container) any { return creator(c) },
		Scope:       scope,
	})
}

func TryAddByType[T any](serviceType reflect.Type, scope di.ServiceScope, creator func(*di.Container) *T) {
	TryAddService(di.ServiceDescriptor{
		ServiceType: serviceType,
		Creator:     func(c *di.Container) any { return creator(c) },
		Scope:       scope,
	})
}

func Get[T any](p *di.Provider) *T {
	return p.Get(getServiceType[T]()).(*T)
}

func GetInterface[T any](p *di.Provider) T {
	return *p.Get(getServiceType[T]()).(*T)
}

func GetOrNil[T any](p *di.Provider) (*T, bool) {
	v, ok := p.GetOrNil(getServiceType[T]())
	if !ok {
		return nil, ok
	}
	return v.(*T), true
}

func GetByType[T any](p *di.Provider, serviceType reflect.Type) *T {
	return p.Get(serviceType).(*T)
}

func GetArray[T any](p *di.Provider) []T {
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

func convertCreator[T any](creator func(*di.Container) *T) di.ServiceCreator {
	if creator == nil {
		return nil
	}
	return func(c *di.Container) any { return creator(c) }
}
