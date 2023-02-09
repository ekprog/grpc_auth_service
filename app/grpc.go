package app

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"net"
	"os"
)

var (
	grpcServer *grpc.Server
	grpcMux    *runtime.ServeMux
)

// GRPC
func InitGRPCServer() (*grpc.Server, *runtime.ServeMux, error) {

	tslEnable := os.Getenv("TSL_ENABLE") == "true"

	if tslEnable {
		crt := "./cert/service.pem"
		key := "./cert/service.key"
		caN := "./cert/ca.cert"

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, nil, errors.New("cannot initialize GRPC Server")
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caN)
		if err != nil {
			return nil, nil, errors.New("cannot initialize GRPC Server")
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, nil, errors.New("failed to append client certs")
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		grpcServer = grpc.NewServer(grpc.Creds(creds))

	} else {
		grpcServer = grpc.NewServer()
	}

	if grpcServer == nil {
		return nil, nil, errors.New("cannot initialize GRPC Server")
	}

	grpcMux = runtime.NewServeMux()

	return grpcServer, grpcMux, nil
}

func RunGRPCServer() {

	gRPCPort := os.Getenv("GRPC_PORT")

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", "localhost:"+gRPCPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("GRPC server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

// Register

type DeliveryService interface {
	Init() error
}

func InitDelivery(d ...DeliveryService) error {
	for _, service := range d {
		err := service.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func InitGRPCService[T any](s func(grpc.ServiceRegistrar, T), src T) {
	s(grpcServer, src)
}
