package protocol

import (
    "io"
    "encoding/binary"
    "errors"
)

var (
    ErrorVarIntTooBig  = errors.New("VarInt is too big")
    ErrorVarLongTooBig = errors.New("VarLong is too big")
)

type Type interface {
    Read(r io.Reader) error
    Write(w io.Writer) error
}

type Boolean       bool
type Byte          int8
type UnsignedByte  uint8
type Short         int16
type UnsignedShort uint16
type Int           int32
type Long          int64
type Float         float32
type Double        float64
type VarInt        int32
type VarLong       int64
type String        string
type Position      uint64
type Angle         int8

func (b *Boolean) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, b)
}

func (b *Boolean) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, b)
}

func (b *Byte) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, b)
}

func (b *Byte) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, b)
}

func (b *UnsignedByte) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, b)
}

func (b *UnsignedByte) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, b)
}

func (s *Short) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, s)
}

func (s *Short) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, s)
}

func (s *UnsignedShort) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, s)
}

func (s *UnsignedShort) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, s)
}

func (i *Int) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, i)
}

func (i *Int) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, i)
}

func (l *Long) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, l)
}

func (l *Long) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, l)
}

func (f *Float) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, f)
}

func (f *Float) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, f)
}

func (d *Double) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, d)
}

func (d *Double) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, d)
}

func (i *VarInt) Read(r io.Reader) error {
    var tmp VarInt = 0
    var num_read uint8 = 0

    for {
        var b uint8
        err := binary.Read(r, binary.BigEndian, &b)
        if err != nil {
            return err
        }

        tmp |= (VarInt(b & 0x7f) << (7 * num_read))

        num_read++
        if num_read > 5 {
            return ErrorVarIntTooBig
        }

        if (b & 0x80) == 0 {
            break
        }
    }

    *i = tmp
    return nil
}

func (i *VarInt) Write(w io.Writer) error {
    tmp := uint32(*i)

    for {
        b := uint8(tmp & 0x7f)
        tmp >>= 7

        if tmp != 0 {
            b |= 0x80
        }

        err := binary.Write(w, binary.BigEndian, b)
        if err != nil {
            return err
        }

        if tmp == 0 {
            break
        }
    }

    return nil
}

func (l *VarLong) Read(r io.Reader) error {
    var tmp VarLong = 0
    var num_read uint8 = 0

    for {
        var b uint8
        err := binary.Read(r, binary.BigEndian, &b)
        if err != nil {
            return err
        }

        tmp |= (VarLong(b & 0x7f) << (7 * num_read))

        num_read++
        if num_read > 10 {
            return ErrorVarLongTooBig
        }

        if (b & 0x80) == 0 {
            break
        }
    }

    *l = tmp
    return nil
}

func (l *VarLong) Write(w io.Writer) error {
    tmp := uint64(*l)

    for tmp != 0 {
        b := uint8(tmp & 0x7f)
        tmp >>= 7

        if tmp != 0 {
            b |= 0x80
        }

        err := binary.Write(w, binary.BigEndian, b)
        if err != nil {
            return err
        }
    }

    return nil
}

func (s *String) Read(r io.Reader) error {
    size := new(VarInt)
    err := size.Read(r)
    if err != nil {
        return err
    }

    buf := make([]byte, *size)
    _, err = r.Read(buf)
    if err != nil {
        return err
    }

    *s = String(buf)
    return err
}

func (s *String) Write(w io.Writer) error {
    buf := []byte(*s)
    
    size := VarInt(len(buf))
    err := (&size).Write(w)
    if err != nil {
        return err
    }

    _, err = w.Write(buf)

    return err
}

func (p *Position) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, p)
}

func (p *Position) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, p)
}

func (a *Angle) Read(r io.Reader) error {
    return binary.Read(r, binary.BigEndian, a)
}

func (a *Angle) Write(w io.Writer) error {
    return binary.Write(w, binary.BigEndian, a)
}