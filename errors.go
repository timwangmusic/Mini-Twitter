package main

import log "github.com/sirupsen/logrus"

func checkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

func checkFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
