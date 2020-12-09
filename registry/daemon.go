package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"ykjam/doc-registry-go/api"
	"ykjam/doc-registry-go/config"
	"ykjam/doc-registry-go/datastore"
	"ykjam/doc-registry-go/web"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	quitChan := make(chan interface{})

	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "01-02 15:04:05.000"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.DebugLevel)

	err := config.ReadConfig("config.json")
	if err != nil {
		log.WithError(err).Panic("error reading config file")
	}

	setupServer(quitChan, signalChan, config.Conf)
}

func setupServer(quit chan interface{}, signalChan chan os.Signal, conf *config.Config) {
	var err error
	var access datastore.Access
	access, err = datastore.NewPgAccess(conf)
	if err != nil {
		log.WithError(err).Panic("Could not initialize datastore.Access")
		return
	}

	apiController := api.NewAPIController(access)
	if apiController == nil {
		log.Panic("API Controller is nil")
		return
	}

	s := web.NewServer(apiController)
	r := mux.NewRouter()

	r.HandleFunc("/api/organization", s.HandleOrganizationList)

	srv := &http.Server{
		Addr:         conf.ListenAddress,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 45 * time.Second,
	}

	var listener net.Listener
	listener, err = net.Listen("tcp", conf.ListenAddress)
	if err != nil {
		log.WithError(err).Panic("Error in setting up listener")
		return
	}

	log.WithField("listen", conf.ListenAddress).Info("Starting HTTP API Server")
	go startServer(srv, listener)

	for {
		select {
		case <-quit:
			log.Warn("quit channel closed, closing listener")
			err = srv.Close()
			if err != nil {
				log.WithError(err).Error("error during HTTP Server close")
			}
			err = listener.Close()
			if err != nil {
				log.WithError(err).Error("error during TCP Listener close")
			}
			return
		case sig := <-signalChan:
			switch sig {
			case os.Interrupt, os.Kill, syscall.SIGTERM:
				log.Info("interrupt signal received, sending Quit signal")
				close(quit)
			default:
				log.WithField("signal", sig).Info("signal received")
			}
		}
	}
}

func startServer(srv *http.Server, listener net.Listener) {
	err := srv.Serve(listener)
	if err != nil {
		log.WithError(err).Error("HTTP server Error")
	}
}
