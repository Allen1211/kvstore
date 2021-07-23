package main

import (
	"context"
	"kvstore-server/base"
	rand "math/rand"
	"net"
	"time"
)

type ClientConf struct {
	addr	string		// server addr (Default: localhost:9000)

	i, u, d, g, r uint	// operation insert, update, delete, get, range

	kr	string		// key from,to example "-100,100"
	vr	string		// val from,to example "1,100000"
	rr	string		// begin:end

	keyFrom		int64
	keyTo		int64
	ValFrom		int32
	ValTo		int32

	numThreads int // num of threads

	numQueries   int  // num of query sent
	enabledBatch bool // enabled batch op
	batchNum     int  // num of operations per batch query

	enabledLog bool   // enabled result log
	logFile    string // result log file path
}

var Conf ClientConf

type ClientWorker struct {
	ctx				context.Context
	threadId		int
	randSource		*rand.Rand
	randGenerator 	*base.RandGenerator
	randOpProducer  *base.RandOpProducer
}

func NewClientWorker(ctx context.Context) *ClientWorker {
	threadId := ctx.Value("threadId").(int)
	randSource := rand.New(rand.NewSource(int64(threadId)))
	return &ClientWorker{
		ctx: ctx,
		threadId: threadId,
		randSource: randSource,
		randGenerator: base.NewRandGenerator(randSource),
		randOpProducer: base.NewRandomOpProducer(Conf.i, Conf.u, Conf.d, Conf.g, Conf.r),
	}
}


func (worker *ClientWorker) KeepRequest() {

	conn, err := net.Dial("tcp", Conf.addr)
	if err != nil {
		base.Logger.Printf("Client Thread %d: %v\n", worker.threadId, err)
		return
	}
	defer conn.Close()
	base.Logger.Printf("Client Thread %d Connected to Server\n", worker.threadId)

	if _, err := conn.Write([]byte("hello")); err != nil {
		base.Logger.Println("write hello err")
	}
	var buf = make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		base.Logger.Println("write hello err")
	}

	//encoderConfig := goframe.EncoderConfig{
	//	ByteOrder:                       base.ByteOrder,
	//	LengthFieldLength:               4,
	//	LengthAdjustment:                0,
	//	LengthIncludesLengthFieldLength: false,
	//}
	//decoderConfig := goframe.DecoderConfig{
	//	ByteOrder:           base.ByteOrder,
	//	LengthFieldOffset:   0,
	//	LengthFieldLength:   4,
	//	LengthAdjustment:    0,
	//	InitialBytesToStrip: 4,
	//}
	//fc := goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)

	tc := new(base.TimeCounter)

	for {
		select {
		case <-worker.ctx.Done():
			base.Logger.Printf("Client Thread %d Stopped\n", worker.threadId)
			break
		default:
		}

		tc.Reset()

		//reqOp := worker.createRandomOpReq()
		//reqData, err := base.EncodeRequest(1, []base.ReqOp{reqOp})
		//if err != nil {
		//	base.Logger.Printf("Client Thread %d: %v\n", worker.threadId, err)
		//	return
		//}

		//if err := fc.WriteFrame(reqData); err != nil {
		//	base.Logger.Printf("Client Thread %d: write err: %v\n",worker.threadId, err)
		//	return
		//}

		var buf []byte = make([]byte, 1024)
		if _, err = conn.Read(buf); err != nil {
			base.Logger.Printf("Client Thread %d: read err: %v\n",worker.threadId, err)
			return
		}
		//base.Logger.Println(buf)

		resp, err := base.DecodeResponse(buf)
		if err != nil {
			base.Logger.Println(err)
		}
		base.Logger.Printf("Client Thread %d: %v", worker.threadId, resp)

		base.Logger.Printf("Client Thread %d: Request time cost: %d\n",worker.threadId, tc.Count())

		time.Sleep(time.Second)
	}

}

func (worker *ClientWorker) createRandomOpReq() base.ReqOp {
	randGenerator := worker.randGenerator
	op := worker.randOpProducer.Next()
	switch op {
	case base.OpInsert:	return randGenerator.CreateInsertOpReq()
	case base.OpUpdate:	return randGenerator.CreateUpdateOpReq()
	case base.OpDelete: return randGenerator.CreateDeleteOpReq()
	case base.OpGet:	return randGenerator.CreateGetOpReq()
	case base.OpRange:	return randGenerator.CreateRangeOpReq(-1000000, 1000000)
	}
	panic("not impl")
}
