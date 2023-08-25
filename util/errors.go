package util

import log "github.com/sirupsen/logrus"

func CheckErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

func CheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
