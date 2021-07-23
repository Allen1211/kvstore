package base

type OpType int

const (
	OpInsert = iota
	OpUpdate
	OpDelete
	OpGet
	OpRange
)

const NumOfOp = 5

func (opType OpType) String() string {
	switch opType {
	case OpInsert: return "insert"
	case OpUpdate: return "update"
	case OpDelete: return "delete"
	case OpGet:    return "get"
	case OpRange:  return "range"
	default: 	   return "invalid operation"
	}
}

var OpMap map[int]OpType

func init() {
	OpMap = make(map[int]OpType)
	OpMap[OpInsert] = OpInsert
	OpMap[OpUpdate] = OpUpdate
	OpMap[OpDelete] = OpDelete
	OpMap[OpGet] = OpGet
	OpMap[OpRange] = OpRange
}