# A cgo-free implementation of os/user that parses /etc/passwd

Golang's `os/user` is a great package that provides an API to access
information for the current user/group. However, the Unix implementation relies
in cgo since it uses `getpwuid` [1] to obtain the user/group information. This
library is avoiding that by parsing the `/etc/passwd` file instead.

**NOTICE** this may not be enough for your application. `getpwuid` does also
use LDAP and NIS databases to obtain the user information. If your system uses
one of these technologies, this implementation is not for you and you
unfortunately have to use cgo.

## Work in progress

This package is currently only implementing the following methods:

* `Current() (*User, error)`
* `Lookup(username string) (*User, error)`
* `LookupId(uid string) (*User, error)`

It does not yet implement:

* `(*User) GroupIds() ([]string, error)`
* `LookupGroup(name string) (*Group, error)`
* `LookupGroupId(gid string) (*Group, error)`

[1] https://linux.die.net/man/3/getpwuid
