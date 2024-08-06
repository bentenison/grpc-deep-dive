package main

import (
	"log"
	"net"
	"net/rpc"
)

type User struct {
	Email string
	Name  string
	Id    string
}
type Args struct{}
type RPCServer struct {
	users map[string]*User
}

func (s *RPCServer) GetAllUsers(args *Args, reply *[]*User) error {
	user := []*User{}
	for _, v := range s.users {
		user = append(user, v)
	}
	log.Println(user)
	*reply = user
	return nil
}

//	func (s *RPCServer) GetUserById(args string, reply *User) error {
//		user, ok := s.users[args]
//		if !ok {
//			log.Println("No user found for the IP")
//		}
//		*reply = *user
//		return nil
//	}
func (s *RPCServer) CreateUser(args *User, reply *User) error {
	log.Println("got hit from rpc client")
	s.users[args.Id] = args
	*reply = *args
	return nil
}

// func (s *RPCServer)
func main() {
	log.Println("Starting the RPC server")
	//register the rpc server service

	rpcServer := new(RPCServer)
	rpcServer.users = map[string]*User{}
	err := rpc.Register(rpcServer)
	if err != nil {
		log.Fatal("error while registering the service", err)
	}
	//start the listener
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("error occured while starting the listner\n")
	}
	log.Println("RPC server listening on port 4000")
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("error occured while accepting the connections")
			continue
		}
		go rpc.ServeConn(conn)
	}

}
