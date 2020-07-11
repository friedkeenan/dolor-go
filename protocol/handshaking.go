package protocol

/* ID: 0x00 */
type HandshakeServerbound struct {
    ProtocolVersion VarInt
    ServerAddress String
    ServerPort UnsignedShort
    NextState State
}