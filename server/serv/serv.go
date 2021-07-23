package serv

type KVStoreServer interface {
	Start(port int)
	Stop()
}
