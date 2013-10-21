// +build linux
package screenwatch

import (
	"fmt"
	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/introspect"
	//"os"
)

const (
	Name      = "com.github.hdonnay.Watch"
	Path      = "/com/github/hdonnay/Watch"
	Signature = "s"
	header    = `<!DOCTYPE node PUBLIC
 "-//freedesktop//DTD D-BUS Object Introspection 1.0//EN"
 "http://www.freedesktop.org/standards/dbus/1.0/introspect.dtd">`
)

var (
	Introspect *introspect.Node
	Connect    string
	Disconnect string
)

type Watch interface {
	Introspect() (string, *dbus.Error)
	Connect(string) (bool, *dbus.Error)
	Disconnect(string) (bool, *dbus.Error)
}

func init() {
	Introspect = &introspect.Node{
		Name: Path,
		Interfaces: []introspect.Interface{
			introspect.Interface{
				Name: Name,
				Methods: []introspect.Method{
					introspect.Method{
						Name: "Connect",
						Args: []introspect.Arg{
							introspect.Arg{"display", "s", "in"},
							introspect.Arg{"exit", "b", "out"},
						},
					},
					introspect.Method{
						Name: "Disconnect",
						Args: []introspect.Arg{
							introspect.Arg{"display", "s", "in"},
							introspect.Arg{"exit", "b", "out"},
						},
					},
				},
			},
		},
	}
	Connect = fmt.Sprint(Name, ".Connect")
	Disconnect = fmt.Sprint(Name, ".Disconnect")
}

func Export(s Watch, conn *dbus.Conn) error {
	var err error
	reply, err := conn.RequestName(Name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("name already taken: %v")
	}

	err = conn.Export(s, Path, Name)
	if err != nil {
		return err
	}
	err = conn.Export(s, Path, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}
	return nil
}
