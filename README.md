# QMP - quick message protocol

## HandShake

```
 +-------------+                            +-------------+
 |   Client    |        TCP/IP Network      |    Server   |
 +-------------+             |              +-------------+
        |                    |                     |
 Uninitialized               |               Uninitialized
        |           C        |                     |
        |------------------->|         C           |
        |                    |-------------------->|
        |                    |                  Version sent
        |                    |          S          |
        |                    |<--------------------|
        |         S          |                 Handshake Done
        |<-------------------|                     |
   Version sent              |                     |
        |                    |                     |
  Handshake Done             |                     |
        |                    |                     |
```

### HandShake C/S

```
  0 1 2 3 4 5 6 7
 +-+-+-+-+-+-+-+-+
 |    version    |
 +-+-+-+-+-+-+-+-+
```

- Version (8 bits): Allowed version 1, other number 2-255 and 0 not allowed in current time

## Message Format

```
 +-------+---------+------+
 | Setup | Header  | Data |
 +-------+---------+------+
 |                |
 |<-Chunk Header->|
```

### Setup

```
  0              1                2                
  0 1  2 3 4 5 6 7  0 1 2 3 4 5 6 7 0 1
 +-+--+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |type| body lenght | header lenght   |
 +-+--+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

- Type (2 bits): initial notification of the type of subsequent message
- Body length (3 bytes): length of the subsequent message body
- Header length (3 bytes): length of the subsequent message header

| Number | Description    |
|--------|----------------|
| 0      | Empty message  |
| 1      | String message |
| 2      | JSON message   |
| 3      | amf0 message   |

### Header

It represents data in amf0 format. Depending on the previously set header length, it is necessary to decode them.

### Body

The body is decoded according to the specified condition in the installation header
