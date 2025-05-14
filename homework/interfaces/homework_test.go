package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	services   map[string]func() interface{}
	singletons map[string]interface{}
}

func NewContainer() *Container {
	return &Container{services: make(map[string]func() interface{}), singletons: make(map[string]interface{})}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	serviceConstructor, ok := constructor.(func() interface{})
	if !ok {
		return
	}
	c.services[name] = serviceConstructor
}

func (c *Container) Resolve(name string) (interface{}, error) {
	service, ok := c.services[name]
	if !ok {
		return nil, errors.New("service not found")
	}
	return service(), nil
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) (interface{}, error) {
	serviceConstructor, ok := constructor.(func() interface{})
	if !ok {
		return nil, errors.New("service not found")
	}

	singleton, ok := c.singletons[name]
	if ok {
		return singleton, nil
	}

	c.singletons[name] = serviceConstructor()

	return c.singletons[name], nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	singleton1, err := container.RegisterSingletonType("SingletonUserService", func() interface{} {
		return &UserService{}
	})
	assert.NoError(t, err)

	singleton2, err := container.RegisterSingletonType("SingletonUserService", func() interface{} {
		return &UserService{}
	})
	assert.NoError(t, err)

	assert.True(t, singleton1 == singleton2)
}
