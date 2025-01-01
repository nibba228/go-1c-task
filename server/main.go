package main

import (
	"context"
  "errors"
  "flag"
  "fmt"
  "net"
  "sync"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nibba228/go-1c-task/experiment"
	"google.golang.org/grpc"
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
  AttemptCount uint64
  GuessStatus uint64
}

type Server struct {
  pb.UnimplementedExperimentsServer
  scoreTable map[string]*Score
  mutex sync.Mutex
  start chan struct{}
  number uint64
}

func newServer(number uint64) (*Server) {
  return &Server{
    scoreTable: make(map[string]*Score),
    mutex: sync.Mutex{},
    start: make(chan struct{}),
    number: number,
  }
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
  ok := true
  select {
  case _, ok = <-s.start:
    ok = false
  default:
  }

  if !ok {
    return &pb.RegisterResponse{
      Status: "experiment has already started. you cannot join",
    }, errors.New("experiment started")
  }

  username := req.GetUsername()
  s.mutex.Lock()
  _, ok = s.scoreTable[username]
  if ok {
    s.mutex.Unlock()
    return &pb.RegisterResponse{
      Status: "a user with the same name has already registered",
    }, errors.New("inappropriate username")
  }

  s.scoreTable[username] = &Score{AttemptCount: 0, GuessStatus: guessInProgress}
  s.mutex.Unlock()
  
  <-s.start
  return &pb.RegisterResponse{
    Status: "starting!",
  }, nil
}

func (s *Server) Start(ctx context.Context, _ *emptypb.Empty) (*pb.StartResponse, error) {
  ok := true
  select {
  case _, ok = <-s.start:
    ok = false
  default:
  }

  if !ok {
    return &pb.StartResponse{
      Msg: "already started",
    }, nil
  }

  close(s.start)
  return &pb.StartResponse{
    Msg: "started",
  }, nil
}

func (s *Server) MakeGuess(ctx context.Context, request *pb.GuessRequest) (*pb.GuessResponse, error) {
  number := request.GetGuess()
  username := request.GetUsername()
  
  s.mutex.Lock()
  _, ok := s.scoreTable[username]
  s.mutex.Unlock()
  if !ok {
    return &pb.GuessResponse{Result: "error"}, errors.New("no such user")
  }

  go func(s *Server, number uint64, username string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    score := s.scoreTable[username]
    if score.GuessStatus == guessFinished {
      return
    }
    score.AttemptCount++
    if number == s.number {
      score.GuessStatus = guessFinished
    }
  }(s, number, username)

  if number == s.number {
    return &pb.GuessResponse{Result: "you guessed the number!"}, nil
  }

  if number < s.number {
    return &pb.GuessResponse{Result: "your number is less"}, nil
  }

  return &pb.GuessResponse{Result: "your number is greater"}, nil
}

func (s *Server) GetUsers(_ *emptypb.Empty, stream pb.Experiments_GetUsersServer) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()

  for username, _ := range s.scoreTable {
    err := stream.Send(&pb.UserResponse{Username: username})
    if err != nil {
      return err
    }
  }
  return nil
}

func (s *Server) GetScores(_ *emptypb.Empty, stream pb.Experiments_GetScoresServer) error {
  s.mutex.Lock()
  defer s.mutex.Unlock()

  for username, score := range s.scoreTable {
    err := stream.Send(&pb.ScoreResponse{
      Username: username,
      AttemptCount: score.AttemptCount,
      Enum: pb.GuessStatus(score.GuessStatus),
    })
    if err != nil {
      return err
    }
  }
  return nil
}

func main() {
  flag.Parse()
  listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

  server := grpc.NewServer()
  const number = 123
  pb.RegisterExperimentsServer(server, newServer(number))
  reflection.Register(server)
  server.Serve(listener)
}
