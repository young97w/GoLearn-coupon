package internal

import (
	"account/log"
	"fmt"
	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Host string `mapstruscture:"host"`
	Port int    `mapstructure:"port"`
}

type ConsulRegistry struct {
	Host   string
	Port   int
	Method string
	Name   string
	Id     string
	Tags   []string
}

func NewConsulRegistry(host, method, name, id string, tags []string, port int) ConsulRegistry {
	return ConsulRegistry{
		Host:   host,
		Port:   port,
		Method: method,
		Name:   name,
		Id:     id,
		Tags:   tags,
	}
}

func Register(cr ConsulRegistry) error {
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", AppConf.ConsulConfig.Host, AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	serverAddr := fmt.Sprintf("%s:%d", cr.Host, cr.Port)
	check := api.AgentServiceCheck{
		Interval:                       "1s",
		Timeout:                        "3s",
		DeregisterCriticalServiceAfter: "6s",
	}
	if cr.Method == "http" {
		check.HTTP = serverAddr
	} else {
		check.GRPC = serverAddr
	}
	agentServiceRegistration := api.AgentServiceRegistration{}

	agentServiceRegistration.Address = cr.Host
	agentServiceRegistration.Port = cr.Port
	agentServiceRegistration.ID = cr.Id
	agentServiceRegistration.Name = cr.Name
	agentServiceRegistration.Tags = cr.Tags
	agentServiceRegistration.Check = &check

	err = client.Agent().ServiceRegister(&agentServiceRegistration)
	return err
}

func Deregister(cr ConsulRegistry) error {
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", AppConf.ConsulConfig.Host, AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	err = client.Agent().ServiceDeregister(cr.Id)
	return err
}
