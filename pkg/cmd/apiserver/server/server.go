package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nsyszr/ariane/pkg/api/handler"
	"github.com/nsyszr/ariane/pkg/cmd/apiserver/config"
	"github.com/spf13/cobra"
)

type server struct {
	Config *config.Config
	nc     *nats.Conn
	errCh  chan error
	wg     sync.WaitGroup
	h      *handler.Handler
}

func newServer(c *config.Config) (*server, error) {
	s := &server{
		Config: c,
		errCh:  make(chan error, 1),
		wg:     sync.WaitGroup{},
	}

	nc, err := nats.Connect(nats.DefaultURL,
		nats.DrainTimeout(10*time.Second),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			fmt.Printf("\n\nerror handler: %s\n\n", err)
			s.errCh <- err
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			fmt.Printf("\n\nclosed handler\n\n")
			s.wg.Done()
		}),
		nats.DisconnectHandler(func(_ *nats.Conn) {
			// TODO(DGL) this method is called twice when NATS server is going
			// offline. 1st when server gone and 2nd when the shutdown/drain is
			// initiated.
			fmt.Printf("\n\ndisconnect handler\n\n")
			// s.wg.Done()
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			//s.quitCh <- os.Interrupt
		}))
	if err != nil {
		return nil, err
	}

	s.nc = nc
	s.h = handler.NewHandler(nc)

	return s, nil
}

func (s *server) Serve() error {
	log.Print("Starting server")

	s.wg.Add(1)

	// Subscribe
	/*if _, err := s.nc.Subscribe("api.core.v1.namespace", func(msg *nats.Msg) {
		data, err := s.h.Handle(msg.Subject, msg.Data)
		if err != nil {
			log.Print("handle request error: ", err)
		}
		s.nc.Publish(msg.Reply, data)
	}); err != nil {
		log.Fatal(err)
	}*/
	if err := s.h.SubscribeAPIGroups(); err != nil {
		// TODO(DGL) we should return an error and not kill the proccess
		log.Fatal(err)
	}

	log.Print("Server started successfully")

	s.wg.Wait()

	// Check if there was an error
	select {
	case err := <-s.errCh:
		log.Print("Received an error: ", err)
		return err
	default:
		return nil
	}
}

func (s *server) Shutdown() {
	log.Print("Shutting down server")
	if s.nc != nil {
		s.nc.Drain()
	}
}

func (s *server) Close() {
	if s.nc != nil {
		s.nc.Close()
	}
	log.Print("Server shutdown successfully")
}

func RunServe(c *config.Config) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Create a new server
		s, err := newServer(c)
		if err != nil {
			log.Fatal(err)
			// os.Exit(1)
		}
		defer s.Close()

		// Run main loop
		go func() {
			if err := s.Serve(); err != nil {
				log.Fatal(err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server
		quitCh := make(chan os.Signal)
		signal.Notify(quitCh, os.Interrupt)
		<-quitCh

		// Shutdown the server
		s.Shutdown()
	}
}
