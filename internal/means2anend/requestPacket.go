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

func (p *RequestPacket) SetTimeStamp(val int32) {
	binary.BigEndian.PutUint32(p[1:5], uint32(val))
}

func (p *RequestPacket) SetPrice(val int32) {
	binary.BigEndian.PutUint32(p[5:9], uint32(val))
}

func NewInsertionPacket(time, price int32) (p RequestPacket) {
	p[0] = byte(Insert)
	p.SetTimeStamp(time)
	p.SetPrice(price)
	return p
}

func (p *RequestPacket) StartTime() int32 {
	return int32(binary.BigEndian.Uint32(p[1:5]))
}

func (p *RequestPacket) EndTime() int32 {
	return int32(binary.BigEndian.Uint32(p[5:9]))
}

func (p *RequestPacket) SetStartTime(val int32) {
	binary.BigEndian.PutUint32(p[1:5], uint32(val))
}

func (p *RequestPacket) SetEndTime(val int32) {
	binary.BigEndian.PutUint32(p[5:9], uint32(val))
}

func NewQueryPacket(startTime, endTime int32) (p RequestPacket) {
	p[0] = byte(Query)
	p.SetStartTime(startTime)
	p.SetEndTime(endTime)
	return p
}

func (p *RequestPacket) SendOverConnection(conn net.Conn) error {
	_, err := conn.Write(p[:])
	return err
}
