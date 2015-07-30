package main

import (
	"bytes"
	"fmt"
	"testing"
	//"time"
)

func Test_encrypt(t *testing.T) {

	demo := &Demo{}
	demoHead := DemoHead{
		Version: []byte("1.1.1"),
		Cmd:     uint8(1),
	}
	demoHead.VerLen = uint32(len(demoHead.Version))

	demo.Head = demoHead

	demobody := DemoBody{
		Content: []byte("helloworld"),
	}

	demobody.Len = uint32(len(demobody.Content))

	demo.Body = demobody

	fmt.Printf("demo0:%+v\n", demo)

	var buf bytes.Buffer

	err := demo.WriteToBuffer(&buf)
	if err != nil {
		t.Errorf("err:%s", err)
		return
	}

	res, err := ZlibEncode(&buf)
	if err != nil {
		t.Errorf("err in zlib encode:%s", err)
		return
	}

	offset := 4 - len(res)%4
	variant := &Variant{
		OffsetLen:  uint8(offset),
		Offset:     make([]byte, offset),
		ContentLen: uint32(len(res)),
		Content:    res,
		Tail:       make([]byte, 3),
	}

	var buf1 bytes.Buffer
	err = variant.WriteToBuffer(&buf1)
	if err != nil {
		t.Errorf("variant write to buffer err:%s", err)
		return
	}

	uint16_res, err := Byte2Uint16(buf1.Bytes())
	if err != nil {
		t.Errorf("byte to uint16 err:%s", err)
		return
	}

	encryptKey := "worldhello"
	uint16_key, err := Byte2Uint16([]byte(encryptKey))
	if err != nil {
		t.Errorf("byte to uint16 err:%s", err)
		return
	}

	InjectKey(uint16_res, uint16_key)

	byte_res := Uint16ToByte(uint16_res)

	fmt.Printf("after encrypt,result:%s\n", string(byte_res))

	fmt.Println("start to decrypt")

	uint16_res, err = Byte2Uint16(byte_res)
	if err != nil {
		t.Errorf("Byte2Uint16 err:%s", err)
		return
	}

	decryptKey := "worldhello"
	uint16_decryptKey, err := Byte2Uint16([]byte(decryptKey))
	if err != nil {
		t.Errorf("byte to uint16 err:%s", err)
		return
	}

	RelieveKey(uint16_res, uint16_decryptKey)
	byte_res = Uint16ToByte(uint16_res)
	var buf2 = bytes.NewBuffer(byte_res)

	var variant1 = &Variant{}
	err = variant1.ReadFromBuf(buf2)
	if err != nil {
		t.Errorf("variant read from buf err:%s", err)
		return
	}

	buf3, err := ZlibDecode(variant1.Content)
	if err != nil {
		t.Errorf("err in zlib decode:%s", err)
		return
	}

	var demo1 = &Demo{}
	err = demo1.ReadFromBuf(buf3)
	if err != nil {
		t.Errorf("err:%s", err)
	}

	fmt.Printf("demo1:%+v\n", demo1)

}
