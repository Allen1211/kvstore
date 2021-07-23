package serv

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"kvstore-server/base"
	"kvstore-server/server/api"
)

type GnetKVStoreServer struct {
	*gnet.EventServer
	addr       string
	multicore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *goroutine.Pool
}

func (kvServer *GnetKVStoreServer) Start(port int) {
	goroutinePool := goroutine.Default()
	defer goroutinePool.Release()

	server := &GnetKVStoreServer{
		addr:       fmt.Sprintf("tcp://:%d", port),
		multicore:  true,
		async:      true,
		codec:      initCodec(),
		workerPool: goroutinePool,
	}
	err := gnet.Serve(server, server.addr,
		gnet.WithMulticore(server.multicore),
		//gnet.WithTCPKeepAlive(time.Minute*5),
		gnet.WithCodec(server.codec),
	)
	if err != nil {
		base.Logger.Fatalf("failed to start server! err: %v", err)
	}
}

func (kvServer *GnetKVStoreServer) Stop() {
	base.Logger.Printf("server stop")
}

func (kvServer *GnetKVStoreServer) OnInitComplete(server gnet.Server) (action gnet.Action) {
	base.Logger.Printf("Server Init\n")
	return
}

func (kvServer *GnetKVStoreServer) OnShutdown(server gnet.Server) {
	base.Logger.Println("Server shutdown")
}

func (kvServer *GnetKVStoreServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	base.Logger.Printf("connection %s opened\n", c.RemoteAddr())
	return
}

func (kvServer *GnetKVStoreServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	base.Logger.Printf("connection %s closed err: %v", c.RemoteAddr(), err)
	return
}

func (kvServer *GnetKVStoreServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, frame...)
	base.Logger.Printf("receive frame: %v\n", frame)
	// Use ants pool to unblock the event-loop.
	//_ = kvServer.workerPool.Submit(func() {
	//	tc := new(base.TimeCounter)
	//	tc.Reset()
	//	processRequest(data, c)
	//	base.Logger.Printf("Request Process time: %d\n", tc.Count())
	//})

	return data, action
}

func initCodec() gnet.ICodec {
	encoderConfig := gnet.EncoderConfig{
		ByteOrder:                       base.ByteOrder,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}
	decoderConfig := gnet.DecoderConfig{
		ByteOrder:           base.ByteOrder,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}
	return gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
}


func processRequest(data []byte, c gnet.Conn) {
	if len(data) == 0 {
		base.Logger.Println("receive empty package")
		return
	}

	request, err := base.DecodeRequest(data)
	if err != nil {
		base.Logger.Println(err.Error())
		return
	}
	base.Logger.Printf("receive request: %v\n", request)

	respOps := make([]base.RespOp, len(request.ReqOps))
	for i, op := range request.ReqOps {
		var respOp = base.RespOp {
			Head: base.RespOpHead {
				OpType: op.OpType,
			},
		}

		switch int(op.OpType) {
		case base.OpInsert:
			{
				key := op.Var1
				val := int32(op.Var2)
				ok := api.Insert(key, val)
				base.Logger.Printf("Execute operation <insert> k:%d, v:%d, res:%v\n", key, val, ok)

				respOp.PayLoad.Ok = bool2Uint8(ok)
				respOp.Head.PayloadLen = 1
			}
		case base.OpUpdate:
			{
				key := op.Var1
				val := int32(op.Var2)
				ok := api.Update(key, val)
				base.Logger.Printf("Execute operation <update> k:%d, v:%d, res:%v\n", key, val, ok)

				respOp.PayLoad.Ok = bool2Uint8(ok)
				respOp.Head.PayloadLen = 1
			}
		case base.OpDelete:
			{
				key := op.Var1
				ok := api.Delete(key)
				base.Logger.Printf("Execute operation <delete> k:%d res:%v\n", key, ok)

				respOp.PayLoad.Ok = bool2Uint8(ok)
				respOp.Head.PayloadLen = 1
			}
		case base.OpGet:
			{
				key := op.Var1
				val := api.Get(key)
				base.Logger.Printf("Execute operation <get> k:%d res:%v\n", key, val)

				respOp.PayLoad.Ok = 1
				respOp.PayLoad.Vals = append(respOp.PayLoad.Vals, val)
				respOp.Head.PayloadLen = 5
			}
		case base.OpRange:
			begin := op.Var1
			end := op.Var2
			vals := api.Range(begin, end)
			base.Logger.Printf("Execute operation <range> begin:%d end:%d res:%v\n", begin, end, vals)

			respOp.PayLoad.Ok = 1
			respOp.PayLoad.Vals = vals
			respOp.Head.PayloadLen = 1 + 4*uint32(len(vals))
		}

		respOps[i] = respOp
	}

	respData, err := base.EncodeResponse(0, request.Head.Seq, respOps)
	if err != nil {
		base.Logger.Println(err)
		return
	}

	if err := c.AsyncWrite(respData); err != nil {
		base.Logger.Fatalf(err.Error())
	}
	base.Logger.Println("response write success")
}

func bool2Uint8(b bool) uint8 {
	if b {
		return 1
	} else {
		return 0
	}
}