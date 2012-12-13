# SVN Daemon

## Synopsis
SVN Daemon is a small tool that can manipulate a svn checkout. We use it at trivago to
manage some of our dev servers.

## Installation
You will need a go compiler as long as I do not provide any downloads (sorry).

```
cd src
go build
mv src svn-daemon
chmod +x svn-daemon
./svn-daemon --config /path/to/config
```