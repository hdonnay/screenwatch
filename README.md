screenwatch
===========

Screenwatch is a small utility to send drm change events over the system DBus
and a daemon to receive the and run `xrandr`.

There is a udev rules file, a DBus policy file, a script to install
everything in the proper place, and an XDG autostart file.

screenwatch-event
-----------------

screenwatch-event just calls the proper method.

screenwatchd
------------

screenwatchd establishes a connection to the system bus, and waits for Change()
calls. On change, it looks at all the displays it's configured for (via `-p`)
and makes `xrandr` calls depending on the status read out of `/sys/class/drm`.

It is intended to be run as a normal user.

Installation
------------

    % ./install.sh
