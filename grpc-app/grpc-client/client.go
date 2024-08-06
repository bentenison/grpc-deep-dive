package main

import (
	"context"
	"grpc-app/grpc-client/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// connect to the grpc server
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("unable to connect to grpc server service", err)
	}
	defer conn.Close()
	//register the user service client

	client := proto.NewUserServiceClient(conn)
	log.Println("connected to the client")
	//call the method on client
	resp, err := client.CreateUser(context.Background(), &proto.CreateUserRequest{User: &proto.User{
		Name:  "tanvir",
		Id:    "123",
		Email: "tanvirs.mkcl.org",
	}})
	if err != nil {
		log.Println("error recieving response from the grpc server", err)
		return
	}
	log.Println("respose recieved is", resp)
	select {}
}
