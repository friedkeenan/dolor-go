package protocol

import "reflect"

type PacketListener func (p Packet, s Server, c Connection) error

type Server interface {
    NewConnection() (Connection, error)
    HandleConnection(c Connection)

    RegisterPacketListener(packet_type reflect.Type, ln PacketListener)
}