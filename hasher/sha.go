package hasher

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"github.com/go-errors/errors"
)

type SHA256 struct {

}

type SHA512 struct {

}

//returns an integer representation of the hashed key using sha256 algorithm
func (s SHA256)Hash(ctx context.Context,key string)(value int64,err error){
	h := sha256.New()
	_,err = h.Write([]byte(key))
	if err != nil{
		return 0,err
	}

	v := h.Sum(nil)
	vi, n := binary.Varint(v)
	if !(n>0){
		return 0,errors.New("error hashing, couldn't convert byte array to int")
	}

	return vi,nil
}


//returns an integer representation of the hashed key using sha512 algorithm
func (s SHA512)Hash(ctx context.Context,key string)(value int64,err error) {
	h := sha512.New()
	_,err = h.Write([]byte(key))
	if err != nil{
		return 0,err
	}

	v := h.Sum(nil)
	vi, n := binary.Varint(v)
	if !(n>0){
		return 0,errors.New("error hashing, couldn't convert byte array to int")
	}

	return vi,nil
}




