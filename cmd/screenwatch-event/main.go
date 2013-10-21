package main

import (
	"flag"
	//"fmt"
	"github.com/guelfey/go.dbus"
	"github.com/hdonnay/screenwatch"
	"log"
	"log/syslog"
	"os"
	//"strings"
)

var (
	verbose = flag.Bool("v", false, "Be verbose")
	display string
	l       *log.Logger
)

func init() {
	var err error
	flag.Parse()
	l, err = syslog.NewLogger(syslog.LOG_USER|syslog.LOG_WARNING, 0)
	if err != nil {
		os.Exit(-1)
	}
	if *verbose {
		l.Printf("sending signal\n")
	}
}

func main() {
	var conn *dbus.Conn
	var err error
	var call *dbus.Call
	conn, err = dbus.SystemBus()
	if err != nil {
		l.Fatalf("Error connecting to dbus System Bus: %v\n", err)
	}
	obj := conn.Object(screenwatch.Name, screenwatch.Path)
	if *verbose {
		c := obj.Call("Ping", 0)
		if c.Err != nil {
			l.Println("pinged object")
		} else {
			l.Println("unable to ping object")
		}
	}
	call = obj.Call("Change", 0)
	if call.Err != nil {
		l.Fatalln(call.Err)
	}
	if call.Body[0].(bool) {
		os.Exit(0)
	}
	os.Exit(1)
}
