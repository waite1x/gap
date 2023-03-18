package di

import (
	"context"
	"reflect"
	"sync"
)

func NewContainer() *Container {
	return &Container{
		initors:  sync.Map{},
		Services: make([]ServiceDescriptor, 0),
		values:   sync.Map{},
	}
}

type ConfigureFunc func(container *Container, instance any)

type Container struct {
	initors  sync.Map
	values   sync.Map
	Services []ServiceDescriptor
}

func (c *Container) Add(descriptor ServiceDescriptor) {
	c.Services = append(c.Services, descriptor)
}

func (c *Container) TryAdd(descriptor ServiceDescriptor) {
	index := c.findIndex(&descriptor.ServiceType)
	if index < 0 {
		c.Services = append(c.Services, descriptor)
	}
}

func (c *Container) Configure(t reflect.Type, initFunc ConfigureFunc) {
	v, ok := c.initors.Load(t)
	if ok {
		c.initors.Store(t, append(v.([]ConfigureFunc), initFunc))
	} else {
		inits := []ConfigureFunc{initFunc}
		c.initors.Store(t, inits)
	}
}

func (c *Container) FirstOrDefault(serviceType reflect.Type) *ServiceDescriptor {
	for i := 0; i < len(c.Services); i++ {
		v := &c.Services[i]
		if v.ServiceType == serviceType {
			return v
		}
	}
	return nil
}

func (c *Container) CreateScope(ctx context.Context) *Provider {
	return NewProvider(ctx, c)
}

func (c *Container) findIndex(serviceType *reflect.Type) int {
	for i := 0; i < len(c.Services); i++ {
		v := &c.Services[i]
		if v.ServiceType == *serviceType {
			return i
		}
	}
	return -1
}
