package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/zeromicro/go-zero/core/netx"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	allEths  = "0.0.0.0"
	envPodIP = "POD_IP"

	defautlTTl    = 10
	defaultTicker = time.Second
)

type Conf struct {
	Host string
	Key  string
	Tags []string          `json:",optional"`
	Meta map[string]string `json:",optional"`
	TTL  int               `json:"ttl,optional"`
}

func Register(conf Conf, listenOn string) error {
	lo := figureOutListenOn(listenOn)
	host, pt, err := net.SplitHostPort(lo)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(pt)
	if err != nil {
		return err
	}
	client, err := api.NewClient(&api.Config{
		Scheme:  "http",
		Address: conf.Host,
	})
	if err != nil {
		return err
	}
	if conf.TTL < 0 {
		conf.TTL = defautlTTl
	}

	ttl := fmt.Sprintf("%ds", conf.TTL)
	expireTTL := fmt.Sprintf("%ds", conf.TTL*3)
	id := genID(conf.Key, host, port)
	&api.AgentServiceRegistration{
		ID:                id,
		Name:              conf.Key,
		Tags:              nil,
		Port:              0,
		Address:           "",
		SocketPath:        "",
		TaggedAddresses:   nil,
		EnableTagOverride: false,
		Meta:              nil,
		Weights:           nil,
		Check:             nil,
		Checks:            nil,
		Proxy:             nil,
		Connect:           nil,
		Namespace:         "",
		Partition:         "",
		Locality:          nil,
	}
	return nil
}

func genID(key, host string, port int) interface{} {
	return fmt.Sprintf("%s-%d-%d", key, host, port)
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields[0]) > 0 && fields[0] != allEths {
		return listenOn
	}
	ip := os.Getenv(envPodIP)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}
	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
