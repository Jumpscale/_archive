package main

import (
	"net/url"
	"os"

	"github.com/Jumpscale/aydostorex/config"
	"github.com/Jumpscale/aydostorex/fs"
	"github.com/Jumpscale/aydostorex/rest"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
)

var (
	log *logging.Logger
)

type Options struct {
	LogLevel   int
	DebugMode  bool
	ConfigPath string
}

func configureLogging(options Options) {
	logging.SetLevel(logging.Level(options.LogLevel), "")
	formatter := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{module} %{level:.1s} > %{message} %{color:reset}")
	logging.SetFormatter(formatter)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(backend)
	log = logging.MustGetLogger("aydostorex")
}

func main() {
	opts := Options{}

	app := cli.NewApp()
	app.Name = "AydoStoreX"
	app.Usage = ""
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       "config.toml",
			Usage:       "Path to the configuration file",
			Destination: &opts.ConfigPath,
		},
		cli.IntFlag{
			Name:        "log, l",
			Value:       4,
			Usage:       "Level of logging (0 less verbose, to 5 most verbose) default to 4",
			Destination: &opts.LogLevel,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Launch web server in debug mode",
			Destination: &opts.DebugMode,
		},
	}
	app.Action = func(c *cli.Context) {
		cfg, err := config.LoadConfig(opts.ConfigPath)
		if err != nil {
			log.Fatalf("Error while loading config file: %v", err)
		}

		if !opts.DebugMode {
			gin.SetMode(gin.ReleaseMode)
		}

		backendStores := []*url.URL{}
		for _, source := range cfg.Sources {
			u, err := url.Parse(source)
			if err != nil {
				log.Fatalf("URL %s not valid : %v", source, err)
			}
			backendStores = append(backendStores, u)
		}
		fsStore := fs.NewStore(cfg.StoreRoot)
		restService := rest.NewService(fsStore, nil, backendStores)

		r := rest.Router(cfg.Authentification, restService)

		if cfg.Tls.Certificate != "" && cfg.Tls.Key != "" {
			log.Info("Enable TLS")
			log.Errorf("%s", r.RunTLS(cfg.ListenAddr, cfg.Tls.Certificate, cfg.Tls.Key).Error())
		} else {
			log.Errorf("%s", r.Run(cfg.ListenAddr).Error())
		}
	}

	configureLogging(opts)
	log.Info("Start server")
	app.Run(os.Args)
}
