package goredis

import "bytes"

// redis sdshdr
type sdshdr struct {
	len  int
	free int
	data []byte
}

// new
// redis sdsnew
func newsds(s string) *sdshdr {
	data := []byte(s)
	return &sdshdr{
		len:  len(data),
		free: 0,
		data: data,
	}
}

// 复制,底层数组copy
// redis sdsdup
func (sds *sdshdr) dup() *sdshdr {
	data := make([]byte, len(sds.data))
	copy(data, sds.data)
	return &sdshdr{
		len:  sds.len,
		free: sds.free,
		data: data,
	}
}

func (sds *sdshdr) clear() {

}

// 拼接一个字符串
// redis sdscat
func (sds *sdshdr) catString(s string) {
	data := []byte(s)
	sds.catSliceByte(data)
}

// 拼接一个byte slice
// 为减少多次拼接造成多次内存分配，使用预留空间
// 拼接后实际长度大于等于1M，则预留1M的空间
// 拼接后实际长度<1M，则预留实际长度的空间
func (sds *sdshdr) catSliceByte(data []byte) {
	if sds.free >= len(data) {
		sds.data = append(sds.data, data...)
	} else { // free不够
		newlen := sds.len + len(data)
		var growlen int
		// 如果新长度newlen>=1M,则最终长度growlen为新长度newlen+1M
		// 否则的话，最终长度growlen=新长度newlen*2
		if newlen >= 1>>20 {
			growlen = newlen + 1>>20
		} else {
			growlen = newlen * 2
		}
		newdata := make([]byte, 0, growlen)
		newdata = append(newdata, sds.data...)
		newdata = append(newdata, data...)
		sds.data = newdata
		sds.len = newlen
		sds.free = growlen - newlen
	}
}

// 将dst的内容，增加到sds上
// redis sdscatsds
func (sds *sdshdr) catsds(dst *sdshdr) {
	data := dst.data
	sds.catSliceByte(data)
}

func (sds *sdshdr) compare(dst *sdshdr) bool {
	if sds.len != dst.len {
		return false
	}
	result := bytes.Compare(sds.data, dst.data)
	return result == 0
}

// 去掉两端的s
// redis sdstrim
func (sds *sdshdr) trim(s string) {
	newdata := bytes.Trim(sds.data, s)
	var newlen int

	if len(newdata) >= 1>>20 {
		newlen = len(newdata) + 1>>20
	} else {
		newlen = newlen * 2
	}
	data := make([]byte, 0, newlen)
	data = append(data, newdata...)
	sds.len = len(newdata)
	sds.free = newlen - len(newdata)
	sds.data = data
}

// 覆盖
// redis sdscpy

func (sds *sdshdr) set(s string) {
	data := []byte(s)
	var newlen int
	if len(data) >= 1>>20 {
		newlen = len(data) + 1>>20
	} else {
		newlen = len(data) * 2
	}
	sds.len = len(data)
	sds.free = newlen - len(data)
	newdata := make([]byte, newlen)
	newdata = append(newdata, data...)
	sds.data = newdata
}

// redis sdsrange
func (sds *sdshdr) setrange() {

}

func (sds *sdshdr) tostring() string {
	data := sds.data[:sds.len]
	return string(data)
}
