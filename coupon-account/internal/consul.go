package internal

import (
	"account/log"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

func (cr ConsulRegistry) Register(server *grpc.Server) error {
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
		grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	}
	fmt.Println("check is :", check)
	bytes, _ := json.Marshal(check)
	fmt.Println("check bytes is :", string(bytes))
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

func (cr ConsulRegistry) Deregister() error {
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", AppConf.ConsulConfig.Host, AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	err = client.Agent().ServiceDeregister(cr.Id)
	return err
}
