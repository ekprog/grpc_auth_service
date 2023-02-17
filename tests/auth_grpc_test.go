package tests

import (
	"auth_service/domain"
	pb "auth_service/pkg/pb/api"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

func makeClient() *grpc.ClientConn {

	addr := "localhost:" + os.Getenv("GRPC_PORT")

	tslEnable := os.Getenv("TSL_ENABLE") == "true"
	if tslEnable {

		crt := "../cert/ca.cert"
		key := "../cert/ca.key"
		caN := "../cert/ca.cert"

		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			log.Fatalf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caN)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append ca certs")
		}

		creds := credentials.NewTLS(&tls.Config{
			ServerName:   addr, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		// Create a connection with the TLS credentials
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("could not dial %s: %s", addr, err)
		}
		return conn
	} else {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return conn
	}
}

func TestAuth_RegisterGRPC(t *testing.T) {

	conn := makeClient()
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Credentials
	testUsername, testPassword, err := generateCredentials()
	require.NoError(t, err, "should be success credentials generation process")

	res, err := c.Register(ctx, &pb.RegisterRequest{
		Username: testUsername,
		Password: testPassword,
	})
	log.Infof("%v", res)
	require.NoError(t, err, "should be success register ucase")
	require.Equal(t, res.Status.Code, domain.Success, "should be success status code")
}

func TestAuth_LoginGRPC(t *testing.T) {

	conn := makeClient()
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Credentials
	testUsername, testPassword, err := generateCredentials()
	require.NoError(t, err, "should be success credentials generation process")

	registerRes, err := c.Register(ctx, &pb.RegisterRequest{
		Username: testUsername,
		Password: testPassword,
	})
	require.NoError(t, err, "should be success register ucase")
	require.Equal(t, registerRes.Status.Code, domain.Success, "should be success status code with registration")

	r, err := c.Login(ctx, &pb.LoginRequest{
		Username: testUsername,
		Password: testPassword,
	})
	require.NoError(t, err, "should be success login ucase")
	require.Equal(t, r.Status.Code, domain.Success, "should be success status code with login")
	require.NotEmpty(t, r.JwtAccess.AccessToken, "access token should not be empty")
	require.NotEmpty(t, r.JwtAccess.RefreshExpiredAt, "access token should not be empty")

}

func TestAuth_ValidateGRPC(t *testing.T) {
	conn := makeClient()
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Credentials
	testUsername, testPassword, err := generateCredentials()
	require.NoError(t, err, "should be success credentials generation process")

	registerRes, err := c.Register(ctx, &pb.RegisterRequest{
		Username: testUsername,
		Password: testPassword,
	})
	require.NoError(t, err, "should be success register ucase")
	require.Equal(t, registerRes.Status.Code, domain.Success, "should be success status code with registration")

	r, err := c.Login(ctx, &pb.LoginRequest{
		Username: testUsername,
		Password: testPassword,
	})
	require.NoError(t, err, "should be success login ucase")
	require.Equal(t, r.Status.Code, domain.Success, "should be success status code with login")

	token := r.JwtAccess.AccessToken
	verify, err := c.Verify(ctx, &pb.VerifyRequest{AccessToken: token})
	require.NoError(t, err, "should be success verify ucase")
	require.Equal(t, verify.Status.Code, domain.Success, "status should be success")
	require.Equal(t, verify.User.Username, testUsername, "usernames after verification should be the same")
}

func TestAuth_RevokeGRPC(t *testing.T) {
	conn := makeClient()
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Credentials
	testUsername, testPassword, err := generateCredentials()
	require.NoError(t, err, "should be success credentials generation process")

	registerRes, err := c.Register(ctx, &pb.RegisterRequest{
		Username: testUsername,
		Password: testPassword,
	})
	require.NoError(t, err, "should be success register ucase")
	require.Equal(t, registerRes.Status.Code, domain.Success, "should be success status code with registration")

	r, err := c.Login(ctx, &pb.LoginRequest{
		Username: testUsername,
		Password: testPassword,
	})
	require.NoError(t, err, "should be success login ucase")
	require.Equal(t, r.Status.Code, domain.Success, "should be success status code with login")

	token := r.JwtAccess.AccessToken
	verify, err := c.Verify(ctx, &pb.VerifyRequest{AccessToken: token})
	require.NoError(t, err, "should be success verify ucase")
	require.Equal(t, verify.Status.Code, domain.Success, "status should be success with verifying")
	require.Equal(t, verify.User.Username, testUsername, "usernames after verification should be the same")

	// Revoking
	revokeRes, err := c.Revoke(ctx, &pb.RevokeRequest{AccessToken: token})
	require.NoError(t, err, "should be success revoke ucase")
	require.Equal(t, revokeRes.Status.Code, domain.Success, "status should be success with revoking")

	// Again verify
	verifyRes, err := c.Verify(ctx, &pb.VerifyRequest{AccessToken: token})
	require.NoError(t, err, "should be success verify ucase")
	require.Equal(t, verifyRes.Status.Code, domain.IncorrectToken, "token should be incorrect after revoking")
}
