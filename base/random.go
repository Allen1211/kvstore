package base

import (
	"math/rand"
)

// TODO rand对象的线程安全性

type RandGenerator struct {
	source  *rand.Rand
}

func NewRandGenerator(source *rand.Rand) *RandGenerator {
	return &RandGenerator{
		source: source,
	}
}

func (generator *RandGenerator) SetSeed(seed int64)  {
	generator.source.Seed(seed)
}

func (generator *RandGenerator) RandomKey() int64 {
	return RandInt64(generator.source)
}

func (generator *RandGenerator) RandomKeyByRange(from, to int64) int64 {
	return rand.Int63n(to-from) + from
}

func (generator *RandGenerator) RandomVal() int32 {
	return RandInt32(generator.source)
}

func (generator *RandGenerator) RandomValByRange(begin, end int32) int32 {
	return rand.Int31n(end) + begin
}

func (generator *RandGenerator) RandomOp() OpType {
	return OpType(rand.Intn(NumOfOp))
}

func (generator *RandGenerator) RandomBeginEnd() (int64, int64) {
	x := generator.RandomKey()
	y := generator.RandomKey()
	if x > y {
		return y, x
	} else {
		return x, y
	}
}

func (generator *RandGenerator) RandomBeginEndByRange(from, to int64) (int64, int64) {
	x := generator.RandomKeyByRange(from, to)
	y := generator.RandomKeyByRange(from, to)
	if x > y {
		return y, x
	} else {
		return x, y
	}
}

/* 随机创建请求包 */

func (generator *RandGenerator) CreateInsertOpReq() ReqOp {
	return ReqOp {
		OpType: OpInsert,
		Var1: generator.RandomKey(),
		Var2: int64(generator.RandomVal()),
	}
}

func (generator *RandGenerator) CreateUpdateOpReq() ReqOp {
	return ReqOp {
		OpType: OpUpdate,
		Var1: generator.RandomKey(),
		Var2: int64(generator.RandomVal()),
	}
}

func (generator *RandGenerator) CreateDeleteOpReq() ReqOp {
	return ReqOp {
		OpType: OpDelete,
		Var1: generator.RandomKey(),
	}
}

func (generator *RandGenerator) CreateGetOpReq() ReqOp {
	return ReqOp {
		OpType: OpGet,
		Var1: generator.RandomKey(),
	}
}

func (generator *RandGenerator) CreateRangeOpReq(from, to int64) ReqOp {
	begin, end := generator.RandomBeginEndByRange(from, to)
	return ReqOp {
		OpType: OpRange,
		Var1: begin,
		Var2: end,
	}
}



type RandOpProducer struct {
	opDistribution 	[]uint
	rightEnd		int
}

func NewRandomOpProducer(i, u, d, g, r uint) *RandOpProducer {
	opArr := []uint{i, u, d, g, r}
	allGcd := gcdN(opArr)
	for j := 0; j < len(opArr); j++ {
		opArr[j] /= allGcd
	}
	opDist := make([]uint, len(opArr))

	var right uint = 0
	for j := 0; j < len(opArr); j++ {
		if opArr[j] != 0 {
			right += opArr[j]
			opDist[j] = right
		}
	}

	return &RandOpProducer{
		opDistribution: opDist,
		rightEnd: int(right),
	}
}

func (p *RandOpProducer) Next() OpType {
	randNum := uint(rand.Intn(p.rightEnd))
	var from, to uint = 0, 0
	for j := 0; j < len(p.opDistribution); j++ {
		if p.opDistribution[j] == 0 {
			continue
		}
		to = from + p.opDistribution[j]
		if randNum >= from && randNum < to {
			return OpMap[j]
		}
	}
	return -1
}

func gcdN(nums []uint) uint {
	if len(nums) == 0 {
		return 0
	}
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		res = gcd(res, nums[i])
	}
	return res
}

func gcd(x, y uint) uint {
	if x < y {
		x, y = y, x
	}
	if y == 0 {
		return x
	} else {
		return gcd(y, x % y)
	}
}


func RandInt32(source *rand.Rand) int32 {
	return int32(source.Uint32())
}

func RandInt64(source *rand.Rand) int64 {
	//return int64(rand.Uint32()) | (int64(rand.Uint32()) << 32)
	return int64(source.Uint64())
}


/* 随机生成请求包 */
