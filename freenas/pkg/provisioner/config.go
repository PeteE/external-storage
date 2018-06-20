package provisioner

import (
    "net/http"
    "crypto/tls"
    "github.com/Sirupsen/logrus"
    "github.com/spf13/viper"
)

type FreeNasConfig struct {
    Uri string
    Pool string
    Username string
    Password string
    Portal string
    IQN string
}

var log = logrus.New()

var httpClient = &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
}

func init() {
    // Disable TLS validation temporarily
    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func initLog() {
    var err error
    log.Level, err = logrus.ParseLevel(viper.GetString("log-level"))
    if err != nil {
        log.Fatalln(err)
    }
}
