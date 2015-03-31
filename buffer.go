package golis

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//自定义缓存
type Buffer struct {
	b    []byte
	rOff int //读取位置
	wOff int //写入位置
}

//创建Buffer
func NewBuffer() *Buffer {
	return new(Buffer)
}

//获取缓存当前容量
func (b *Buffer) BufferCap() int {
	return cap(b.b)
}

//写入bytes数据
func (b *Buffer) putBytes(buffer []byte) {
	b.b = append(b.b, buffer...)
	b.wOff = len(buffer) + b.wOff
}

//指定位置写入,如果指定写入位置超出了wOff位置,则抛出异常
//如果指定位置已经存在数据并写入数据超出wOff位置则覆盖之前数据，wOff变更最新
//如果指定位置已经存在数据并写入数据没有超出wOff位置则覆盖之前数据，wOff不变
func (b *Buffer) putBytesAt(pos int, buffer []byte) error {
	willPos := pos + len(buffer)
	if pos > b.wOff {
		return errors.New("pos is out of wOff")
	}
	if willPos > b.wOff {
		copy(b.b[b.wOff:], buffer)
		b.wOff = willPos
	} else {
		copy(b.b[b.wOff:], buffer)
	}
	return nil
}

//将int数据存入缓存
func (b *Buffer) putInt(i int) {
	x := int32(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	b.putBytes(bytesBuffer.Bytes())
}

//将uint32数据放入内存
func (b *Buffer) putUint32(i uint32) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, i)
	b.putBytes(bytesBuffer.Bytes())
}

//将字符串存入buffer
func (b *Buffer) putString(s string) {
	b.putBytes([]byte(s))
}

//读取指定位置开始，指定长度的bytes数据
//如果读取数据位置超出了写入数据的位置，则返回错误
func (b *Buffer) ReadBytesAt(pos, length int) ([]byte, error) {
	if pos > b.wOff {
		return nil, errors.New("pos is out of wOff")
	}
	buffer := make([]byte, length)
	if pos+length > b.wOff {
		copy(buffer, b.b[pos:b.wOff])
		b.rOff = b.wOff
	} else {
		p := pos + length
		copy(buffer, b.b[pos:p])
		b.rOff = p
	}

	return buffer, nil
}

//读取int
func (b *Buffer) ReadInt() (int, error) {
	rpos := b.rOff + 4
	if rpos > b.wOff {
		return 0, errors.New("ReadInt out off wOff")
	}
	bytesBuffer := bytes.NewBuffer(b.b[b.rOff:rpos])
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	b.rOff = rpos
	return int(x), nil
}

//读取uint32数据
func (b *Buffer) ReadUint32() (uint32, error) {
	rpos := b.rOff + 4
	if rpos > b.wOff {
		return 0, errors.New("ReadUint32 out off wOff")
	}
	bytesBuffer := bytes.NewBuffer(b.b[b.rOff:rpos])
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	b.rOff = rpos
	return x, nil
}

//读取字符串
func (b *Buffer) ReadString(length int) (string, error) {
	rpos := b.rOff + length
	if rpos > b.wOff {
		return "", errors.New("ReadString out of wOff")
	}
	s := string(b.b[b.rOff:rpos])
	b.rOff = rpos
	return s, nil
}
