package common

import (
	"code.google.com/p/gcfg"
	"net"
	"flag"
	"fmt"
)

const DB_KEY_URI = "uri"

type Config struct{

	Auth struct {
		 User string
		 Pass string
	}

	Service struct {
		Redis    string
		Restcomm string
        Host string
	}

	Redis struct{
		Channel string
	}

	ServerPort struct{
		Main int
		Conference int
		Advertising int
		Portal int
	}

	Callback struct{
		Phone string
        Conference string
	}

	Messages struct{
		DialFrom string
		ConferenceWelcome string
		Question string
		Answer1 string
		ThanksForAttention string
		ThanksForAnswer string
	}
}

func NewConfig()(cfg Config){
	err := gcfg.ReadFileInto(&cfg, "demo.gcfg")
	if(err != nil){
		panic(err)
	}
    flag.StringVar(&cfg.Service.Host, "host", getLocalIp().String(), "host ip Address");
    flag.StringVar(&cfg.Service.Restcomm, "restcomm", cfg.Service.Restcomm, "Restcomm ip Address");
    flag.StringVar(&cfg.Service.Redis, "redis", cfg.Service.Redis, "Redis ip Address");
    flag.Parse()
	return cfg;
}

func (cfg Config)GetExternalAddress(port  int)(string){
	//ip := net.LookupIP(cfg.Service.Host)[0]
	ip := cfg.Service.Host
	return fmt.Sprintf("%s:%d", ip, port)
}

func getLocalIp() net.IP {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return nil
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP
            }
        }
    }
    return nil
}
