#!/bin/sh -x
go get -u github.com/hdonnay/screenwatch/cmd/screenwatchd
go get -u github.com/hdonnay/screenwatch/cmd/screenwatch-event

cat << EOF | sudo sh
set -x
install policy.conf /etc/dbus-1/system.d/screenwatch.conf
install 99-xrandr-trigger.rules /etc/udev/rules.d
install $GOPATH/src/github.com/hdonnay/screenwatch/cmd/screenwatch-event/screenwatch-event /lib/udev
install $GOPATH/src/github.com/hdonnay/screenwatch/cmd/screenwatchd/screenwatchd /usr/local/bin
EOF
