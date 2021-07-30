// Package DemoParse is a simple and easy library used to read data on the bit-level from a buffer.
package bitbuffer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// BitBuffer represents a buffer, which is filled with bytes where each bit can be read as a single unit.
type BitBuffer struct {
	buffer    []byte
	bitPos    uint8
	bytePos   int32
	ByteOrder binary.ByteOrder
}

// TooManyBitsError occurs when attempting to read too many bits from a buffer.
type TooManyBitsError uint8

func (err TooManyBitsError) Error() string {
	return fmt.Sprintf("DemoParse: too many bits requested: %v", err)
}

// NewBitBuffer constructs a new BitBuffer.
func NewBitBuffer(byteOrder binary.ByteOrder) (bitBuffer *BitBuffer) {
	return &BitBuffer{
		bitPos:    0,
		bytePos:   0,
		ByteOrder: byteOrder,
	}
}

// Feed data bytes into the buffer.
func (bitBuffer *BitBuffer) Feed(data []byte) {
	bitBuffer.buffer = append(bitBuffer.buffer, data...)
	//bitBuffer.buffer = append(bitBuffer.buffer, []byte)
	return
}

// Return cur bitPos.
func (bitbuffer *BitBuffer) BitPos() uint8 {
	return bitbuffer.bitPos
}

// Current Byte Position
func (bitbuffer *BitBuffer) BytePos() int32 {
	return bitbuffer.bytePos
}

func (bitbuffer *BitBuffer) Size() int {
	return len(bitbuffer.buffer)
}


// Seek
// whence : 0 -> from head
//		  : 1 -> from cur
func (bitBuffer *BitBuffer) Seek(pos int32, whence uint8) {
	if whence == 0 {
		bitBuffer.bytePos = pos
	} else if whence == 1 {
		bitBuffer.bytePos += pos
	}
}


// Clear buffer.
func (bitBuffer *BitBuffer) Clear() {
	bitBuffer.buffer = []byte{}
	bitBuffer.bitPos = 0
	bitBuffer.bytePos = 0
}

// Read a number of bits from the buffer and return them as a byte array.
func (bitBuffer *BitBuffer) Read(numBits uint64) (data []byte, err error) {
	if uint64((len(bitBuffer.buffer) - int(bitBuffer.bytePos))*8-int(bitBuffer.bitPos)) < numBits {
		err = io.EOF
	}

	for numBits > 0 && len(bitBuffer.buffer) - int(bitBuffer.bytePos) > 0 {
		data = append(data, bitBuffer.buffer[bitBuffer.bytePos])
		data[len(data)-1] <<= bitBuffer.bitPos

		if len(bitBuffer.buffer) - int(bitBuffer.bytePos) > 1 {
			shifter := bitBuffer.buffer[bitBuffer.bytePos + 1] >> (8 - bitBuffer.bitPos)
			data[len(data)-1] ^= shifter
		}

		if numBits < 8 {
			data[len(data)-1] >>= (8 - numBits)
			data[len(data)-1] <<= (8 - numBits)

			if uint64(bitBuffer.bitPos)+numBits > 7 {
				bitBuffer.bytePos++
				//bitBuffer.buffer = bitBuffer.buffer[1:]
			}

			bitBuffer.bitPos = uint8((uint64(bitBuffer.bitPos) + numBits) % 8)
			numBits = 0

			return
		}

		//bitBuffer.buffer = bitBuffer.buffer[1:]
		bitBuffer.bytePos++
		numBits -= 8
	}

	return
}

// ReadUint64 reads a uint64 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint64(numBits uint8) (data uint64, err error) {
	if numBits > 64 {
		err = TooManyBitsError(numBits)

		return
	}

	dataBytes, err := bitBuffer.Read(uint64(numBits))

	if err != nil {
		return
	}

	if bitBuffer.ByteOrder == binary.BigEndian {
		dataBytes = append(make([]byte, 8-len(dataBytes)), dataBytes...)
	} else if bitBuffer.ByteOrder == binary.LittleEndian {
		dataBytes = append(dataBytes, make([]byte, 8-len(dataBytes))...)
	}

	err = binary.Read(bytes.NewBuffer(dataBytes), bitBuffer.ByteOrder, &data)

	if err != nil {
		return
	}

	shifter := uint8(0)

	if numBits%8 > 0 {
		shifter = 8 - (numBits % 8)
	}

	data >>= shifter

	return
}

// ReadUint a uint from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint(numBits uint8) (data uint, err error) {
	wordSize := 32 << (^uint(0) >> 32 & 1)

	if numBits > uint8(wordSize) {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint(rawData)

	return
}

// ReadBit reads a single bit as a boolean and returns the value.
func (bitBuffer *BitBuffer) ReadBit() (data bool, err error) {
	rawData, err := bitBuffer.ReadUint8(1)

	if err != nil {
		return
	}

	data = rawData != 0

	return
}

// ReadUint8 reads a uint8 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint8(numBits uint8) (data uint8, err error) {
	if numBits > 8 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint8(rawData)

	return
}

// ReadUint16 reads a uint16 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint16(numBits uint8) (data uint16, err error) {
	if numBits > 16 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint16(rawData)

	return
}

// ReadUint32 reads a uint32 from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadUint32(numBits uint8) (data uint32, err error) {
	if numBits > 32 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = uint32(rawData)

	return
}

func (bitBuffer *BitBuffer) ReadInt32(numBits uint8) (data int32, err error) {
	if numBits > 32 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = int32(rawData)
	return
}



// ReadByte reads a byte from the buffer of numBits size and return the integer value.
func (bitBuffer *BitBuffer) ReadByte(numBytes uint8) (data byte, err error) {
	numBits := numBytes * 8
	if numBits > 8 {
		err = TooManyBitsError(numBits)
	}

	rawData, err := bitBuffer.ReadUint64(numBits)

	if err != nil {
		return
	}

	data = byte(rawData)

	return
}

// ReadString reads numBits bits from the buffer and returns the value as a string.
func (bitBuffer *BitBuffer) ReadString(numBytes uint64) (data string, err error) {
	dataBytes, err := bitBuffer.Read(numBytes * 8)

	if err != nil {
		return
	}

	data = string(dataBytes)

	return
}

func (bitBuffer *BitBuffer) ReadFloat32() (data float32, err error) {
	rawData, err := bitBuffer.ReadUint32(32)
	if err != nil {
		return
	}
	data = math.Float32frombits(rawData)
	return
}

func (bitBuffer *BitBuffer) ReadFloat64() (data float64, err error) {
	rawData, err := bitBuffer.ReadUint64(64)
	if err != nil {
		return
	}
	data = math.Float64frombits(rawData)
	return
}