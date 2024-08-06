package main

import (
	"log"
	"net/rpc"
)

type User struct {
	Email string
	Name  string
	Id    string
}
type Args struct{}

// type RPCServer struct {
// 	users map[string]*User
// }

func main() {
	log.Println("starting the rpc client")
	conn, err := rpc.Dial("tcp", ":4000")
	if err != nil {
		log.Fatal("error while getting rpc connection", err)
	}
	defer conn.Close()
	log.Println("connnected to the server")
	user := &User{
		Name:  "tanvirs",
		Email: "tanvirs.example.com",
		Id:    "1",
	}
	var result *User

	err = conn.Call("RPCServer.CreateUser", user, &result)
	if err != nil {
		log.Println("error occured while calling the method", err)
	}
	log.Println("RPCServer.CreateUser returned response", *result)
	var users []*User
	// var data interface{}
	args := &Args{}
	err = conn.Call("RPCServer.GetAllUsers", args, &users)
	if err != nil {
		log.Println("error occured while calling the method", err)
	}
	for _, v := range users {
		log.Println("users", v)
	}
	// log.Printf("RPCServer.GetAllUsers returned response %v\n", users)
}
