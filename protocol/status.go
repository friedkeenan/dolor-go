package protocol

type NameIdPair struct {
    Name string `json:"name"`
    Id string `json:"id"`
}

type ResponseJson struct {
    Version struct {
        Name string `json:"name"`
        Protocol int `json:"protocol"`
    } `json:"version"`

    Players struct {
        Max int `json:"max"`
        Online int `json:"online"`
        Sample []NameIdPair `json:"sample"`
    } `json:"players"`

    Description map[string]interface{} `json:"description"` // TODO: Make Chat type

    Favicon string `json:"favicon,omitempty"`
}

/* ID: 0x00 */
type RequestServerbound struct {}

/* ID: 0x01 */
type PingServerbound struct {
    Payload Long
}

/* ID: 0x00 */
type ResponseClientbound struct {
    Response ResponseJson `as:"json"`
}

/* ID: 0x01 */
type PongClientbound struct {
    Payload Long
}