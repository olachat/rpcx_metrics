package prom

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
	consulapi "github.com/hashicorp/consul/api"
)

// ConsulDiscovery consul服务发现
type ConsulDiscovery struct {
	serviceName string   // 服务名称
	Ipv4        string   // 注册IP
	Port        int      // 注册端口
	consulAddr  []string // consul地址
	// consulPath  string   // consulPath
	healthCheckUseTcp bool // 健康检查是否使用tcp
	closed            bool
}

// NewConsulDiscovery 创建consul服务发现
func NewConsulDiscovery(serviceName string, port int, healthCheckUseTcp ...bool) *ConsulDiscovery {
	ipv4s, err := IP.LocalIPv4s()
	if err != nil {
		panic(err)
	}

	return &ConsulDiscovery{
		serviceName:       serviceName,
		Ipv4:              ipv4s[0],
		Port:              port,
		consulAddr:        g.Cfg().GetStrings("rpc.discover.Addr"),
		healthCheckUseTcp: len(healthCheckUseTcp) > 0 && healthCheckUseTcp[0],
		closed:            false,
	}
}

// getClient 获取consul客户端
func (c *ConsulDiscovery) getClient() (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddr[0]
	config.Scheme = "http"
	return consulapi.NewClient(config)
}

// getID 获取服务ID
func (c *ConsulDiscovery) getID() string {
	return fmt.Sprintf("%s:%d", c.Ipv4, c.Port)
}

// Register 执行注册逻辑
func (c *ConsulDiscovery) Register(tags []string, meta map[string]string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}
	if meta == nil {
		meta = map[string]string{}
	}
	meta["prefix"] = c.serviceName

	// 创建一个新服务。
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = c.getID()
	registration.Name = c.serviceName
	registration.Tags = tags
	registration.Address = c.Ipv4
	registration.Port = c.Port
	registration.Meta = meta

	// 增加check。
	if c.healthCheckUseTcp {
		registration.Check = &consulapi.AgentServiceCheck{
			TCP:                            registration.ID,
			Interval:                       "3s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "30s", // check失败后30秒删除本服务，注销时间，相当于过期时间
		}
	} else {
		registration.Check = &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/ping"),
			Interval:                       "3s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "30s", // check失败后30秒删除本服务，注销时间，相当于过期时间
		}
	}

	return client.Agent().ServiceRegister(registration)
}

// Deregister 解除注册
func (c *ConsulDiscovery) Deregister() error {
	if !c.closed {
		g.Log().Infof("consul unregister")
		client, err := c.getClient()
		if err != nil {
			g.Log().Println("consul getClient error", err)
			return err
		}
		if err = client.Agent().ServiceDeregister(c.getID()); err == nil {
			g.Log().Println("nginx close from", c.getID())
		}
		c.closed = true
	}
	return nil
}
