package passwduser

import (
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"testing"
)

func TestLookup(t *testing.T) {
	const passwdContent = `
root:x:0:0:root user:/root:/bin/bash
adm:x:42:43:adm:/var/adm:/bin/false
111:x:222:333::/home/111:/bin/false
this is just some garbage data
`

	tests := []struct {
		testDescription string
		username        string
		expected        User
	}{
		{
			testDescription: "RootUser",
			username:        "root",
			expected: User{
				UID:      "0",
				GID:      "0",
				Username: "root",
				Name:     "root",
				HomeDir:  "/root",
			},
		},
		{
			testDescription: "NonRootUser",
			username:        "adm",
			expected: User{
				UID:      "42",
				GID:      "43",
				Username: "adm",
				Name:     "adm",
				HomeDir:  "/var/adm",
			},
		},
		{
			testDescription: "NumericUsernames",
			username:        "111",
			expected: User{
				UID:      "222",
				GID:      "333",
				Username: "111",
				Name:     "111",
				HomeDir:  "/home/111",
			},
		},
	}

	passwdFile, err := writePasswdFile(passwdContent)
	if err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
	passwdFilePath = passwdFile.Name()

	t.Run("Group", func(t *testing.T) {
		for _, test := range tests {
			test := test
			t.Run(test.testDescription, func(t *testing.T) {
				t.Parallel()

				actual, err := Lookup(test.username)
				if err != nil {
					t.Logf(
						"got unexpected error when looking up user '%s': %s",
						test.username, err.Error(),
					)
					t.Fail()
					return
				}

				if !reflect.DeepEqual(test.expected, *actual) {
					t.Logf("username: %v", test.username)
					t.Logf("got:      %#v", actual)
					t.Logf("expected: %#v", test.expected)
					t.Fail()
				}
			})
		}
	})

	if err := os.Remove(passwdFile.Name()); err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
}

func TestLookupID(t *testing.T) {
	const passwdContent = `
root:x:0:0:root user:/root:/bin/bash
adm:x:42:43:adm:/var/adm:/bin/false
111:x:222:333::/home/111:/bin/false
this is just some garbage data
`

	tests := []struct {
		testDescription string
		uid             string
		expected        User
	}{
		{
			testDescription: "RootUser",
			uid:             "0",
			expected: User{
				UID:      "0",
				GID:      "0",
				Username: "root",
				Name:     "root",
				HomeDir:  "/root",
			},
		},
		{
			testDescription: "NonRootUser",
			uid:             "42",
			expected: User{
				UID:      "42",
				GID:      "43",
				Username: "adm",
				Name:     "adm",
				HomeDir:  "/var/adm",
			},
		},
		{
			testDescription: "NumericUsernames",
			uid:             "222",
			expected: User{
				UID:      "222",
				GID:      "333",
				Username: "111",
				Name:     "111",
				HomeDir:  "/home/111",
			},
		},
	}

	passwdFile, err := writePasswdFile(passwdContent)
	if err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
	passwdFilePath = passwdFile.Name()

	t.Run("Group", func(t *testing.T) {
		for _, test := range tests {
			test := test
			t.Run(test.testDescription, func(t *testing.T) {
				t.Parallel()

				actual, err := LookupID(test.uid)
				if err != nil {
					t.Logf(
						"got unexpected error when looking up user id '%s': %s",
						test.uid, err.Error(),
					)
					t.Fail()
					return
				}

				if !reflect.DeepEqual(test.expected, *actual) {
					t.Logf("uid:      %v", test.uid)
					t.Logf("got:      %#v", actual)
					t.Logf("expected: %#v", test.expected)
					t.Fail()
				}
			})
		}
	})

	if err := os.Remove(passwdFile.Name()); err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
}

func TestLookupErrors(t *testing.T) {
	const passwdContent = `
root:x:0:0:root user:/root:/bin/bash
adm:x:42:43:adm:/var/adm:/bin/false
111:x:222:333::/home/111:/bin/false
this is just some garbage data
`

	tests := []struct {
		testDescription string
		username        string
		expectedError   string
	}{
		{
			"NonExistingUsername",
			"test",
			user.UnknownUserError("test").Error(),
		},
		{
			"NonExistingNumbericUsername",
			"222",
			user.UnknownUserError("222").Error(),
		},
	}

	passwdFile, err := writePasswdFile(passwdContent)
	if err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
	passwdFilePath = passwdFile.Name()

	t.Run("Group", func(t *testing.T) {
		for _, test := range tests {
			test := test
			t.Run(test.testDescription, func(t *testing.T) {
				t.Parallel()

				_, err := Lookup(test.username)
				if err == nil {
					t.Logf("expected error, got nil when looking up username '%s'",
						test.username,
					)
					t.Fail()
					return
				}

				actualError := err.Error()
				if actualError != test.expectedError {
					t.Logf("username:				%v", test.username)
					t.Logf("got error:      %v", actualError)
					t.Logf("expected error: %v", test.expectedError)
					t.Fail()
				}
			})
		}
	})

	if err := os.Remove(passwdFile.Name()); err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
}

func TestLookupIDErrors(t *testing.T) {
	const passwdContent = `
root:x:0:0:root user:/root:/bin/bash
adm:x:42:43:adm:/var/adm:/bin/false
111:x:222:333::/home/111:/bin/false
this is just some garbage data
`

	tests := []struct {
		testDescription string
		uid             string
		expectedError   string
	}{
		{
			"NonExistingUsername",
			"-20",
			user.UnknownUserIdError(-20).Error(),
		},
		{
			"NonExistingNumbericUsername",
			"111",
			user.UnknownUserIdError(111).Error(),
		},
	}

	passwdFile, err := writePasswdFile(passwdContent)
	if err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
	passwdFilePath = passwdFile.Name()

	t.Run("Group", func(t *testing.T) {
		for _, test := range tests {
			test := test
			t.Run(test.testDescription, func(t *testing.T) {
				t.Parallel()

				_, err := LookupID(test.uid)
				if err == nil {
					t.Logf("expected error, got nil when looking up user id '%s'",
						test.uid,
					)
					t.Fail()
					return
				}

				actualError := err.Error()
				if actualError != test.expectedError {
					t.Logf("uid:						%v", test.uid)
					t.Logf("got error:      %v", actualError)
					t.Logf("expected error: %v", test.expectedError)
					t.Fail()
				}
			})
		}
	})

	if err := os.Remove(passwdFile.Name()); err != nil {
		t.Logf("got unexpected error: %s", err.Error())
		t.Fatal()
	}
}

func writePasswdFile(passwdContent string) (*os.File, error) {
	passwdFile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	_, err = passwdFile.WriteString(passwdContent)
	if err != nil {
		return nil, err
	}

	if err = passwdFile.Close(); err != nil {
		return nil, err
	}

	return passwdFile, nil
}
