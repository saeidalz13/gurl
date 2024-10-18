package errutils

import (
	"log"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
