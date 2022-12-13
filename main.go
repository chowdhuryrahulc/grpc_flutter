package main

//todo PAGINATION IN DOCKER FLUTTER
//todo My app will have enums. Create another flutter app
//todo stats thread left

//* DIDNT UNDERSTAND:
// 1) grpc.WithStatsHandler(&ocgrpc.ClientHandler{}), : what is ocrpc? StatsHandler? And how to measure size of msg send?
// "google.golang.org/grpc/stats"?
// 2) stats package is left, to be done by employees perhaps. Can we clculate size with it?
// Calculate wirelength of the payload, set tags, etc. Employees do that
// 3) 

import (
	"context"
	"fmt"
	proto "lco/gen"
	"lco/statr"

	"log"
	"net"
	"os"
	"sync"
	prt "google.golang.org/protobuf/proto"		//! WOWW, this implements the same stuff as done by "github.com/golang/protobuf/proto"

	"google.golang.org/grpc"

	_ "google.golang.org/grpc/encoding/gzip"
	glog "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/stats"
)

var grpcLog glog.LoggerV2 //! NECESSARY FOR BERLINGER

func init() {
	// NewLoggerV2(info, warning, error)
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)	
}

// todo Why did we create this? Connection struct? What is its purpose?
type Connection struct {
	// proto.Broadcast_CreateStreamServer: helps stream msgs btw server and client
	stream proto.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error // channel error (go channels)(error type channel)
}

// we implement grpc on top of Server struct
type Server struct {
	Connection []*Connection
	proto.UnimplementedBroadcastServer
}

// Defining 2 protobuf methords (CreateStream and BroadcastMessage)
func (s *Server) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	// rpc CreateStream(Connect) returns (stream Message);
	conn := &Connection{ //todo Where are we using CreateStream and BroadcastMessage function?
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error), // making a new error channel
	}

	// adds this connection into Server.connection list
	s.Connection = append(s.Connection, conn)
	return <-conn.error // returns connection channel error
}

// ctx: grpc context
func (s *Server) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
	// rpc BroadcastMessage(Message) returns (Close);
	// grpcLog.Info("message", msg.Content) // only for testing

	wait := sync.WaitGroup{} // to implement go routines
	done := make(chan int)   // to know when all the go routines are finished
	
	size := prt.Size(msg) 	//! gets msg size
	grpcLog.Info("SIZEEEE BROD", size)	// returns 126, 184, 291 etc with compression and without compression
	// Size is measured in bytes
	//todo Why is the size same? With and without compression?

	for _, conn := range s.Connection {
		wait.Add(1)
		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done() // wait - 1

			if conn.active {
				err := conn.stream.Send(msg)                      // send messages back to client (client attached to conn)
				grpcLog.Info("Sending message to: ", conn.stream) // to show in CLI that msg has been send

				if err != nil {
					grpcLog.Errorf("Error with Stream: %s - Error: %v", &conn.stream, err)
					conn.active = false // if we fail to send the msg, connection is no longer active
					conn.error <- err
				}
			}
		}(msg, conn) //todo What does these arguments do?
	}

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done // this needs to return an item before the next line (return) is executed
	return &proto.Close{}, nil
}

func main() {
	var connections []*Connection

	server := &Server{connections, proto.UnimplementedBroadcastServer{}} //! Why did we write unimplement...?

	// Create a stats handler function.
	// statsHandler := func(statsk stats.RPCStats) {
	// 	switch stat := statsk.(type) {
	// 	case *stats.OutPayload:
	// 		// A message was sent.
	// 		fmt.Printf("Sent message of size %d bytes\n", stat.Length)
	// 	case *stats.InPayload:
	// 		// A message was received.
	// 		fmt.Printf("Received message of size %d bytes\n", stat.Length)
	// 	}
	// }

	// statsHandler := func(stats stats.Handler) stats.Handler {
	// 	outgoingTags := stats.OutgoingTags()
	// 	// Use the outgoing tags to record statistics about the RPC.
	// 	// ...
	//   }

	// func createStatsHandler() func(stats.StatsHandler) {
	// 	return func(stats stats.StatsHandler) {
	// 		switch stat := stats.(type) {
	// 		case *stats.OutPayload:
	// 			// A message was sent.
	// 			fmt.Printf("Sent message of size %d bytes\n", stat.Length)
	// 		case *stats.InPayload:
	// 			// A message was received.
	// 			fmt.Printf("Received message of size %d bytes\n", stat.Length)
	// 		}
	// 	}
	// }

	//// Find example of grpc.WithStatsHandler. And how to implement it
	//// ctx := context.Background()

	// handler, err := statshandler.NewServerHandler()
	// if err != nil {}

	// Create a gRPC server with the stats handler.
	grpcServer := grpc.NewServer()

	// grpc.Creds(creds),			//todo from GarageDoor or other

		// grpc.StatsHandler(&statr.StatrHandler{}),
		// grpcServer.RegisterService()
		// server.Use(statsHandler)
		// grpc.WithStatsHandler(),


		//? STATS NOTES (now not necessary probably. Check when you get time)
		// This worked bcox it has all the 4 functions. 
		// Even if 1 is missing, it wont work
		// So what .. can replace that. With official stats package
		// Also this is what &statr.StatrHandler{} and &ocgrpc.ClientHandler{} provided
		// grpc.StatsHandler(&myStatsHandler{}), // uncomment

		//// Deprecated grpc.WithCompressor(grpc.NewGZIPCompressor()), //? Also mosty done client side? Whyy?
		// grpc.UseCompressor(gzip.Name), //? GZIP is mostly used client side. Why?? Bcoz we also need it server side
		// grpc.UseCompressor(grpc.NewGZIPCompressor().Type()), // WTF is Type()?
		
		// grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5*time.Second}), //todo WTF is this?
		
		// grpc.WithStatsHandler(&ocgrpc.ClientHandler{}), // todo Make Improvements How the Fuck will this work? Bcoz everyone has used it
		// grpc.WithStatsHandler(stats.Handler.HandleRPC(ctx, stats.RPCStats)),
		// grpc.WithStatsHandler(stats.OutgoingTags(ctx context.Context)),
		// grpc.WithStatsHandler(createStatsHandler()),
	// )

	grpcLog.Info("Starting server at port :8080")
	// ERROR caused: in client, we use proto.BroadcastClient.
	// And in server we use proto.BroadcastServer

	// Register a gRPC service with the server.
	proto.RegisterBroadcastServer(grpcServer, server)
	//// func RegisterBroadcastServer(s grpc.ServiceRegistrar, srv BroadcastServer) {

	// Start the server.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}
	grpcServer.Serve(listener)
}

type myStatsHandler struct{}

// Implement the stats.Handler interface
func (h *myStatsHandler) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	fmt.Println(stat)
    // Handle RPC stats here
}

func (h *myStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
    // Tag the RPC here
    return ctx
}

func (g *myStatsHandler) TagConn(ctx context.Context, ct *stats.ConnTagInfo) context.Context {
	return ctx
}

func (g *myStatsHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
}


// Implement the other stats.Handler methods as needed
// This example creates a grpc.Server and adds a custom stats.Handler implementation called myStatsHandler as the stats handler. The myStatsHandler struct implements the stats.Handler interface, which requires you to implement the HandleRPC() and TagRPC() methods.
// You can then register your gRPC service with the server and start it by calling grpcServer.Serve(). When the server receives an RPC request, it will call the HandleRPC() and TagRPC() methods of the myStatsHandler to handle the stats for that request.
// This allows you to collect and process custom stats for your gRPC server.

/*
import (
    "context"
    "io"

    "google.golang.org/grpc/stats"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
)

startServer starts a gRPC server and enables stats reporting.
func startServer() {
    // Create a new gRPC server.
    server := grpc.NewServer()

    // Register your gRPC service with the server.
    // ...

    // Create a new stats handler.
    statsHandler := &stats.Handler{}

    // Enable stats reporting for the server.
    server.Use(statsHandler)

    // Start the server.
    // ...
}

// Your gRPC service implementation.
type MyService struct{}

func (s *MyService) MyMethod(ctx context.Context, req *MyRequest) (*MyResponse, error) {
    // Report stats for the RPC.
    statsHandler.HandleRPC(ctx, &stats.OutHeader{
        Client: true,
        Code:   codes.OK,
        Ctx:

*/