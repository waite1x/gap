package di

import (
	"context"
	"fmt"
	"reflect"
)

type Provider struct {
	Ctx context.Context
	c   *Container
}

func NewProvider(ctx context.Context, c *Container) *Provider {
	return &Provider{
		c:   c,
		Ctx: ctx,
	}
}

func (p *Provider) Get(serviceType reflect.Type) interface{} {
	v, ok := p.GetOrNil(serviceType)
	if !ok {
		panic(fmt.Sprintf("service %s not found", serviceType))
	}
	return v
}

// 获取对象，如果没有，返回nil
func (p *Provider) GetOrNil(serviceType reflect.Type) (interface{}, bool) {
	descriptor := p.c.FirstOrDefault(serviceType)
	if descriptor == nil {
		return nil, false
	}
	return p.create(serviceType, descriptor), true
}

func (p *Provider) GetArray(baseType reflect.Type) []interface{} {
	instances := make([]interface{}, 0)
	for i := 0; i < len(p.c.Services); i++ {
		v := &p.c.Services[i]
		implType := reflect.PtrTo(v.ServiceType)
		if v.ServiceType == baseType || implType == baseType || implType.AssignableTo(baseType) {
			instances = append(instances, p.create(v.ServiceType, v))
		}
	}
	return instances
}

func (p *Provider) create(serviceType reflect.Type, descriptor *ServiceDescriptor) interface{} {
	var instance any
	if descriptor.Scope == Singleton {
		v, ok := p.c.values.Load(serviceType)
		if !ok {
			v, _ = p.c.values.LoadOrStore(serviceType, p.createInstance(descriptor))
		}
		instance = v
		p.instanceInit(serviceType, instance, true)
	} else if descriptor.Scope == Transient {
		instance = p.createInstance(descriptor)
		p.instanceInit(serviceType, instance, false)
	} else {
		panic("invalid scope")
	}
	return instance
}

func (p *Provider) createInstance(descriptor *ServiceDescriptor) interface{} {
	if descriptor.Value != nil {
		return descriptor.Value
	}
	if descriptor.Creator != nil {
		return descriptor.Creator(p.c)
	}
	return reflect.New(descriptor.ServiceType).Interface()
}

func (p *Provider) instanceInit(serviceType reflect.Type, v any, delete bool) {
	if delete {
		if initor, ok := p.c.initors.LoadAndDelete(serviceType); ok {
			p.exucteInit(initor, v)
		}
	} else {
		if initor, ok := p.c.initors.Load(serviceType); ok {
			p.exucteInit(initor, v)
		}
	}
}

func (p *Provider) exucteInit(initor any, v any) {
	if init, ok := initor.([]ConfigureFunc); ok {
		for i := 0; i < len(init); i++ {
			init[i](p.c, v)
		}
	}
}
