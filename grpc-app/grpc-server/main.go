package main

import (
	"context"
	"library/proto"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// register server type for our user service server
type server struct {
	// method promotion for user service server
	proto.UnimplementedUserServiceServer
	// simulate users in database for our app
	users map[string]*proto.User
}

func (s *server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user := req.GetUser()
	s.users[user.GetId()] = user
	return &proto.CreateUserResponse{User: user}, nil
}
func (s *server) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	id := req.GetId()
	user, ok := s.users[id]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "no user with the id found")
	}
	return &proto.GetUserResponse{User: user}, nil
}
func main() {
	log.Println("Starting the grpc server")
	// start a net listner
	list, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatal("error occured while starting listner")
	}
	// create a grpc server
	s := grpc.NewServer()
	log.Println("server is listning on the port 3000")
	// register the user service server with created server
	proto.RegisterUserServiceServer(s, &server{users: make(map[string]*proto.User)})
	// start listning for connctions on the grpc server
	if err := s.Serve(list); err != nil {
		log.Fatal("error starting the server")
	}
}
