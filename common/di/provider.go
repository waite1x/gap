package di

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/waite1x/gap/common/log"
)

const ProviderKey = "di:Provider"

type Provider struct {
	Ctx    context.Context
	c      *Container
	values sync.Map
}

func NewProvider(ctx context.Context, c *Container) *Provider {
	return &Provider{
		c:   c,
		Ctx: ctx,
	}
}

func CuurentProvider(ctx context.Context) *Provider {
	v := ctx.Value(ProviderKey)
	if v == nil {
		panic("di: provider not found")
	}
	return v.(*Provider)
}

func (p *Provider) Get(serviceType reflect.Type) any {
	v, ok := p.GetOrNil(serviceType)
	if !ok {
		panic(fmt.Sprintf("service %s not found", serviceType))
	}
	return v
}

// 获取对象，如果没有，返回nil
func (p *Provider) GetOrNil(serviceType reflect.Type) (any, bool) {
	descriptor := p.c.FirstOrDefault(serviceType)
	if descriptor == nil {
		return nil, false
	}
	return p.create(serviceType, descriptor), true
}

func (p *Provider) GetArray(baseType reflect.Type) []any {
	instances := make([]any, 0)
	for i := 0; i < len(p.c.Services); i++ {
		v := &p.c.Services[i]
		implType := reflect.PtrTo(v.ServiceType)
		if v.ServiceType == baseType || implType == baseType || implType.AssignableTo(baseType) {
			instances = append(instances, p.create(v.ServiceType, v))
		}
	}
	return instances
}

func (p *Provider) Close() {
	p.values.Range(func(key, value interface{}) bool {
		if closer, ok := value.(io.Closer); ok {
			err := closer.Close()
			if err != nil {
				log.Warn("close instance failed", err)
			}
		}
		return true
	})
}

func (p *Provider) create(serviceType reflect.Type, descriptor *ServiceDescriptor) any {
	if descriptor.Scope == Scoped {
		if v, ok := p.values.Load(serviceType); ok {
			return v
		}
		v, loaded := p.values.LoadOrStore(serviceType, p.createInstance(descriptor))
		if !loaded {
			p.instanceInit(descriptor.ServiceType, v, false)
		}
		return v
	}
	if descriptor.Scope == Singleton {
		if v, ok := p.c.singletons.Load(serviceType); ok {
			return v
		}
		v, loaded := p.c.singletons.LoadOrStore(serviceType, p.createInstance(descriptor))
		if !loaded {
			p.instanceInit(descriptor.ServiceType, v, true)
		}
		return v
	}
	v := p.createInstance(descriptor)
	p.instanceInit(descriptor.ServiceType, v, false)
	return v
}

func (p *Provider) createInstance(descriptor *ServiceDescriptor) any {
	if descriptor.Value != nil {
		return descriptor.Value
	}
	var instance any
	if descriptor.Creator != nil {
		instance = descriptor.Creator(p.c)
	} else {
		instance = reflect.New(descriptor.ServiceType).Interface()
	}
	return instance
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
