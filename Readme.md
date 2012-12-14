# SVN Daemon

## Synopsis
SVN Daemon is a small web based tool that can manipulate a svn checkout. We use it at trivago to
manage some of our dev servers.

## Installation
You will need a go compiler as long as I do not provide any downloads (sorry).

```
# Build it
go build main.go

# Rename it
mv main svn-daemon

# Make it executable
chmod +x svn-daemon

# Run it
./svn-daemon --config /path/to/config
```
## Functions
* Update current checkout (svn up)
* Switch to any other branch or tag (with support for bootstrap's typeahead plugin)
* optionally revert before each action
* four hooks configurable: pre_up, post_up, pre_sw, post_sw

## Screenshot
[[https://dl.dropbox.com/u/230202/projects/svn-daemon/view1.png]]
