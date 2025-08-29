package dev

import (
	"os"
)

const slaveAddr int8 = 0x25

type Device struct {
	f *os.File
}

// Write 對 I²C slave 寫資料
func (d *Device) Write(b []byte) (int, error) {
	return d.f.Write(b)
}

// Read 從 I²C slave 讀資料
func (d *Device) Read(buf []byte) (int, error) {
	return d.f.Read(buf)
}

// Close 關閉裝置
func (d *Device) Close() error {
	return d.f.Close()
}
