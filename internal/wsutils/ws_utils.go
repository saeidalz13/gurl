package wsutils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const (
	ProtocolWS = iota
	ProtocolWSS
)

const (
	frameFirstByte125OrLess = iota + 125
	frameFirstByte126
	frameFirstByte127
)

/*
Parses the WebSocker frame read from TCP
connection. Strips the extra bytes and
returns the payload.

This is the structure of a websocket frame (each data that comes)
  - FIN (1 bit): Indicates if this is the final fragment of a message.
  - Opcode (4 bits): Specifies the type of data (e.g., text, binary).
  - Mask (1 bit): Indicates if the payload is masked. If 1, the payload is masked.
  - Payload length (7 bits): Specifies the length of the payload data.
    If the length is 0-125, it directly represents the payload length.
    If it is 126, the next 2 bytes represent the payload length.
    If it is 127, the next 8 bytes represent the payload length.
  - Masking key (4 bytes, if Mask is 1): A key used to unmask the payload.
  - Payload data: The actual message data (text or binary).
*/
func ParseWsFrame(frame []byte) ([]byte, error) {
	// The first two bytes are needed to determine
	// the FIN bit, Opcode, Mask bit, and Payload length
	// This is mandatory for a standard websocket frame.
	if len(frame) < 2 {
		return nil, fmt.Errorf("ws frames must be at least 2 bytes")
	}

	// AND operation on 128 (0x80) and second byte of frame.
	// Checks the most significant bit. If that's 1, it means
	// that the payload is masked (needs a key to unmask)
	isMasked := (frame[1] & 0x80) != 0

	// AND operation on 127 (0x7F) and second byte of frame.
	// Extracts the lower 7 bits of the second byte which
	// shows the length of the payload.
	//
	// if 0-125 => It is itself the length of payload
	// if 126 => the next 2 bytes represent the length
	// if 127 => the next 8 bytes represent the length
	payloadLen := int(frame[1] & 0x7F)

	// Declaring an offset to trace our position
	// in frame slice of bytes.
	offset := 2

	switch payloadLen {
	case frameFirstByte125OrLess:
		// No action required

	case frameFirstByte126:
		// Reads the next 2 bytes and interpret them as
		// a 16-bit integer
		//
		// alternative solution:
		// int(frame[2])<<8 | int(frame[3])
		if len(frame) < 4 {
			return nil, fmt.Errorf("frame too short for extended payload length")
		}
		payloadLen = int(binary.BigEndian.Uint16(frame[2:4]))
		offset += 2

	case frameFirstByte127:
		// Reads the next 8 bytes and interpret them as
		// a 64-bit integer
		if len(frame) < 10 {
			return nil, fmt.Errorf("frame too short for 64-bit payload length")
		}

		payloadLen = int(binary.BigEndian.Uint64(frame[2:10]))
		offset += 2
	}

	// Now we have the payload length, we need to make sure
	// frame is long enough to have the payload!
	if payloadLen+offset > len(frame) {
		return nil, fmt.Errorf("frame too short to read payload")
	}

	// we want to extract masking ket ([]byte) if
	// the payload was masked. It should exist in
	// the next 4 bytes after payload length.
	var maskKey []byte
	if isMasked {
		if len(frame) < offset+4 {
			return nil, fmt.Errorf("frame too short to contain mask key")
		}

		maskKey = frame[offset : offset+4]
		offset += 4
	}

	// Now we're at the position of payload
	payload := frame[offset : offset+payloadLen]

	// If isMasked, then we use the key to unmask
	// unmasking is done through XORing each byte
	// of payload with a corresponding byte from
	// the 4-byte mask key. If the length of payload
	// is more, the key is used in a cyclic manner.
	if isMasked {
		for i := 0; i < payloadLen; i++ {
			// ^ is XOR sign
			payload[i] ^= maskKey[i%4]
		}
	}

	return payload, nil
}

func CreateWsFrame(payload []byte) []byte {
	payloadLen := len(payload)
	if payloadLen > 125 {
		fmt.Println("for now it's too long of a payload")
		os.Exit(1)
	}

	var frame bytes.Buffer
	// FIN + RSV (1 2 3) + OPCODE (1 is for text)
	// TODO: Adding feature for other types of messages
	firstByte := byte(0b10000001)
	frame.WriteByte(firstByte)

	// MASK + len of payload
	secondByte := byte(0b10000000 | payloadLen)
	// masking
	maskingKey := []byte{0x37, 0xfa, 0x21, 0x3d}
	for i := range payload {
		payload[i] ^= maskingKey[i%4]
	}
	frame.WriteByte(secondByte)

	// Add the masking key to frame
	frame.Write(maskingKey)

	// Append the payload to the end of the frame
	frame.Write(payload)

	return frame.Bytes()
}
