package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"divoom-gallery-client/internal/config"
	"divoom-gallery-client/internal/device"
	"divoom-gallery-client/internal/gallery"
	"divoom-gallery-client/internal/handlers"
	"divoom-gallery-client/internal/httpx"
	"divoom-gallery-client/internal/media"
	"divoom-gallery-client/internal/metrics"
	"divoom-gallery-client/internal/playback"
	"divoom-gallery-client/internal/rotation"
	"divoom-gallery-client/internal/settings"
	"divoom-gallery-client/internal/source"
)

const (
	exitFailure            = 1
	defaultInterval        = 5 * time.Minute
	defaultDisplayDuration = 10 * time.Second
)

func main() {
	if err := run(); err != nil && !errors.Is(err, context.Canceled) {
		log.Print(err)
		os.Exit(exitFailure)
	}
}
func run() error {
	if len(os.Args) > 1 {
		return errors.New("command-line parameters are not supported; use environment variables and the web UI")
	}
	listen, err := config.LoadWebAddress()
	if err != nil {
		return err
	}
	dataDir := config.LoadDataDirectory()
	account, err := config.LoadGalleryAccount()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	password, err := config.LoadWebPassword()
	if err != nil {
		return err
	}
	d, err := config.LoadDeviceAPI()
	if err != nil {
		return err
	}
	metricsConfig, err := config.LoadMetrics()
	if err != nil {
		return err
	}
	defaults := settings.Defaults(0, defaultInterval, defaultDisplayDuration)
	prefs, err := settings.New(filepath.Join(dataDir, "settings.json"), defaults)
	if err != nil {
		return err
	}
	mediaStore, err := media.New(filepath.Join(dataDir, "media"))
	if err != nil {
		return err
	}
	var requestMetrics httpx.Metrics
	if metricsConfig.URL != "" {
		client := metrics.NewClient(metricsConfig)
		defer client.Close()
		requestMetrics = client
	}
	api := gallery.New(gallery.Options{Email: account.Email, Password: account.Password, Metrics: requestMetrics})
	animationSource := source.New(api, mediaStore, prefs)
	deviceDisplay := device.New(d, requestMetrics)
	if err := deviceDisplay.Health(ctx); err != nil {
		return fmt.Errorf("device API health: %w", err)
	}
	runner := rotation.New(prefs, animationSource, deviceDisplay)
	commands := playback.New(prefs, runner, animationSource, deviceDisplay)
	ui := handlers.New(animationSource, api, commands, prefs, runner, int(defaults.SourceID), password)
	return serveWeb(ctx, listen, ui, runner)
}
