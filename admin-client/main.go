package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nibba228/go-1c-task/experiment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
  port = flag.Int("port", 50051, "port to connect to")
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

  var input string
  for {
    _, err := fmt.Scanf("%s", &input)
    if err != nil {
      log.Fatalf("Could not read from stdin")
    }

    if input == "start" {
      response, err := client.Start(ctx, &emptypb.Empty{})
      if err != nil {
        log.Fatalf("Could not start experiment: %v", err)
      }
      fmt.Println(response.GetMsg())
      break
    }
  }

  // Writer for printing GetScores response
  writer := csv.NewWriter(os.Stdout)
  writer.Comma = '\t'

  for {
    _, err := fmt.Scanf("%s", &input)
    if err != nil {
      log.Fatalf("Could not read from stdin")
    }

    switch input {
    case "users":
      fmt.Println("Getting the users, participating in the experiment:")
      stream, err := client.GetUsers(ctx, &emptypb.Empty{})
      if err != nil {
        log.Fatalf("could not get users: %v", err)
      }

      for {
        user, err := stream.Recv();
        if err == io.EOF {
          break
        }

        if err != nil {
          log.Fatalf("GetUsers error: %v", err)
        }

        fmt.Println(user.GetUsername())
      }

      fmt.Println("Got all the users\n")

    case "score":
      fmt.Println("Getting the score of each user:")
      err := writer.Write([]string{
        "User",
        "Attempts",
        "Status",
      })
      
      if err != nil {
        log.Fatalf("could print table: %v", err)
      }

      stream, err := client.GetScores(ctx, &emptypb.Empty{})
      if err != nil {
        log.Fatalf("could not get scores: %v", err)
      }

      for {
        score, err := stream.Recv()
        if err == io.EOF {
          break
        }

        if err != nil {
          log.Fatalf("GetScores error: %v", err)
        }
        attempts := strconv.FormatUint(score.GetAttemptCount(), 10)
        err = writer.Write([]string{
          score.GetUsername(),
          attempts,
          score.GetEnum().String(),
        })       
        writer.Flush()
      }

      
      fmt.Println("Got all the scores\n")

    default:
      fmt.Println("Incorrect request")
    }
  }
}
