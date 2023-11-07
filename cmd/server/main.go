package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	hh "github.com/MicahParks/httphandle"
	hhconst "github.com/MicahParks/httphandle/constant"
	"github.com/MicahParks/httphandle/middleware"

	aswu "github.com/MicahParks/aws-ses-web-ui"
	"github.com/MicahParks/aws-ses-web-ui/config"
	"github.com/MicahParks/aws-ses-web-ui/handle/api"
	"github.com/MicahParks/aws-ses-web-ui/handle/template"
	"github.com/MicahParks/aws-ses-web-ui/server"
)

func main() {
	setupArgs := hh.SetupArgs{
		Static:    aswu.Static,
		Templates: aswu.Templates,
	}
	setupRes, err := hh.Setup[config.Config](setupArgs)
	if err != nil {
		log.Fatalf(hhconst.LogFmt, "Failed to setup.", err)
	}
	conf := setupRes.Conf
	l := setupRes.Logger
	slog.SetDefault(l.With("global", true))
	l.Debug("Setup complete.")
	mux := http.NewServeMux()

	srv, err := server.NewServer(conf, l, setupRes.Templater)
	if err != nil {
		l.Error("Failed to create server.",
			hhconst.LogErr, err,
		)
		os.Exit(1)
	}

	apiHandlers := []hh.API[server.Server]{
		&api.Compose{},
	}
	generalHandlers := make([]hh.General[server.Server], 0)
	templateHandlers := []hh.Template[server.Server]{
		&template.Index{},
	}

	middlewareOpts := middleware.GlobalOptions{
		MaxReqSize: 10 * 1024 * 1024, // 10 MB
		ReqTimeout: 10 * time.Second,
	}
	attachArgs := hh.AttachArgs[server.Server]{
		API:            apiHandlers,
		Files:          setupRes.Files,
		General:        generalHandlers,
		MiddlewareOpts: middlewareOpts,
		Template:       templateHandlers,
		Templater:      setupRes.Templater,
	}
	err = hh.Attach(attachArgs, srv, mux)
	if err != nil {
		l.Error("Failed to attach handlers.",
			hhconst.LogErr, err,
		)
		os.Exit(1)
	}
	l.Debug("Handlers attached.")

	l.Info("Server listening on http://localhost:8080")
	serveArgs := hh.ServeArgs{
		Logger:          l.With("httphandle", true),
		Port:            8080,
		ShutdownFunc:    srv.Shutdown,
		ShutdownTimeout: 10 * time.Second,
	}
	hh.Serve(serveArgs, mux)
}
