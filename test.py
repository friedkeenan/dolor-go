#!/usr/bin/env python3

import socket
import json
import code

s = socket.socket()
s.connect(("localhost", 25565))

h = b"\x00\x00\x01h\x63\xdd\x01"
h = bytes([len(h)]) + h

r = b"\x01\x00"

code.interact("", local=locals())

s.close()