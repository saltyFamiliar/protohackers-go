package means2anend

import (
	"encoding/binary"
	"net"
)

type RequestType byte

const (
	Insert RequestType = 'I'
	Query              = 'Q'
)

type RequestPacket [9]byte

func (p *RequestPacket) RequestType() RequestType {
	return RequestType(p[0])
}

func (p *RequestPacket) TimeStamp() int32 {
	return int32(binary.BigEndian.Uint32(p[1:5]))
}

func (p *RequestPacket) Price() int32 {
	return int32(binary.BigEndian.Uint32(p[5:9]))
}

func (p *RequestPacket) StartTime() int32 {
	return int32(binary.BigEndian.Uint32(p[1:5]))
}

func (p *RequestPacket) EndTime() int32 {
	return int32(binary.BigEndian.Uint32(p[5:9]))
}

func NewInsertionPacket() (p RequestPacket) {
	p[0] = byte(Insert)
	return p
}

func NewQueryPacket() (p RequestPacket) {
	p[0] = byte(Query)
	return p
}

func (p *RequestPacket) SetTimeStamp(val int32) {
	binary.BigEndian.PutUint32(p[1:5], uint32(val))
}

func (p *RequestPacket) SetPrice(val int32) {
	binary.BigEndian.PutUint32(p[5:9], uint32(val))
}

func (p *RequestPacket) SetStartTime(val int32) {
	binary.BigEndian.PutUint32(p[1:5], uint32(val))
}

func (p *RequestPacket) SetEndTime(val int32) {
	binary.BigEndian.PutUint32(p[5:9], uint32(val))
}

func (p *RequestPacket) SendOverConnection(conn net.Conn) {
	_, err := conn.Write(p[:])
	if err != nil {
		panic(err)
	}
}

type PriceHistory []RequestPacket

func InsertRequestPacket(ph PriceHistory, rp RequestPacket) PriceHistory {
	newPriceHistory := make(PriceHistory, len(ph)+1)

	if len(ph) == 0 {
		newPriceHistory[0] = rp
		return newPriceHistory
	}

	timeStamp := rp.TimeStamp()
	var insertionPoint int
	for i := 0; i < len(ph); i++ {
		if ph[i].TimeStamp() < timeStamp {
			newPriceHistory[i] = ph[i]
		}
		if ph[i].TimeStamp() >= timeStamp {
			insertionPoint = i
			newPriceHistory[i] = rp
		}
	}

	for i := insertionPoint + 1; i <= len(ph); i++ {
		newPriceHistory[i] = ph[i-1]
	}

	return newPriceHistory
}
