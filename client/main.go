package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "github.com/nibba228/go-1c-task/experiment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
  port = flag.Int("port", 50051, "port to connect to")
  name = flag.String("name", "noname", "username")
)

func main() {
  flag.Parse()
  opts := []grpc.DialOption{
    grpc.WithTransportCredentials(insecure.NewCredentials()),
  }

  conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", *port), opts...)
  if err != nil {
    log.Fatalf("failed to connect to server. Error: %v", err)
  }
  defer conn.Close()
  client := pb.NewExperimentsClient(conn)
  ctx := context.Background()

  response, err := client.Register(ctx, &pb.RegisterRequest{Username: *name})
  if err != nil {
    log.Fatalf("failed to call Register: %v", err)
  }
  fmt.Println(response)

  var number uint64
  for {
    _, err := fmt.Scanf("%d", &number)
    if err != nil {
      log.Fatalf("could not scan the number: %v", err)
    }

    response, err := client.MakeGuess(ctx, &pb.GuessRequest{Guess: number, Username: *name})
    if err != nil {
      log.Fatalf("could not make a guess: %v", err)
    }

    fmt.Println(response.GetResult())
  }
}
