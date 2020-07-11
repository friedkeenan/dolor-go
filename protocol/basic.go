package protocol

import (
    "net"
    "reflect"

    "fmt"
)

var (
    PacketIds = map[State]map[Bound]map[VarInt]reflect.Type {
        StateHandshaking: map[Bound]map[VarInt]reflect.Type {
            Serverbound: map[VarInt]reflect.Type {
                0x00: reflect.TypeOf(HandshakeServerbound{}),
            },
        },

        StateStatus: map[Bound]map[VarInt]reflect.Type {
            Serverbound: map[VarInt]reflect.Type {
                0x00: reflect.TypeOf(RequestServerbound{}),
                0x01: reflect.TypeOf(PingServerbound{}),
            },

            Clientbound: map[VarInt]reflect.Type {
                0x00: reflect.TypeOf(ResponseClientbound{}),
                0x01: reflect.TypeOf(PongClientbound{}),
            },
        },
    }
)

type BasicServer struct {
    Ln net.Listener

    PacketListeners map[reflect.Type][]PacketListener
}

type BasicConnection struct {
    Conn net.Conn
    Server *BasicServer

    CurrentState State
}

func NewBasicServer(port string) (*BasicServer, error) {
    ln, err := net.Listen("tcp", port)
    if err != nil {
        return nil, err
    }

    s := &BasicServer{
        Ln: ln,
        PacketListeners: map[reflect.Type][]PacketListener {},
    }

    s.RegisterPacketListener(reflect.TypeOf(HandshakeServerbound{}), HandleHandshake)

    s.RegisterPacketListener(reflect.TypeOf(RequestServerbound{}), HandleRequest)
    s.RegisterPacketListener(reflect.TypeOf(PingServerbound{}), HandlePing)

    return s, nil
}

func (s *BasicServer) RegisterPacket(state State, bound Bound, id VarInt, packet_type reflect.Type) {
    PacketIds[state][bound][id] = packet_type
}

func (s *BasicServer) GetPacketType(state State, bound Bound, id VarInt) (reflect.Type, error) {
    packet_type, ok := PacketIds[state][bound][id]

    if !ok {
        return nil, ErrorInvalidPacketId
    }

    return packet_type, nil
}

func (s *BasicServer) GetPacketId(state State, packet_type reflect.Type) (VarInt, error) {
    for id, typ := range PacketIds[state][Clientbound] {
        if typ == packet_type {
            return id, nil
        }
    }

    for id, typ := range PacketIds[state][Serverbound] {
        if typ == packet_type {
            return id, nil
        }
    }

    return -1, ErrorInvalidPacketType
}

func (s *BasicServer) NewConnection() (Connection, error) {
    conn, err := s.Ln.Accept()
    if err != nil {
        return nil, err
    }

    return &BasicConnection{
        Conn: conn,
        Server: s,
        CurrentState: StateHandshaking,
    }, nil
}

func (s *BasicServer) HandleConnection(c Connection) {
    defer c.Close()

    var err error = nil
    for {
        var p Packet
        p, err = ReadPacket(c, Serverbound)
        if err != nil {
            break
        }

        lns, ok := s.PacketListeners[reflect.TypeOf(p)]
        if ok {
            for _, ln := range lns {
                go ln(p, s, c)
            }
        }
    }

    if err != nil {
        fmt.Println(err)
    }
}

func (s *BasicServer) RegisterPacketListener(packet_type reflect.Type, ln PacketListener) {
    s.PacketListeners[packet_type] = append(s.PacketListeners[packet_type], ln)
}

func (c *BasicConnection) Read(buf []byte) (int, error) {
    return c.Conn.Read(buf)
}

func (c *BasicConnection) Write(buf []byte) (int, error) {
    return c.Conn.Write(buf)
}

func (c *BasicConnection) Close() error {
    return c.Conn.Close()
}

func (c *BasicConnection) RegisterPacket(state State, bound Bound, id VarInt, packet_type reflect.Type) {
    c.Server.RegisterPacket(state, bound, id, packet_type)
}

func (c *BasicConnection) SetState(state State) {
    c.CurrentState = state
}

func (c *BasicConnection) GetPacketType(bound Bound, id VarInt) (reflect.Type, error) {
    return c.Server.GetPacketType(c.CurrentState, bound, id)
}

func (c *BasicConnection) GetPacketId(packet_type reflect.Type) (VarInt, error) {
    return c.Server.GetPacketId(c.CurrentState, packet_type)
}

/* Packet listeners */

func HandleHandshake(p Packet, s Server, c Connection) error {
    fmt.Println("Handshake:", p)

    c.SetState(p.(HandshakeServerbound).NextState)

    return nil
}

func HandleRequest(p Packet, s Server, c Connection) error {
    fmt.Println("Request:", p)

    resp := ResponseClientbound{}

    resp.Response.Version.Name = "1.15.2"
    resp.Response.Version.Protocol = 578

    resp.Response.Players.Max = 20
    resp.Response.Players.Online = 0
    resp.Response.Players.Sample = []NameIdPair {}

    resp.Response.Description = map[string]interface{} {
        "text": "Test",
        "bold": true,
    }

    return WritePacket(c, resp)
}

func HandlePing(p Packet, s Server, c Connection) error {
    fmt.Println("Ping:", p)

    pong := PongClientbound{p.(PingServerbound).Payload}
    return WritePacket(c, pong)
}
