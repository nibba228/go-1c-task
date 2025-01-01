package main

import (
	"context"
  "flag"
  "fmt"
  "net"

	pb "github.com/nibba228/go-1c-task/experiment"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

const (
  guessInProgress = iota
  guessFinished
)

type Score struct {
  guessCount int
  guessStatus int
}

type Server struct {
  scoreTable map[string]Score
  mutex sync.Mutex
  start chan struct{}
}


func (s *Server) Register(ctx context.Context) (*pb.RegisterResponse, error) {
  _, ok <- s.start
  if !ok {
    return &pb.RegisterResponse{
      Status: "experiment has already started. you cannot join"
    }, error("experiment started")
  }

  data, ok := metadata.FromIncomingContext(ctx) 

  if !ok {
    return nil, error("could not extract context data")
  }

  mutex.Lock()
  username, ok = data["username"]
  mutex.Unlock()
  if !ok {
    return nil, error("name is not provided")
  }
  mutex.Lock()
  _, ok = s.scoreTable[username]
  if ok {
    mutex.Unlock()
    return &pb.RegisterResponse{
      Status: "a user with the same name has already registered"
    }, error("inappropriate username")
  }

  s.scoreTable[username] = Score{}
  mutex.Unlock()
  
  <-s.start
  return &RegisterResponse{
    Status: "starting!",
  }, nil
}

func (s *Server) Start(ctx context.Context) (*pb.StartResponse, error) {
  _, ok <- s.start
  if !ok {
    return &pb.StartResponse{
      Msg: "already started"
    }, nil
  }

  close(s.start)
  return &pb.StartResponse{
    Msg: "started"
  }, nil
}

func newServer() (*Server) {
  return &Server{
    scoreTable: make(map[string]Score),
    mutex: sync.Mutex{},
    start: make(chan struct{}),
  }
}

func main() {
  flag.Parse()
  listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))

  server := grpc.NewServer()
  pb.RegisterExperimentsServer(server, newServer())
  reflection.Register(server)
  server.Serve(listener)
}
