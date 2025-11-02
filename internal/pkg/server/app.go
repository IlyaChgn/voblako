package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	authproto "github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/grpc/protobuf"
	authdel "github.com/IlyaChgn/voblako/internal/pkg/auth/delivery/rest"
	authuc "github.com/IlyaChgn/voblako/internal/pkg/auth/usecases"
	"github.com/IlyaChgn/voblako/internal/pkg/config"
	routers "github.com/IlyaChgn/voblako/internal/pkg/server/delivery"

	"github.com/gorilla/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	server *http.Server
}
type serverConfig struct {
	Address string
	Timeout time.Duration
	Handler http.Handler
}

func createServerConfig(addr string, timeout int, handler *http.Handler) serverConfig {
	return serverConfig{
		Address: addr,
		Timeout: time.Second * time.Duration(timeout),
		Handler: *handler,
	}
}

func createServer(config serverConfig) *http.Server {
	return &http.Server{
		Addr:         config.Address,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		Handler:      config.Handler,
	}
}

func (srv *Server) Run() error {
	cfgPath := os.Getenv("CONFIG_PATH")
	cfg := config.ReadConfig(cfgPath)
	if cfg == nil {
		log.Fatal("The config wasn`t opened")
	}

	credentials := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders(cfg.Server.Headers)
	originsOk := handlers.AllowedOrigins(cfg.Server.Origins)
	methodsOk := handlers.AllowedMethods(cfg.Server.Methods)

	authServiceURL := fmt.Sprintf("%s:%s", cfg.Auth.ExternalHost, cfg.Auth.Port)
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	authConn, err := grpc.NewClient(authServiceURL, opts)
	if err != nil {
		log.Fatal("Cannot create client for auth service", err)
	}
	defer authConn.Close()

	authClient := authproto.NewAuthClient(authConn)
	authUsecases := authuc.NewAuthUsecases(authClient)
	authHandler := authdel.NewAuthHandler(authUsecases)

	router := routers.NewRouter(authHandler)
	muxWithCORS := handlers.CORS(credentials, originsOk, headersOk, methodsOk)(router)

	serverURL := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	serverCfg := createServerConfig(serverURL, cfg.Server.Timeout, &muxWithCORS)
	srv.server = createServer(serverCfg)

	log.Printf("Server is listening on %s\n", serverURL)

	return srv.server.ListenAndServe()
}
