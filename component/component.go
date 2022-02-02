package component

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/sustglobal/sust-go/event"
	"github.com/sustglobal/sust-go/httpapi"
)

func NewFromEnv() (*Component, error) {
	cfg := DefaultConfig()
	if err := LoadFromEnv(&cfg); err != nil {
		return nil, err
	}
	return New(cfg)
}

func New(cfg Config) (*Component, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	cmp := Component{
		Config: cfg,

		Logger:              logger,
		HTTPRouter:          mux.NewRouter(),
		InboundEventRouter:  event.NewEventRouter(logger),
		OutboundEventRouter: event.NewEventRouter(logger),
	}

	cmp.httpServer = &http.Server{
		Addr:    cfg.BindHTTPServer,
		Handler: cmp.HTTPRouter,
	}

	if cfg.ExposeMetrics {
		httpapi.NewMetricsHandler().Mount(cmp.HTTPRouter)
	}
	if cfg.ExposeHealth {
		httpapi.NewHealthHandler().Mount(cmp.HTTPRouter)
	}

	return &cmp, nil
}

type Component struct {
	Config Config

	HTTPRouter          *mux.Router
	InboundEventRouter  *event.EventRouter
	OutboundEventRouter *event.EventRouter
	Logger              *zap.Logger

	httpServer *http.Server
	asyncDone  chan struct{}
	asyncError error
}

func (c *Component) Start() error {
	var network, addr string
	if strings.HasPrefix(c.Config.BindHTTPServer, "unix://") {
		network = "unix"
		addr = c.Config.BindHTTPServer[7:]
	} else {
		network = "tcp"
		addr = c.Config.BindHTTPServer
	}

	l, err := net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("failed network bind: %v", err)
	}

	c.asyncDone = make(chan struct{})

	go func() {
		defer func() { close(c.asyncDone) }()

		err := c.httpServer.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			c.asyncError = err
		}
	}()

	return nil
}

func (c *Component) Run() error {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	return c.runUntil(sigc)
}

func (c *Component) runUntil(sigc <-chan os.Signal) error {
	c.Logger.Info("starting component", zap.String("config", fmt.Sprintf("%+v", c.Config)))

	if err := c.Start(); err != nil {
		c.Logger.Error("failed to start component", zap.Error(err))
		return err
	}

	c.Logger.Info("component up and running")

	sig := <-sigc

	c.Logger.Info("stopping component", zap.Reflect("signal", sig))

	if err := c.Stop(); err != nil {
		c.Logger.Error("failed to stop component", zap.Error(err))
		return err
	}

	c.Logger.Info("component stopped gracefully")
	return nil
}

func (c *Component) Stop() error {
	defer c.Logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), c.Config.GracefulShutdownTimeout)
	defer cancel()

	err := c.httpServer.Shutdown(ctx)

	select {
	case <-c.asyncDone:
	case <-ctx.Done():
	}

	return err
}

func (c *Component) Error() error { return c.asyncError }
