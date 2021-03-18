package layeraoi

import (
	"aoi"
	"errors"
)

const deBruijn64ctz = 0x0218a392cd3d5dbf

var deBruijnIdx64ctz = [64]byte{
	0, 1, 2, 7, 3, 13, 8, 19,
	4, 25, 14, 28, 9, 34, 20, 40,
	5, 17, 26, 38, 15, 46, 29, 48,
	10, 31, 35, 54, 21, 50, 41, 57,
	63, 6, 12, 18, 24, 27, 33, 39,
	16, 37, 45, 47, 30, 53, 49, 56,
	62, 11, 23, 32, 36, 44, 52, 55,
	61, 22, 43, 51, 60, 42, 59, 58,
}

func Ctz64(x uint64) int {
	x &= -x
	y := x * deBruijn64ctz >> 58
	i := int(deBruijnIdx64ctz[y])
	z := int((x - 1) >> 57 & 64)
	return i + z
}

func SetNZero(x uint64,num int)uint64{
	return x &^ (1 << num)
}

func getLayerBits(obj aoi.IObject) uint64{
	if o,error := convIObject2LayerIObject(obj);error != nil{
		panic(error)
	} else {
		return o.GetLayerBits()
	}
}

func convIObject2LayerIObject(obj aoi.IObject)(ILayerObject,error){
	io,ok := obj.(ILayerObject)
	if !ok {
		return nil,errors.New("Must use ILayerObject ")
	}

	return io,nil
}

func convIWatch2ILayerAOIBase(){

}