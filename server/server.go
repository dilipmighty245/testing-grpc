package main

import (
	"context"
	pb "github.com/dilipmighty/testing-grpc/proto/greeter"
	"golang.org/x/sync/errgroup"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	ctx := context.Background()
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting server at:%q", addr)
	err = run(ctx, lis)

	switch err {
	case context.Canceled:
		log.Printf("graceful shutdown")
	default:
		log.Fatalf("failed to serve: %v", err)
	}
}

// run will start the server
func run(ctx context.Context, lis net.Listener) error {

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &Server{})
	errGrp, egCtx := errgroup.WithContext(ctx)

	errGrp.Go(func() error {
		return s.Serve(lis)
	})

	errGrp.Go(func() error {
		<-egCtx.Done()
		s.GracefulStop()
		return nil
	})

	return errGrp.Wait()
}
