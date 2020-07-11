package protocol

import (
    "reflect"
    "errors"
)

var (
    ErrorInvalidPacketId = errors.New("Invalid packet id")
    ErrorInvalidPacketType = errors.New("Invalid packet type")
)

type Connection interface {
    /* Implements the io.ReadWriteCloser interface */
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error

    /* Makes the connection aware of the packet */
    RegisterPacket(state State, bound Bound, id VarInt, packet_type reflect.Type)

    SetState(state State)

    GetPacketType(bound Bound, id VarInt) (reflect.Type, error)
    GetPacketId(packet_type reflect.Type) (VarInt, error)
}