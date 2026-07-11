package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kidx45/Project-KK/api"
	db "github.com/kidx45/Project-KK/db/sqlc"
	"github.com/kidx45/Project-KK/grpcs"
	"github.com/kidx45/Project-KK/pb"
	"github.com/kidx45/Project-KK/utils"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	AppConfig, err := utils.LoadEnv(".env")
	if err != nil {
		log.Fatal("Can't load data because: ", err)
	}

	conn, err := sql.Open(AppConfig.DB_DRIVER_NAME, AppConfig.DB_URL)
	if err != nil {
		log.Fatal("Can't start server due to: ", err)
	}

	store := db.NewStore(conn)
	go startGatewayServer(AppConfig, store)
	startGRPCServer(AppConfig, store)
}

func startHTTPServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't create server due to: ", err)
	}
	address := fmt.Sprintf("0.0.0.0:%s", config.PORT)
	err = server.Start(address)
	if err != nil {
		log.Fatal("Can't start server due to: ", err)
	}
}
func startGRPCServer(config utils.Config, store db.Store) {
	gserver := grpc.NewServer()
	server, err := grpcs.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't create server due to: ", err)
	}
	pb.RegisterProjectKKServer(gserver, server)
	reflection.Register(gserver)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", config.GRPC_PORT))
	if err != nil {
		log.Fatal("Can't create listener due to: ", err)
	}
	log.Printf("Server started on: %s", listener.Addr().String())
	err = gserver.Serve(listener)
	if err != nil {
		log.Fatal("Can't start gRPC server due to: ", err)
	}
}

func startGatewayServer(config utils.Config, store db.Store) {
	server, err := grpcs.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't create server due to: ", err)
	}
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pb.RegisterProjectKKHandlerServer(ctx, mux, server)

	httpmux := http.NewServeMux()
	httpmux.Handle("/", mux)

	address := fmt.Sprintf("0.0.0.0:%s", config.PORT)
	log.Printf("Gateway server started on: %s", address)
	err = http.ListenAndServe(address, httpmux)
	if err != nil {
		log.Fatal("Can't start gateway server due to: ", err)
	}
}
