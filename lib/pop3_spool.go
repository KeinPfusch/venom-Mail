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
	mailWrite(VenomMail, "test2")
	log.Printf("[LIB][SPOOL][POP3] Spool test successful: %t  ", len(spoolList(VenomMail)) == 2)
	VenomMail.EraseAll()

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
