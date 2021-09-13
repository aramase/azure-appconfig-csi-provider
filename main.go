package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/aramase/azure-appconfig-csi-provider/pkg/server"
	"github.com/aramase/azure-appconfig-csi-provider/pkg/utils"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sigs.k8s.io/secrets-store-csi-driver/provider/v1alpha1"
)

var (
	endpoint = flag.String("endpoint", "unix:///tmp/azure-appconfig.sock", "CSI gRPC endpoint")
	check    = flag.Bool("check", false, "Check if the binary is working")

	log logr.Logger
)

func main() {
	flag.Parse()

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("unable to initialize logger: %v", err))
	}
	log = zapr.NewLogger(zapLog)
	log.WithName("azure-appconfig-csi-provider")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if *check {
		log.Info("Binary working!")
		return
	}

	// Initialize and run the gRPC server
	proto, addr, err := utils.ParseEndpoint(*endpoint)
	if err != nil {
		panic(err)
	}

	// setup provider gRPC server
	s := &server.Server{
		Log: log,
	}

	// remove the socket file if it already exists
	if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	// only unix domain socket is supported for now
	listener, err := net.Listen(proto, addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(utils.LogInterceptor(log)),
	}

	log.Info("starting azure appconfig provider server", "proto", proto, "addr", addr)

	g := grpc.NewServer(opts...)
	v1alpha1.RegisterCSIDriverProviderServer(g, s)
	go g.Serve(listener)

	<-ctx.Done()
	log.Info("Stopping server")
	g.GracefulStop()
}
