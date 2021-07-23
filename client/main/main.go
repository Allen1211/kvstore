package main

import (
	"fmt"
	"net"
)

func main() {
	//argsWithoutProg := os.Args[1:]
	//fmt.Println(argsWithoutProg)
	//
	//readCommandArgs()
	//base.Logger.Println(Conf)

	conn, err := net.Dial("tcp", "10.132.36.8:9001")
	if err != nil {
		fmt.Printf("Client Thread %d: %v\n", 1, err)
		return
	}
	defer conn.Close()
	fmt.Printf("Client Thread %d Connected to Server\n", 1)

	if _, err := conn.Write([]byte("wtf")); err != nil {
		fmt.Printf("write hello err: %v\n", err)
		return
	}
	var buf = make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		fmt.Printf("read err: %v\n", err)
		return
	}
	fmt.Println(buf)
	//
	//var wg sync.WaitGroup
	//wg.Add(Conf.numThreads)
	//
	//for i := 0; i < Conf.numThreads; i++ {
	//	go func(threadId int) {
	//		defer func() {
	//			wg.Done()
	//			if r := recover(); r != nil {
	//				base.Logger.Printf("Recovered from goroutine %d, err: %v", threadId, r)
	//			}
	//		}()
	//		ctx := context.Background()
	//		ctx = context.WithValue(ctx, "threadId", threadId)
	//		clientWorker := NewClientWorker(ctx)
	//		clientWorker.KeepRequest()
	//	}(i)
	//}
	//wg.Wait()

	//base.Logger.Printf("Client Stopped.")

}


func readCommandArgs() {
	//flag.StringVar(&Conf.addr, "addr", "localhost:9000", "server addr (Default: localhost:9000)")
	//flag.UintVar(&Conf.i, "i", 0, "The share of Insert among all operations")
	//flag.UintVar(&Conf.u, "u", 0, "The share of Update among all operations")
	//flag.UintVar(&Conf.d, "d", 0, "The share of Delete among all operations")
	//flag.UintVar(&Conf.g, "g", 1, "The share of Get among all operations")
	//flag.UintVar(&Conf.r, "r", 0, "The share of Range among all operations")
	//flag.StringVar(&Conf.kr, "kr", "", "Key range")
	//flag.StringVar(&Conf.vr, "vr", "", "Val range")
	//flag.IntVar(&Conf.numThreads, "t", 1, "Num of threads (Default: 1)")
	//flag.IntVar(&Conf.numQueries, "n", -1, "Num of queries sent. input -1 to keep sending (Default: -1)")
	//flag.BoolVar(&Conf.enabledBatch, "b", false, "Enabled batch operations in one query")
	//flag.IntVar(&Conf.batchNum, "bn", 0, "Num of operations per batch query")
	//flag.BoolVar(&Conf.enabledLog, "log", false, "Enabled log result to file")
	//flag.StringVar(&Conf.logFile, "log-file", "", "result log file path")
	//flag.Parse()

	//if Conf.kr != "" {
	//	ss := strings.Split(Conf.kr, ",")
	//	if len(ss) != 2 {
	//		base.Logger.Fatalf("Wrong usage of arg 'kr'")
	//	}
	//	var err error
	//	if Conf.keyFrom, err = strconv.ParseInt(ss[0], 10, 64); err != nil {
	//		base.Logger.Fatalf("Wrong usage of arg 'kr'")
	//	}
	//	if Conf.keyTo, err = strconv.ParseInt(ss[1], 10, 64); err != nil {
	//		base.Logger.Fatalf("Wrong usage of arg 'kr'")
	//	}
	//} else {
	//
	//}
}
