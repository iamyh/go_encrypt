package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	//"time"
)

type Demo struct {
	Head DemoHead
	Body DemoBody
}

type DemoHead struct {
	Cmd     uint8
	VerLen  uint32
	Version []byte
}

type DemoBody struct {
	Len     uint32
	Content []byte
}

func (demo *Demo) WriteToBuffer(buf *bytes.Buffer) (err error) {

	head := &demo.Head
	err = head.WriteToBuffer(buf)
	if err != nil {
		return
	}

	body := &demo.Body
	err = body.WriteToBuffer(buf)
	if err != nil {
		return
	}

	return
}

func (demo *Demo) ReadFromBuf(buf *bytes.Buffer) (err error) {

	head := &demo.Head
	err = head.ReadFromBuf(buf)
	if err != nil {
		return
	}

	body := &demo.Body
	err = body.ReadFromBuf(buf)
	if err != nil {
		return
	}

	return
}

func (dh *DemoHead) WriteToBuffer(buf *bytes.Buffer) (err error) {
	err = binary.Write(buf, binary.BigEndian, dh.Cmd)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, dh.VerLen)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, dh.Version)
	if err != nil {
		return
	}

	return
}

func (dh *DemoHead) ReadFromBuf(buf *bytes.Buffer) (err error) {
	err = binary.Read(buf, binary.BigEndian, &dh.Cmd)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.BigEndian, &dh.VerLen)
	if err != nil {
		return
	}

	dh.Version = make([]byte, dh.VerLen)

	err = binary.Read(buf, binary.BigEndian, &dh.Version)
	if err != nil {
		return
	}

	return
}

func (db *DemoBody) WriteToBuffer(buf *bytes.Buffer) (err error) {
	err = binary.Write(buf, binary.BigEndian, db.Len)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, db.Content)
	if err != nil {
		return
	}

	return
}

func (db *DemoBody) ReadFromBuf(buf *bytes.Buffer) (err error) {
	err = binary.Read(buf, binary.BigEndian, &db.Len)
	if err != nil {
		return
	}

	db.Content = make([]byte, db.Len)
	err = binary.Read(buf, binary.BigEndian, &db.Content)
	if err != nil {
		return
	}

	return
}

func ZlibEncode(buf *bytes.Buffer) (res []byte, err error) {
	var zlibBuf bytes.Buffer

	zlibw := zlib.NewWriter(&zlibBuf)

	_, err = zlibw.Write(buf.Bytes())
	if err != nil {
		return
	}

	zlibw.Close()
	res = zlibBuf.Bytes()
	return
}

func ZlibDecode(src []byte) (buf *bytes.Buffer, err error) {

	var tmp = bytes.NewBuffer(src)

	zlibr, err := zlib.NewReader(tmp)
	if err != nil {
		return
	}

	var res []byte

	res, err = ioutil.ReadAll(zlibr)
	if err != nil {
		return
	}

	zlibr.Close()
	buf = bytes.NewBuffer(res)
	return
}

//为了凑合一定的字节数
type Variant struct {
	OffsetLen  uint8 //1 byte
	Offset     []byte
	ContentLen uint32 //4 byte
	Content    []byte
	Tail       []byte //3 byte
}

func (v *Variant) WriteToBuffer(buf *bytes.Buffer) (err error) {
	err = binary.Write(buf, binary.BigEndian, v.OffsetLen)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, v.Offset)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, v.ContentLen)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, v.Content)
	if err != nil {
		return
	}

	err = binary.Write(buf, binary.BigEndian, v.Tail)
	if err != nil {
		return
	}

	return
}

func (v *Variant) ReadFromBuf(buf *bytes.Buffer) (err error) {

	err = binary.Read(buf, binary.BigEndian, &v.OffsetLen)
	if err != nil {
		return
	}

	v.Offset = make([]byte, v.OffsetLen)
	err = binary.Read(buf, binary.BigEndian, v.Offset)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.BigEndian, &v.ContentLen)
	if err != nil {
		return
	}

	v.Content = make([]byte, v.ContentLen)
	err = binary.Read(buf, binary.BigEndian, v.Content)
	if err != nil {
		return
	}

	err = binary.Read(buf, binary.BigEndian, v.Tail)
	if err != nil {
		return
	}

	return
}

func Byte2Uint16(src []byte) (dst []uint16, err error) {
	length := len(src)
	if length%2 != 0 {
		err = fmt.Errorf("len %d err", length)
		return
	}

	count := length / 2
	dst = make([]uint16, count)
	for i := 0; i < count; i++ {
		dst[i] = uint16(src[i*2]) | uint16(src[i*2+1])<<8
	}

	return
}

func Uint16ToByte(src []uint16) (dst []byte) {

	length := len(src)

	dst = make([]byte, 2*length)

	for i := 0; i < length; i++ {
		dst[i*2] = byte(src[i])
		dst[i*2+1] = byte(src[i] >> 8)
	}

	return
}

func InjectKey(src []uint16, key []uint16) {

	srcLen := len(src)
	keyLen := len(key)

	for i := 0; i < srcLen-1; i++ {

		var p = i
		if i >= keyLen {
			p = i % keyLen
		}

		delta := src[i+1] & key[p]
		src[i] += delta
	}
}

func RelieveKey(src []uint16, key []uint16) {

	srcLen := len(src)
	keyLen := len(key)

	for i := srcLen - 2; i >= 0; i-- {

		var p = i
		if i >= keyLen {
			p = i % keyLen
		}

		delta := src[i+1] & key[p]
		src[i] -= delta
	}
}
