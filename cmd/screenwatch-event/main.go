package main

import (
	"flag"
	"fmt"
	"github.com/guelfey/go.dbus"
	"github.com/hdonnay/screenwatch"
	"os"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "Be verbose")
	display string
)

func init() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr, "not enough arguments")
		os.Exit(-1)
	}
	display = strings.Join(strings.Split(flag.Arg(1), "-")[1:], "")
	if *verbose {
		fmt.Fprintf(os.Stderr, "sending signal %s %s\n", flag.Arg(0), display)
	}
}

func main() {
	var conn *dbus.Conn
	var err error
	var call *dbus.Call
	conn, err = dbus.SystemBus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to dbus System Bus: %v\n", err)
		os.Exit(1)
	}
	obj := conn.Object(screenwatch.Name, screenwatch.Path)
	if *verbose {
		c := obj.Call("Ping", 0)
		if c.Err != nil {
			fmt.Fprintln(os.Stderr, "pinged object")
		} else {
			fmt.Fprintln(os.Stderr, "unable to ping object")
		}
	}
	switch flag.Arg(0) {
	case "connect":
		if *verbose {
			fmt.Fprintln(os.Stderr, "sending connect")
		}
		call = obj.Call("Connect", 0, display)
	case "disconnect":
		if *verbose {
			fmt.Fprintln(os.Stderr, "sending disconnect")
		}
		call = obj.Call("Disconnect", 0, display)
	default:
		fmt.Fprintf(os.Stderr, "unrecognized event: %s\n", flag.Arg(0))
	}
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, call.Err)
		os.Exit(1)
	}
	if call.Body[0].(bool) {
		os.Exit(0)
	}
	os.Exit(1)
}
