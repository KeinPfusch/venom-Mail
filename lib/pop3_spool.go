package lib

import (
	"log"

	"github.com/peterbourgon/diskv"
)

func init() {

	VenomMail := diskv.New(diskv.Options{
		BasePath:     Hpwd() + "/spool/mail",
		Transform:    blockTransform,
		CacheSizeMax: 1024 * 1024, // 1MB
	})

	log.Printf("[LIB][SPOOL][POP3] Spool mailspace created in %s", VenomMail.BasePath)
	mailWrite(VenomMail, "test1")
	log.Printf("[LIB][SPOOL][POP3] Spool test successful: %t  ", mailRead(VenomMail, shasum("test1")) != nil)
	log.Printf("[LIB][SPOOL][POP3] Spool delete successful: %t  ", VenomMail.Erase(shasum("test1")) == nil)

}

func mailWrite(vq *diskv.Diskv, mailin string) {

	vq.Write(shasum(mailin), []byte(mailin))

}

func mailList(vq *diskv.Diskv, mailin string) []string {

	var l []string
	for key := range vq.Keys(nil) {
		l = append(l, key)

	}
	return l

}

func mailRead(vq *diskv.Diskv, key string) []byte {

	val, err := vq.Read(key)
	if err != nil {
		log.Println("%s not in spool", key)
		return nil
	}
	return val

}
