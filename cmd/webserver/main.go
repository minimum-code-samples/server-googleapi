package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"server-googleapi/lg"
	"server-googleapi/server"
	"server-googleapi/tpl"

	"gopkg.in/yaml.v2"
)

const (
	// ConfigFilePath specifies the name and location of the configuration file.
	ConfigFilePath = "./config/web.yaml"
	// LogFilename specifies the name of the log file. The location is determined by the configuration file.
	LogFilename = "web.log"
)

func main() {
	conf, gac, wd := parseFlags()

	tpl.Load(wd)
	config := prepConfig(conf)
	prepLog(config)

	if gac != "" {
		// If file is specified via a flag, override the configuration setting.
		config.GoogleApplicationCredentials = gac
	}

	s := server.NewServer(config)
	// TODO Change the configuration property to another one.
	s.Ready = isGoogleCredentialsAvail(config.GoogleApplicationCredentials)
	s.MakeRouter(false)
	runServer(s)
}

func isGoogleCredentialsAvail(gacPath string) bool {
	info, err := os.Stat(gacPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func parseFlags() (conf, gac, wd string) {
	flag.StringVar(&conf, "conf", ConfigFilePath, "Reads the configuration file for server. Default is 'config/web.yaml'")
	flag.StringVar(&gac, "gac", "", "Reads the Google application credentials file for accessing protected resources.")
	flag.StringVar(&wd, "wd", ".", "Sets the working directory for unit tests.")
	flag.Parse()
	return
}

func prepConfig(configPath string) server.Config {
	var cfg server.Config

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf(lg.FatalConfigParse, err)
	}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf(lg.FatalConfigParse, err)
	}
	return cfg
}

func prepLog(cfg server.Config) {
	filename := cfg.LogFolder + "/" + LogFilename
	if strings.HasSuffix(cfg.LogFolder, "/") {
		filename = cfg.LogFolder + LogFilename
	}
	lg.Init(filename, cfg.LogLevel)
}

// redirectHTTP returns a handler that performs a HTTP redirection to the
// secure port.
func redirectHTTP(port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := strings.Split(r.Host, ":")
		http.Redirect(w, r, fmt.Sprintf("https://%v:%v%v", h[0], port, r.RequestURI), http.StatusPermanentRedirect)
	}
}

func runServer(s *server.Server) error {
	errs := make(chan error)

	if s.PortTLS != "" {
		lg.Info(lg.ServerInitSecure, s.Interface, s.Port, s.PortTLS)
		// Run server with TLS.
	} else {
		lg.Info(lg.ServerInit, s.Interface, s.Port)

		go func() {
			srv := &http.Server{
				Handler: s.Router,
				Addr:    fmt.Sprintf("%s:%s", s.Interface, s.Port),
				// Do not set timeout on the read/write yet.
				// ReadTimeout: s.Config.ReadTimeout,
				// WriteTimeout: s.Config.WriteTimeout,
			}
			if e := srv.ListenAndServe(); e != nil {
				errs <- e
			}
		}()
		lg.Info(lg.ServerStarted)
	}

	return <-errs
}
