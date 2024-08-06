package main

import (
	"context"
	"io"
	"log"

	pb "grpc-streaming/grpc-client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	//
	log.Println("starting the grpc client")
	cc, err := grpc.NewClient(":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error while stating grpc client\n", err)
	}
	client := pb.NewBookServiceClient(cc)
	res, err := client.CreateBook(context.Background(), &pb.CreateBookRequest{Book: &pb.Book{Name: "domain driven Design", Id: "13", Author: "tanvirs", CreatedAt: timestamppb.Now()}})
	if err != nil {
		log.Println("error occured while creating a book", err)
	}
	log.Println("got response from server", res.Success)
	resp, err := client.ListBooks(context.Background(), &pb.ListBookRequest{})
	if err != nil {
		log.Println("error occured while creating a book", err)
	}

	log.Println("got response from server", resp.Books)
	books := []*pb.Book{{Name: "Server Design", Id: "1", Author: "tanvirs", CreatedAt: timestamppb.Now()},
		{Name: "Microservices Design", Id: "2", Author: "tanvirs", CreatedAt: timestamppb.Now()},
		{Name: "System Design", Id: "3", Author: "tanvirs", CreatedAt: timestamppb.Now()}}
	ids := []*pb.GetBookRequest{{Id: "1"},
		{Id: "2"},
		{Id: "3"}}
	//client streaming
	ClientStreaming(books, client)
	//server streaming
	ServerStreaming(client)
	//bidirectional streaming
	BIDIStreaming(ids, client)
}
func ClientStreaming(books []*pb.Book, client pb.BookServiceClient) {
	stream, err := client.CreateSreamingBook(context.Background())
	if err != nil {
		log.Println("error occured while creating stream", err)
	}
	for _, v := range books {
		stream.Send(&pb.CreateBookRequest{Book: v})
	}

	data, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("error occured while recieving stream data", err)
		return
	}
	log.Printf("got response from server %s\n", data.Success)

}
func ServerStreaming(client pb.BookServiceClient) {
	stream, err := client.ListBooksStreaming(context.Background(), &pb.ListBookRequest{})
	if err != nil {
		log.Println("error occured while creating stream", err)
	}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return
		}
		log.Println("data recieved from server from server streaming", data)
		if err != nil {
			log.Println("error while recieving response", err)
			return
		}
	}
	// log.Printf("got response from server)
}
func BIDIStreaming(ids []*pb.GetBookRequest, client pb.BookServiceClient) {
	stream, err := client.GetBookBidirectionalStreaming(context.Background())
	if err != nil {
		log.Println("error occured while getting response from bidi streaming", err)
	}
	go func() {
		// stream
		for _, v := range ids {
			stream.Send(v)
		}
		stream.CloseSend()
	}()
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return
		}
		log.Println("data recieved from server from bidi streaming", data)
		if err != nil {
			status, ok := status.FromError(err)
			if ok {
				log.Println("more meaaningful error with", status.Code(), status.Details())
			}
			log.Println("error while recieving response", err)
			return
		}
	}
	// log.Printf("got response from server)
}
