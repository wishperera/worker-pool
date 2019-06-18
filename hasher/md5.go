package hasher

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"errors"
)

type MD5 struct {

}

//returns an integer representation of the hashed key using md5 algorithm
func (s MD5)Hash(ctx context.Context,key string)(value int64,err error){
	h := md5.New()
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

