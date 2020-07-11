package main

import (
    "log"

    pr "./protocol"
)

func main() {
    s, err := pr.NewBasicServer(":25565")
    if err != nil {
        log.Fatal(err)
    }

    for {
        c, err := s.NewConnection()
        if err != nil {
            log.Fatal(err)
        }

        go s.HandleConnection(c)
    }
}