package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var ByteOrder binary.ByteOrder = binary.LittleEndian

/* request protocol */
type Request struct {
	Head   ReqHeader
	ReqOps []ReqOp
}

type ReqHeader struct { // 10 bytes
	Magic  uint16 // 0x1234
	Seq    uint32 // request sequence number
	Length uint32 // body length
}

type ReqOp struct { // 18 bytes
	Magic  uint8 // 0x59
	OpType uint8 // see base.op
	Var1   int64
	Var2   int64
}

func (req Request) String() string {
	return fmt.Sprintf("seq: %d; ops: %v", req.Head.Seq, req.ReqOps)
}

func (op ReqOp) String() string {
	return fmt.Sprintf("opType: %v; var1: %d, var2: %d", OpType(op.OpType), op.Var1, op.Var2)
}

var ReqHeadLen = uint32(binary.Size(ReqHeader{}))
var ReqOpLen = uint32(binary.Size(ReqOp{}))

func EncodeRequest(seq uint32, reqOps []ReqOp) ([]byte, error) {
	request := &Request{
		Head: ReqHeader{
			Magic:  0x1234,
			Seq:    seq,
			Length: ReqOpLen * uint32(len(reqOps)),
		},
		ReqOps: reqOps,
	}
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, ByteOrder, request.Head); err != nil {
		return nil, fmt.Errorf("failed to encode ReqHead: %v", err)
	}
	for _, reqOp := range reqOps {
		reqOp.Magic = 0x59
		if err := binary.Write(buf, ByteOrder, reqOp); err != nil {
			return nil, fmt.Errorf("failed to encode reqOp: %v", err)
		}
	}
	return buf.Bytes(), nil
}

func DecodeRequest(buf []byte) (*Request, error) {
	if uint32(len(buf)) < ReqHeadLen {
		return nil, fmt.Errorf("failed to decode request: package length(%d) less than expected", len(buf))
	}
	head := new(ReqHeader)
	if err := binary.Read(bytes.NewReader(buf[:ReqHeadLen]), ByteOrder, head); err != nil {
		return nil, fmt.Errorf("failed to decode ReqHeader: %v", err)
	}
	if head.Magic != 0x1234 {
		return nil, fmt.Errorf("failed to decode ReqHeader: magic(%d) wrong", head.Magic)
	}
	numReqOp := head.Length / ReqOpLen
	request := &Request{
		Head:   *head,
		ReqOps: make([]ReqOp, numReqOp),
	}

	ofs := ReqHeadLen
	for i := uint32(0); i < numReqOp; i++ {
		reqOp := new(ReqOp)
		if err := binary.Read(bytes.NewReader(buf[ofs:ofs+ReqOpLen]), ByteOrder, reqOp); err != nil {
			return nil, fmt.Errorf("failed to decode ReqOp: %v", err)
		}
		request.ReqOps[i] = *reqOp
		ofs += ReqOpLen
	}

	return request, nil
}

/* response protocol */
type Response struct {
	Head    RespHeader
	RespOps []RespOp
}

type RespHeader struct { // 12 bytes
	Magic   uint16 // 0x6789
	RetCode uint16 // return code
	Seq     uint32 // request sequence number
	Length  uint32 // body length
}

type RespOp struct { // 6 + PayLoadLen  bytes
	Head    RespOpHead
	PayLoad Payload
}

type RespOpHead struct {
	Magic      uint8 // 0xa5
	OpType     uint8 // see base.op
	PayloadLen uint32
}

type Payload struct { // length: for insert/update/delete:1bytes for get: 5bytes; for range: >= 8 bytes
	Ok   uint8
	Vals []int32
}

func (resp Response) String() string {
	return fmt.Sprintf("seq: %d; ops: %v", resp.Head.Seq, resp.RespOps)
}

func (op RespOp) String() string {
	return fmt.Sprintf("opType: %v; ok: %d; vals: %v", OpType(op.Head.OpType), op.PayLoad.Ok, op.PayLoad.Vals)
}



var RespHeadLen = uint32(binary.Size(RespHeader{}))
var RespOpHeadLen = uint32(binary.Size(RespOpHead{}))

func EncodeResponse(retCode uint16, seq uint32, respOps []RespOp) ([]byte, error) {
	var bodyLength uint32 = 0
	for _, respOp := range respOps {
		bodyLength += RespOpHeadLen + respOp.Head.PayloadLen
	}
	response := &Response{
		Head: RespHeader{
			Magic:  0x6789,
			Seq:    seq,
			RetCode: retCode,
			Length: bodyLength,
		},
		RespOps: respOps,
	}
	buf := &bytes.Buffer{}
	if err := binary.Write(buf, ByteOrder, response.Head); err != nil {
		return nil, fmt.Errorf("failed to encode RespHead: %v", err)
	}
	for _, respOp := range respOps {
		respOp.Head.Magic = 0xa5
		if err := binary.Write(buf, ByteOrder, respOp.Head); err != nil {
			return nil, fmt.Errorf("failed to encode respOp: %v", err)
		}
		if err := binary.Write(buf, ByteOrder, respOp.PayLoad.Ok); err != nil {
			return nil, fmt.Errorf("failed to encode respOp: %v", err)
		}
		for _, val := range respOp.PayLoad.Vals {
			if err := binary.Write(buf, ByteOrder, val); err != nil {
				return nil, fmt.Errorf("failed to encode respOp: %v", err)
			}
		}
	}
	return buf.Bytes(), nil
}

func DecodeResponse(buf []byte) (*Response, error) {
	if uint32(len(buf)) < RespHeadLen {
		return nil, fmt.Errorf("failed to decode response: package length(%d) less than expected", len(buf))
	}
	head := new(RespHeader)
	if err := binary.Read(bytes.NewReader(buf[:RespHeadLen]), ByteOrder, head); err != nil {
		return nil, fmt.Errorf("failed to decode ResqHeader: %v", err)
	}
	if head.Magic != 0x6789 {
		return nil, fmt.Errorf("failed to decode ResqHeader: magic(%d) wrong", head.Magic)
	}

	response := &Response{
		Head:    *head,
		RespOps: make([]RespOp, 0),
	}

	ofs := RespHeadLen
	for ofs < RespHeadLen+head.Length {
		opHead := new(RespOpHead)
		if err := binary.Read(bytes.NewReader(buf[ofs:ofs+RespOpHeadLen]), ByteOrder, opHead); err != nil {
			return nil, fmt.Errorf("failed to decode RespOpHead: %v", err)
		}
		if opHead.Magic != 0xa5 {
			return nil, fmt.Errorf("failed to decode opHead: magic(%d) wrong", head.Magic)
		}
		ofs += RespOpHeadLen

		payloadBytes := buf[ofs : ofs+opHead.PayloadLen]
		response.RespOps = append(response.RespOps, RespOp{
			Head:    *opHead,
			PayLoad: DecodePayload(payloadBytes, opHead.PayloadLen),
		})

		ofs += opHead.PayloadLen
	}

	return response, nil
}

func DecodePayload(payloadBytes []byte, payloadLen uint32) Payload {
	payload := Payload{}
	payload.Ok = payloadBytes[0]
	for j := uint32(1); j < payloadLen; j += 4 {
		val := int32(ByteOrder.Uint32(payloadBytes[j : j+4]))
		payload.Vals = append(payload.Vals, val)
	}
	return payload
}
