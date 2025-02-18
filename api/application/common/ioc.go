package common

import (
	"log"
	"reflect"

	"github.com/shelakel/go-ioc"
)

type IocScope interface {
	Scope() Ioc
}

type iocScopeImpl func() Ioc

func (ioc *iocScopeImpl) Scope() Ioc {
	return (*ioc)()
}

func NewIocScope(getter func() Ioc) IocScope {
	ioc := iocScopeImpl(getter)
	return &ioc
}

// this meant to decouple use cases from ioc container
//
// this doesn't decouple presentation layer which creates scope.
type Ioc interface {
	// note: service should be pointer
	InjectOptional(service any) error
	// this uses Inject under the hood but panics when no service was injected
	// note: service should be pointer
	Inject(service any)
}

type iocImpl struct {
	inject func(any) error
}

func (ioc *iocImpl) InjectOptional(service any) error {
	return ioc.inject(service)
}

func (ioc *iocImpl) Inject(service any) {
	if err := ioc.InjectOptional(service); err != nil {
		log.Panic(err.Error())
	}
}

func NewIoc(inject func(any) error) Ioc {
	return &iocImpl{
		inject: inject,
	}
}

func IocContainer(c *ioc.Container) Ioc {
	return NewIoc(func(a any) error { return c.Resolve(a) })
}

//

type ServiceGroup[T any] struct {
	services *[]T
}

func (services ServiceGroup[T]) Add(service T) {
	*services.services = append(*services.services, service)
}
func (services ServiceGroup[T]) GetAll() []T {
	if services.services == nil {
		return nil
	}
	return *services.services
}
func NewServiceGroup[T any]() ServiceGroup[T] {
	return ServiceGroup[T]{
		services: &[]T{},
	}
}

//

type ServiceStorage[T any] interface {
	Get() *T
	MustGet() T
	Set(T)
}

type serviceStorageImpl[T any] struct {
	service *T
}

func (storage *serviceStorageImpl[T]) Get() *T {
	return storage.service
}

func (storage *serviceStorageImpl[T]) MustGet() T {
	if storage.service == nil {
		log.Panicf("service storage of type \"%s\" is not set", reflect.TypeOf(storage.service).String())
	}
	return *storage.service
}

func (storage *serviceStorageImpl[T]) Set(service T) {
	storage.service = &service
}

func NewServiceStorage[T any]() ServiceStorage[T] {
	return &serviceStorageImpl[T]{}
}
