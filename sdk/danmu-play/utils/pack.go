/*
 * @Author: ww
 * @Date: 2022-07-02 02:27:22
 * @Description:
 * @FilePath: /danmuplay/utils/pack.go
 */
package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
)

// python struct pack/unpack golang 实现
// 改自：https://github.com/roman-kachanovsky/go-binary-pack

type BinaryPack struct{}

func reverse[T string | byte | interface{}](m *[]T) {
	l := len(*m)
	for i := 0; i < l/2; i++ {
		(*m)[i], (*m)[l-1-i] = (*m)[l-1-i], (*m)[i]
	}
}

// Return a byte slice containing the values of msg slice packed according to the given format.
// The items of msg slice must match the values required by the format exactly.
func (bp *BinaryPack) Pack(format []string, msg []interface{}) ([]byte, error) {
	if len(format) > len(msg) {
		return nil, errors.New("Format is longer than values to pack")
	}
	reverse(&format)
	reverse(&msg)
	res := []byte{}

	for i, f := range format {
		switch f {
		case "?":
			casted_value, ok := msg[i].(bool)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (bool)")
			}
			res = append(res, boolToBytes(casted_value)...)
		case "h", "H":
			casted_value, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 2 bytes)")
			}
			res = append(res, intToBytes(casted_value, 2)...)
		case "i", "I", "l", "L":
			casted_value, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 4 bytes)")
			}
			res = append(res, intToBytes(casted_value, 4)...)
		case "q", "Q":
			casted_value, ok := msg[i].(int)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (int, 8 bytes)")
			}
			res = append(res, intToBytes(casted_value, 8)...)
		case "f":
			casted_value, ok := msg[i].(float32)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float32)")
			}
			res = append(res, float32ToBytes(casted_value, 4)...)
		case "d":
			casted_value, ok := msg[i].(float64)
			if !ok {
				return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (float64)")
			}
			res = append(res, float64ToBytes(casted_value, 8)...)
		default:
			if strings.Contains(f, "s") {
				casted_value, ok := msg[i].(string)
				if !ok {
					return nil, errors.New("Type of passed value doesn't match to expected '" + f + "' (string)")
				}
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, []byte(fmt.Sprintf("%s%s",
					casted_value, strings.Repeat("\x00", n-len(casted_value))))...)
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}
	reverse(&res)
	return res, nil
}

// Unpack the byte slice (presumably packed by Pack(format, msg)) according to the given format.
// The result is a []interface{} slice even if it contains exactly one item.
// The byte slice must contain not less the amount of data required by the format
// (len(msg) must more or equal CalcSize(format)).
func (bp *BinaryPack) UnPack(format []string, msg []byte) ([]interface{}, error) {
	expected_size, err := bp.CalcSize(format)

	if err != nil {
		return nil, err
	}

	if expected_size > len(msg) {
		return nil, errors.New("Expected size is bigger than actual size of message")
	}
	reverse(&format)
	reverse(&msg)

	res := []interface{}{}

	for _, f := range format {
		switch f {
		case "?":
			res = append(res, bytesToBool(msg[:1]))
			msg = msg[1:]
		case "h", "H":
			res = append(res, bytesToInt(msg[:2]))
			msg = msg[2:]
		case "i", "I", "l", "L":
			res = append(res, bytesToInt(msg[:4]))
			msg = msg[4:]
		case "q", "Q":
			res = append(res, bytesToInt(msg[:8]))
			msg = msg[8:]
		case "f":
			res = append(res, bytesToFloat32(msg[:4]))
			msg = msg[4:]
		case "d":
			res = append(res, bytesToFloat64(msg[:8]))
			msg = msg[8:]
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, string(msg[:n]))
				msg = msg[n:]
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}
	reverse(&res)
	return res, nil
}

// Return the size of the struct (and hence of the byte slice) corresponding to the given format.
func (bp *BinaryPack) CalcSize(format []string) (int, error) {
	var size int

	for _, f := range format {
		switch f {
		case "?":
			size = size + 1
		case "h", "H":
			size = size + 2
		case "i", "I", "l", "L", "f":
			size = size + 4
		case "q", "Q", "d":
			size = size + 8
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				size = size + n
			} else {
				return 0, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return size, nil
}

func boolToBytes(x bool) []byte {
	if x {
		return intToBytes(1, 1)
	}
	return intToBytes(0, 1)
}

func bytesToBool(b []byte) bool {
	return bytesToInt(b) > 0
}

func intToBytes(n int, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, int64(n))
	return buf.Bytes()[0:size]
}

func bytesToInt(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 2:
		var x int16
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	case 4:
		var x int32
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	default:
		var x int64
		binary.Read(buf, binary.LittleEndian, &x)
		return int(x)
	}
}

func float32ToBytes(n float32, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat32(b []byte) float32 {
	var x float32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &x)
	return x
}

func float64ToBytes(n float64, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, n)
	return buf.Bytes()[0:size]
}

func bytesToFloat64(b []byte) float64 {
	var x float64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &x)
	return x
}

var user_danmaku = map[float64]string{}

func DecodeMessage(message []byte, outMsg chan []byte) {
	packet_len, _ := strconv.ParseUint(fmt.Sprintf("%x", message[:4]), 16, 64)
	ver, _ := strconv.ParseInt(fmt.Sprintf("%x", message[6:8]), 16, 64)
	op, _ := strconv.ParseInt(fmt.Sprintf("%x", message[8:12]), 16, 64)
	for int64(len(message)) > int64(packet_len) {
		DecodeMessage(message[packet_len:], outMsg)
		message = message[:packet_len]
	}
	switch ver {
	case 0:
		outMsg <- message[16:]
		// log.Printf(string(message[16:]))

	case 1:
		if op == 3 {
			popular, _ := strconv.ParseInt(fmt.Sprintf("%x", message[16:]), 16, 64)
			// log.Printf("[人气]  %d", popular)
			popmsg := fmt.Sprintf(`{"cmd":"POP","count": %d}`, popular)
			// fmt.Println(popmsg)
			outMsg <- []byte(popmsg)
		}
	case 2:
		b := bytes.NewReader(message[16:])
		var out bytes.Buffer
		r, _ := zlib.NewReader(b)
		io.Copy(&out, r)
		DecodeMessage(out.Bytes(), outMsg)
	case 3:
		b := bytes.NewReader(message[16:])
		r := brotli.NewReader(b)
		bytess, _ := io.ReadAll(r)
		DecodeMessage(bytess, outMsg)
	default:
		// log.Printf("未知数据包")
		// log.Printf(string(message))
		// log.Printf("packet_len: %d, ver: %d, op: %d", packet_len, ver, op)
	}
}
