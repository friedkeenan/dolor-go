package protocol

import (
    "io"
    "bytes"
    "reflect"
    "encoding/json"
    "errors"
)

var (
    ErrorInvalidPacketField = errors.New("Invalid packet field")
    ErrorInvalidPacketStructTag = errors.New("Invalid packet struct tag")
)

const (
    StateHandshaking = iota
    StateStatus
    StateLogin
    StatePlay
)

type State VarInt

func (s *State) Read(r io.Reader) error {
    i := VarInt(*s)
    err := (&i).Read(r)
    if err != nil {
        return err
    }

    *s = State(i)
    return nil
}

func (s *State) Write(w io.Writer) error {
    i := new(VarInt)
    *i = VarInt(*s)
    return i.Write(w)
}

const (
    Serverbound = iota
    Clientbound
)

type Bound uint8

type Packet interface{}

func ReadPacket(c Connection, bound Bound) (Packet, error) {
    var length, id VarInt
    err := (&length).Read(c)
    if err != nil {
        return nil, err
    }

    tmp_buf := make([]byte, length)
    _, err = c.Read(tmp_buf)
    if err != nil {
        return nil, err
    }

    buf := bytes.NewBuffer(tmp_buf)

    err = (&id).Read(buf)
    if err != nil {
        return nil, err
    }

    packet_type, err := c.GetPacketType(bound, id)
    if err != nil {
        return nil, err
    }

    pack := reflect.New(packet_type).Elem()

    for i := 0; i < pack.NumField(); i++ {
        f := pack.Field(i)
        tag := packet_type.Field(i).Tag

        if tag != "" {
            switch tag.Get("as") {
                case "json":
                    s := new(String)
                    err = s.Read(buf)
                    if err != nil {
                        return nil, err
                    }

                    data := []byte(*s)
                    err = json.Unmarshal(data, f.Addr().Interface())
                    if err != nil {
                        return nil, err
                    }

                    continue

                default: return nil, ErrorInvalidPacketStructTag
            }
        } else {
            tmp, ok := f.Addr().Interface().(Type)
            if !ok {
                return nil, ErrorInvalidPacketField
            }

            err = tmp.Read(buf)
            if err != nil {
                return nil, err
            }
        }
    }

    return pack.Interface().(Packet), nil
}

func WritePacket(c Connection, pack Packet) error {
    packet_type := reflect.TypeOf(pack)
    id, err := c.GetPacketId(packet_type)
    if err != nil {
        return err
    }

    buf := bytes.NewBuffer([]byte {})

    err = (&id).Write(buf)
    if err != nil {
        return err
    }

    pack_v := reflect.ValueOf(pack)
    for i := 0; i < pack_v.NumField(); i++ {
        f := pack_v.Field(i)
        tag := packet_type.Field(i).Tag

        if tag != "" {
            switch tag.Get("as") {
                case "json":
                    tmp_v := reflect.New(f.Type())
                    tmp_v.Elem().Set(f)

                    data, err := json.Marshal(tmp_v.Interface())
                    if err != nil {
                        return err
                    }

                    s := String(data)
                    err = (&s).Write(buf)
                    if err != nil {
                        return err
                    }

                    continue
            }
        } else {
            tmp_v := reflect.New(f.Type())
            tmp_v.Elem().Set(f)

            tmp, ok := tmp_v.Interface().(Type)
            if !ok {
                return ErrorInvalidPacketField
            }

            err = tmp.Write(buf)
            if err != nil {
                return err
            }
        }
    }

    length := VarInt(buf.Len())
    err = (&length).Write(c)
    if err != nil {
        return err
    }

    _, err = buf.WriteTo(c)

    return err
}