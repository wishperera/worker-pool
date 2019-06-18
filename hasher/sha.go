package hasher

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"github.com/go-errors/errors"
)

type SHA256 struct {

}

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


