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
	"os/signal"
	"strings"
)

var (
	verbose        = flag.Bool("v", false, "Be (very) verbose")
	logFile        = flag.String("l", "", "File to log to. If empty, log to stderr")
	positionalArgs = flag.String("p", "card0-DP-1:--right-of LVDS1", "Comma separated list of <sys name>:<position> pairs")

	conn   *dbus.Conn
	l      *log.Logger
	xrandr string
	s      swatch
)

type swatch struct {
	introspect.Introspectable
	pos   map[string][]string
	state map[string]bool
}

func (s swatch) Change() (bool, *dbus.Error) {
	for display, _ := range s.pos {
		var err error
		var stat bool
		if stat, err = s.status(display); err != nil {
			return false, &dbus.Error{Name: "READ_ERR"}
		}
		switch stat {
		case true:
			args := []string{"--output", mangle(display), "--auto"}
			if _, exist := s.pos[display]; exist {
				for _, a := range s.pos[display] {
					args = append(args, a)
				}
			}
			cmd := exec.Command(xrandr, args...)
			l.Printf("connecting %s\n", display)
			if *verbose {
				l.Println(cmd.Args)
			}
			cmd.Run()
			return cmd.ProcessState.Success(), nil
		case false:
			cmd := exec.Command(xrandr, "--output", mangle(display), "--off")
			l.Printf("disconnecting %s\n", display)
			if *verbose {
				l.Println(cmd.Args)
			}
			cmd.Run()
			return cmd.ProcessState.Success(), nil
		}
	}
	return false, nil
}

// return changes in connected, disconnected
func (s swatch) status(dev string) (state bool, err error) {
	var status string
	f, err := os.Open(fmt.Sprint("/sys/class/drm/", dev, "/status"))
	if err != nil {
		return
	}
	_, err = fmt.Fscanf(f, "%s\n", &status)
	if err != nil {
		return
	}
	switch status {
	case "connected":
		state = true
	case "disconnected":
		state = false
	default:
		break
	}
	return
}

func mangle(s string) string {
	return strings.Join(strings.Split(s, "-")[1:], "")
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
	s = swatch{
		introspect.NewIntrospectable(screenwatch.Introspect),
		make(map[string][]string),
		make(map[string]bool),
	}
	for _, a := range strings.Split(*positionalArgs, ",") {
		arg := strings.Split(a, ":")
		if _, exist := s.pos[arg[0]]; exist {
			l.Printf("duplicate display argument '%s:%s', discarding\n", arg[0], arg[1])
			continue
		}
		if *verbose {
			l.Printf("watching for %s(%s)\n", arg[0], mangle(arg[0]))
		}
		s.pos[arg[0]] = strings.Fields(arg[1])
	}
	for d, _ := range s.pos {
		var status string
		f, err := os.Open(fmt.Sprint("/sys/class/drm/", d, "/status"))
		if err != nil {
			l.Fatalln(err)
		}
		_, err = fmt.Fscanf(f, "%s\n", &status)
		if err != nil {
			l.Fatalln(err)
		}
		switch status {
		case "connected":
			s.state[d] = true
		case "disconnected":
			s.state[d] = false
		default:
			l.Fatalf("unknown status '%s', bailing\n", status)
		}
	}
	if *verbose {
		fmt.Printf("positional arguments: %v\n", s.pos)
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
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// TODO: clean up nicely when sent SIGTERM et. al.
	for {
		<-stop
		//teardown()
		os.Exit(0)
	}
}
