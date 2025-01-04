package server

import (
	"log"
	"log/slog"
	"mesh-mqtt/clients/noaa"
	"mesh-mqtt/clients/radio"
	"mesh-mqtt/server/hooks"
	"os"

	"github.com/cockroachdb/pebble"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/hooks/debug"
	ph "github.com/mochi-mqtt/server/v2/hooks/storage/pebble"
	"github.com/mochi-mqtt/server/v2/listeners"
)

type Server struct {
	server *mqtt.Server
	radio  *radio.RadioClient
	noaa   *noaa.NOAAClient
	db     *pebble.DB
}

func (s *Server) Close() {
	s.radio.Close()
	s.server.Close()
}

func (s *Server) Serve() error {
	return s.server.Serve()
}

func (s *Server) Broadcast(msg string, opt radio.BroadcastOpt) {
	if s.radio != nil {
		s.radio.Send(msg, opt.To, opt.Channel)
	}
}

func NewServer(opts ...func(*Server)) *Server {
	s := &Server{
		server: mqtt.New(&mqtt.Options{
			InlineClient: true,
		}),
	}
	s.server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: nil, //nil == info
	}))

	for _, o := range opts {
		o(s)
	}
	return s
}

func WithPebbleDB(db *pebble.DB) func(*Server) {
	return func(s *Server) {
		s.db = db
	}
}

func WithRadioClient(r *radio.RadioClient) func(*Server) {
	return func(s *Server) {
		s.radio = r
	}
}

func WithNOAAClient(n *noaa.NOAAClient) func(*Server) {
	return func(s *Server) {
		s.noaa = n
	}
}

func WithPebbleHook(path string) func(*Server) {
	return func(s *Server) {
		err := s.server.AddHook(new(ph.Hook), &ph.Options{
			Path: path,
			Mode: ph.NoSync,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WithOpenAuthHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(auth.AllowHook), nil); err != nil {
			log.Fatal(err)
		}
	}
}

func WithExampleHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(hooks.ExampleHook), &hooks.ExampleHookOptions{
			Server: s.server,
			Radio:  s.radio,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithUserHook() func(*Server) {
	return func(s *Server) {
		if err := s.server.AddHook(new(hooks.UserHook), &hooks.UserHookOptions{
			Server: s.server,
			Radio:  s.radio,
			DB:     s.db,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithDebugHook() func(*Server) {
	return func(s *Server) {
		level := new(slog.LevelVar)
		level.Set(slog.LevelDebug)
		s.server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))

		err := s.server.AddHook(new(debug.Hook), &debug.Options{
			Enable:         true,
			ShowPacketData: true,
			ShowPasswords:  true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WithWeatherHook() func(*Server) {
	return func(s *Server) {
		level := new(slog.LevelVar)
		level.Set(slog.LevelDebug)
		s.server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
		if err := s.server.AddHook(new(hooks.WeatherHook), &hooks.WeatherHookOptions{
			Server: s.server,
			Radio:  s.radio,
			NOAA:   s.noaa,
			DB:     s.db,
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func WithDefaultListener() func(*Server) {
	return func(s *Server) {
		tcp := listeners.NewTCP(listeners.Config{
			ID:      "t1",
			Address: ":1883",
		})
		err := s.server.AddListener(tcp)
		if err != nil {
			log.Fatal(err)
		}
	}
}
