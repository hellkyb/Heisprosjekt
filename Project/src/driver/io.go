package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func Init() bool{
	return bool(int(C.io_init()) == 1)
}

func SetBit(channel int){
	C.io_set_bit(C.int(channel))
}

func ClearBit(channel int){
	C.io_clear_bit(C.int(channel))
}

func WriteAnalog(channel int, value int){
	C.io_write_analog(C.int(channel), C.int(value))
}

func ReadBit(channel int) int{
	return int(C.io_read_bit(C.int(channel)))
}

func ReadAnalog(channel int) int{
	return int(C.io_read_analog(C.int(channel)))
}