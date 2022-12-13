// stats for grpc:
// service/method
//   request counter #
//   in-wire length counter # gauge
//   out-wire length counter # gauge
//   request latency # gauge
//
// remote
//   addr1: counter, latency
//   addr2: counter, latency
//
// name -> span(all, current_5min)
//

package statr

import (
	"context"
	// "fmt"

	"google.golang.org/grpc/stats"
)

// grpc stats handler
type StatrHandler struct{}

func (g *StatrHandler) TagRPC(ctx context.Context, tag *stats.RPCTagInfo) context.Context {
	return ctx
}
func (g *StatrHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	// fmt.Printf("handle %#v\n", rs)
	// handle &stats.Begin{Client:false, BeginTime:time.Date(2022, time.November, 16, 18, 38, 13, 148314000, time.Local), FailFast:false, IsClientStream:true, IsServerStream:true, IsTransparentRetryAttempt:false}
	// handle &stats.InPayload{Client:false, Payload:(*gen.Servo)(0xc0000c02d0), Data:[]uint8{0xa, 0x6, 0x71, 0x75, 0x65, 0x75, 0x65, 0x64, 0x12, 0x5, 0x3a, 0x35, 0x36, 0x30, 0x31}, Length:15, WireLength:20, RecvTime:time.Date(2022, time.November, 16, 18, 38, 13, 148384000, time.Local)}
	// handle &stats.OutHeader{Client:false, Compression:"", Header:metadata.MD{}, FullMethod:"", RemoteAddr:net.Addr(nil), LocalAddr:net.Addr(nil)}
	// handle &stats.OutPayload{Client:false, Payload:(*gen.Fin)(0xc0000aa6c0), Data:[]uint8{}, Length:0, WireLength:5, SentTime:time.Date(2022, time.November, 16, 18, 38, 13, 148445000, time.Local)}

	// switch v := rs.(type) {
	// case *stats.InHeader:
	// 	// DefaultEngine.Incr("in",
	// 	// 	T("method", v.FullMethod),
	// 	// 	T("client", v.Client),
	// 	// )
	// 	// DefaultEngine.AddSample("in.header-wire", v.WireLength,
	// 	// 	T("method", v.FullMethod),
	// 	// 	T("client", v.Client),
	// 	// )
	// case *stats.InPayload:
	// 	// DefaultEngine.AddSample("in.payload-wire", v.WireLength,
	// 	// 	T("client", v.Client),
	// 	// )
	// case *stats.OutHeader:
	// 	// DefaultEngine.Incr("out.head",
	// 	// 	T("method", v.FullMethod),
	// 	// 	// T("local", v.LocalAddr.String()),
	// 	// 	// T("remote", v.RemoteAddr.String()),
	// 	// )
	// case *stats.OutPayload:
	// 	// DefaultEngine.AddSample("out.payload-wire", v.WireLength,
	// 	// 	T("client", v.Client),
	// 	// )

	// case *stats.Begin:
	// 	if v.IsClientStream || v.IsServerStream {
	// 		// DefaultEngine.Incr("stream",
	// 		// 	T("client", v.IsClientStream),
	// 		// 	T("server", v.IsServerStream),
	// 		// )
	// 	}
	// case *stats.End:
	// 	// DefaultEngine.AddSample("latency", v.EndTime.Sub(v.BeginTime).Microseconds(),
	// 	// 	T("client", v.Client),
	// 	// 	T("error", fmt.Sprintf("%v", v.Error)),
	// 	// )
	// }
}
func (g *StatrHandler) TagConn(ctx context.Context, ct *stats.ConnTagInfo) context.Context {
	return ctx
}

// type conn struct {
// 	counter   int  `metric:"grpc.conn" type:"counter"`
// 	is_client bool `metric:"grpc.conn.client" type:"gauge"`
// }

func (g *StatrHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	// switch cs.(type) {
	// case *stats.ConnBegin:
	// 	// DefaultEngine.Incr("active")
	// case *stats.ConnEnd:
	// 	// DefaultEngine.Decr("active")
	// 	// DefaultEngine.Incr("conn")
	// }
}
