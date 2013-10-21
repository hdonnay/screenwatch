package main

import (
	"flag"
	"fmt"
	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/introspect"
	"github.com/hdonnay/screenwatch"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	verbose        = flag.Bool("v", false, "Be (very) verbose")
	logFile        = flag.String("l", "", "File to log to. If empty, log to stderr")
	daemonize      = flag.Bool("d", false, "Daemonize after starting")
	positionalArgs = flag.String("p", "DP1:--right-of LVDS1", "Comma separated list of <display>:<position> pairs")

	conn     *dbus.Conn
	l        *log.Logger
	xrandr   string
	s        swatch
	position = make(map[string][]string)
)

type swatch struct {
	introspect.Introspectable
	pos map[string][]string
}

func (s swatch) Connect(display string) (bool, *dbus.Error) {
	args := []string{"--output", display, "--auto"}
	if _, exist := s.pos[display]; exist {
		for _, a := range s.pos[display] {
			args = append(args, a)
		}
	}
	cmd := exec.Command(xrandr, args...)
	if *verbose {
		l.Printf("connecting %s\n", display)
		l.Println(cmd.Args)
	}
	cmd.Run()
	return cmd.ProcessState.Success(), nil
}

func (s swatch) Disconnect(display string) (bool, *dbus.Error) {
	cmd := exec.Command(xrandr, "--output", display, "--off")
	cmd.Run()
	if *verbose {
		l.Printf("disconnecting %s\n", display)
		l.Println(cmd.Args)
	}
	return cmd.ProcessState.Success(), nil
}

func init() {
	var err error
	flag.Parse()
	if *logFile != "" {
		f, err := os.Open(*logFile)
		defer f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l = log.New(f, "", log.LstdFlags)
	} else {
		l = log.New(os.Stderr, "screenwatchd ", log.LstdFlags)
	}

	conn, err = dbus.SystemBus()
	if err != nil {
		l.Fatalln(err)
	}

	xrandr, err = exec.LookPath("xrandr")
	if err != nil {
		l.Fatalf("couldn't find xrandr: %v\n", err)
	}

	if *verbose {
		l.Printf("parsing '%s'\n", *positionalArgs)
	}
	s = swatch{introspect.NewIntrospectable(screenwatch.Introspect), make(map[string][]string)}
	for _, a := range strings.Split(*positionalArgs, ",") {
		arg := strings.Split(a, ":")
		if _, exist := s.pos[arg[0]]; exist {
			l.Printf("duplicate display argument '%s:%s', discarding\n", arg[0], arg[1])
			continue
		}
		s.pos[arg[0]] = strings.Fields(arg[1])
	}
	if *verbose {
		fmt.Printf("positional arguments: %v\n", s.pos)
	}

	if *daemonize {
		l.Println("FIXME: I would daemonize here, but I'm too dumb right now.")
	}
}

func main() {
	var err error
	err = screenwatch.Export(s, conn)
	if err != nil {
		l.Fatalln("unable to setup d-bus", err)
	}
	if *verbose {
		l.Printf("Listening on %s / %s\n", screenwatch.Path, screenwatch.Name)
	}

	// TODO: clean up nicely when sent SIGTERM et. al.
	select {
	//<-stop
	//teardown()
	//break
	}
	os.Exit(0)
}
