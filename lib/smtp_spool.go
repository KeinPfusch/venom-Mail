package lib

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"

	"github.com/peterbourgon/diskv"
)

var VenomQueue = diskv.New(diskv.Options{
	BasePath:     Hpwd() + "/spool/queue",
	Transform:    blockTransform,
	CacheSizeMax: 1024 * 1024, // 1MB
})

func init() {

	log.Printf("[LIB][SPOOL][SMTP] Spool queue created in %s", VenomQueue.BasePath)
	SpoolWrite(VenomQueue, "test1")
	log.Printf("[LIB][SPOOL][SMTP] Spool read test successful: %t  ", SpoolRead(VenomQueue, shasum("test1")) != nil)
	log.Printf("[LIB][SPOOL][SMTP] Spool delete successful: %t  ", VenomQueue.Erase(shasum("test1")) == nil)

}

func blockTransform(s string) []string {
	var (
		sliceSize = len(s) / transformBlockSize
		pathSlice = make([]string, sliceSize)
	)
	for i := 0; i < sliceSize; i++ {
		from, to := i*transformBlockSize, (i*transformBlockSize)+transformBlockSize
		pathSlice[i] = s[from:to]
	}
	return pathSlice
}

const transformBlockSize = 32 // grouping of chars per directory depth

func SpoolWrite(vq *diskv.Diskv, mailin string) {

	vq.Write(shasum(mailin), []byte(mailin))

}

func SpoolList(vq *diskv.Diskv) []string {

	var l []string
	for key := range vq.Keys(nil) {
		l = append(l, key)

	}
	return l

}

func SpoolRead(vq *diskv.Diskv, key string) []byte {

	val, err := vq.Read(key)
	if err != nil {
		log.Println("%s not in spool", key)
		return nil
	}
	return val

}

func shasum(s string) string {
	h := sha256.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
