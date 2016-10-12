package passwduser

import (
	"bufio"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
)

var passwdFilePath = "/etc/passwd"

// Current finds and returns the current user.
func Current() (*User, error) {
	uid := strconv.Itoa(os.Getuid())
	return LookupID(uid)
}

// Lookup finds a user by her username.
func Lookup(username string) (*User, error) {
	passwdFile, err := os.Open(passwdFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = passwdFile.Close()
	}()

	users, err := parsePasswdFilter(passwdFile, func(u User) bool {
		return u.Username == username
	})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, user.UnknownUserError(username)
	}

	return &users[0], nil
}

// LookupID finds a user by her UID.
func LookupID(uid string) (*User, error) {
	uidInt, err := strconv.ParseInt(uid, 10, 32)
	if err != nil {
		return nil, err
	}

	passwdFile, err := os.Open(passwdFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = passwdFile.Close()
	}()

	users, err := parsePasswdFilter(passwdFile, func(u User) bool {
		return u.UID == uid
	})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, user.UnknownUserIdError(uidInt)
	}

	return &users[0], nil

}

// `(*User) GroupIds() ([]string, error)`
// `LookupGroup(name string) (*Group, error)`
// `LookupGroupId(gid string) (*Group, error)`

func parseLine(line string) User {
	user := User{}

	// see: man 5 passwd
	//  name:password:UID:GID:GECOS:directory:shell
	parts := strings.Split(line, ":")
	if len(parts) >= 1 {
		user.Username = parts[0]
		user.Name = parts[0]
	}
	if len(parts) >= 3 {
		user.UID = parts[2]
	}
	if len(parts) >= 4 {
		user.GID = parts[3]
	}
	if len(parts) >= 6 {
		user.HomeDir = parts[5]
	}

	return user
}

func parsePasswdFilter(r io.Reader, filter func(User) bool) ([]User, error) {
	out := []User{}

	s := bufio.NewScanner(r)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nil, err
		}

		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		p := parseLine(line)
		if filter == nil || filter(p) {
			out = append(out, p)
		}
	}

	return out, nil
}
