package passwduser

// User represents a user account.
type User struct {
	UID      string // user ID
	GID      string // primary group ID
	Username string
	Name     string
	HomeDir  string
}
