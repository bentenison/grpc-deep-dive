package main

import (
	"context"
	"fmt"
	pb "grpc-streaming/grpc-server/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// method promotuion for our grpc server
type server struct {
	pb.UnimplementedBookServiceServer
	books map[string]*pb.Book
}

func (s *server) ListBooks(ctx context.Context, req *pb.ListBookRequest) (*pb.ListBookResponse, error) {
	var bks []*pb.Book
	for _, v := range s.books {
		// log.Println("sending the book", *v)
		bks = append(bks, v)
	}
	return &pb.ListBookResponse{Books: bks}, nil
}
func (s *server) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	id := req.GetId()
	book, ok := s.books[id]
	if !ok {
		return &pb.GetBookResponse{Book: book}, status.Errorf(codes.NotFound, fmt.Sprintf("book not found for id %v", id))
	}
	return &pb.GetBookResponse{Book: book}, nil
}
func (s *server) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	book := req.GetBook()
	s.books[book.GetId()] = book
	return &pb.CreateBookResponse{Success: "book created successfully"}, nil
}
func (s *server) CreateSreamingBook(stream grpc.ClientStreamingServer[pb.CreateBookRequest, pb.CreateBookResponse]) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.CreateBookResponse{Success: "all books recieved!!!"})
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("error occured while recieving clioent stream %v", err))
		}
		s.books[data.Book.Id] = data.Book
		log.Println("recieving data", data.Book)
	}
}
func (s *server) ListBooksStreaming(req *pb.ListBookRequest, stream grpc.ServerStreamingServer[pb.Book]) error {
	for _, v := range s.books {
		if err := stream.Send(v); err != nil {
			log.Println("error", err)
			return err
		}
	}
	return nil
}
func (s *server) GetBookBidirectionalStreaming(stream grpc.BidiStreamingServer[pb.GetBookRequest, pb.GetBookResponse]) error {
	for {
		data, err := stream.Recv()
		if err == io.EOF && data == nil {
			return nil
		}
		if err != nil && err != io.EOF {
			log.Println("error occured while recieving stream from client", err)
			return err
		}
		log.Println("recieved data from client streaming", data)
		book, ok := s.books[data.GetId()]
		if !ok {
			return status.Errorf(codes.NotFound, "no data found for id")
		}
		if err := stream.Send(&pb.GetBookResponse{Book: book}); err != nil {
			log.Println(err)
			return err
		}
	}
}
func main() {
	log.Println("starting the server")
	//start the net listener
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("error while starting the listner svc", err)
	}
	defer lis.Close()
	//create a new grpc server
	s := grpc.NewServer()
	//register the service server with implemented methods
	pb.RegisterBookServiceServer(s, &server{books: make(map[string]*pb.Book)})
	log.Println("Started listening on 5000")
	//strat serving on grpc server
	if err := s.Serve(lis); err != nil {
		log.Fatal("error while listening for grpc transport")
	}
	select {}
}
