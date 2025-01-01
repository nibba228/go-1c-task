package main

import (
	"context"
	"fmt"

	pb "github.com/nibba228/go-1c-task/experiment"
	"google.golang.org/grpc"
)

const (
  port = flag.String("port", 50051, "port to connect to")
  name = flag.String("name", "noname", "username")
)

func main() {
  flag.Parse()
  opts := []grpc.DialOption{
    grpc.WithPerRPCCredentials(*name),
  }

  conn, _ := grpc.NewClient(fmt.Sprintf("localhost:%d", *port), opts...)
  defer conn.Close()

  client := pb.NewExperimentClient(conn)
  ctx := context.Background()

  response, err := client.Register(ctx)
  fmt.Printf(response)
}
