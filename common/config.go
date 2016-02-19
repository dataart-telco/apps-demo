package common

import (
	"github.com/scalingdata/gcfg"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

const DB_KEY_URI = "uri"
const CHANNEL_CALL_STATUS = "call_status"
const CHANNEL_CONF_DROPPED = "conf_dropped"

type Config struct {
	Auth struct {
		User string
		Pass string
	}

	Service struct {
		Redis    string
		Restcomm string
		Host     string
		Recorder string
		Opencell string
	}

	Redis struct {
		MainChannel   string
		ConfChannel   string
		InProgressKey string
	}

	ServerPort struct {
		Main        int
		Conference  int
		Advertising int
		Portal      int
		Sms         int
	}

	Callback struct {
		Phone      string
		Conference string
		Sms        string
	}

	Messages struct {
		DialFrom           string
		ConferenceWelcome  string
		Question           string
		Answer1            string
		ThanksForAttention string
		ThanksForAnswer    string
		SmsMessage         string
		SmsMessageSimple   string
	}

	Sip struct {
		DidProvider string
	}
}

func NewConfig() (cfg Config) {
	err := gcfg.ReadFileInto(&cfg, "demo.gcfg")
	if err != nil {
		panic(err)
	}
	flag.StringVar(&cfg.Service.Host, "host", getLocalIp().String(), "host ip Address")
	flag.StringVar(&cfg.Service.Restcomm, "restcomm", cfg.Service.Restcomm, "Restcomm ip Address")
	flag.StringVar(&cfg.Service.Redis, "redis", cfg.Service.Redis, "Redis ip Address")
	flag.StringVar(&cfg.Service.Recorder, "rec", cfg.Service.Recorder, "Recorder host")

	flag.StringVar(&cfg.Auth.User, "r-user", cfg.Auth.User, "Restcomm user")
	flag.StringVar(&cfg.Auth.Pass, "r-pass", cfg.Auth.Pass, "Restcomm password")

	flag.StringVar(&cfg.Callback.Phone, "r-phone-incom", cfg.Callback.Phone, "Incoming phone number")
	flag.StringVar(&cfg.Callback.Conference, "r-phone-conf", cfg.Callback.Conference, "Conference phone number")
	flag.StringVar(&cfg.Callback.Sms, "r-phone-sms", cfg.Callback.Sms, "Sms phone number")
	flag.StringVar(&cfg.Sip.DidProvider, "dp", cfg.Sip.DidProvider, "Did provider sip domain")

	l := flag.String("l", "INFO", "Log level: TRACE, INFO")

	flag.Parse()

	var traceHandle io.Writer
	if *l == "TRACE" {
		traceHandle = os.Stdout
	} else {
		traceHandle = ioutil.Discard
	}

	InitLog(traceHandle, os.Stdout, os.Stdout, os.Stdout)
	return cfg
}

func (cfg Config) GetExternalAddress(port int) string {
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
