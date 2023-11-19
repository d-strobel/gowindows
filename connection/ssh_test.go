package connection

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

// Testing fixtures
const (
	errMsgParams string = "ssh client: SSHConfig parameter 'SSHHost', 'SSHUsername' and one of 'SSHPassword', 'SSHPrivateKey', 'SSHPrivateKeyPath' must be set"

	privateKeyED25519 string = `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCNKqLkDaaa4KGp+xaT0X94XVxGiwG6RHsymEc9/m39hwAAAJjpeDkr6Xg5
KwAAAAtzc2gtZWQyNTUxOQAAACCNKqLkDaaa4KGp+xaT0X94XVxGiwG6RHsymEc9/m39hw
AAAEAMT15+Ut2N+m9HW9wXgIeVR+qKeoT3UlVCxxnPsnoA5o0qouQNpprgoan7FpPRf3hd
XEaLAbpEezKYRz3+bf2HAAAAD2RzdHJvYmVsQE5CMDc4NAECAwQFBg==
-----END OPENSSH PRIVATE KEY-----
`

	privateKeyRSA string = `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`
)

func TestNewSSHClient(t *testing.T) {
	t.Run("NoParameters", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})

	t.Run("Case1", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHHost:     "test.example.intern",
			SSHUsername: "root",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
	t.Run("Case2", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHHost:     "test.example.intern",
			SSHPassword: "Test",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
	t.Run("Case3", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHHost:       "test.example.intern",
			SSHPrivateKey: "Test",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
	t.Run("Case4", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHHost:           "test.example.intern",
			SSHPrivateKeyPath: "Test",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
	t.Run("Case5", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHUsername: "root",
			SSHPassword: "Test",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
	t.Run("Case6", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHUsername:   "root",
			SSHPrivateKey: "Test",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
	t.Run("Case7", func(t *testing.T) {
		// Test case with missing required fields
		config := &SSHConfig{
			SSHUsername:       "root",
			SSHPrivateKeyPath: "Test",
		}

		client, err := newSSHClient(config)
		if err == nil {
			t.Error("Expected an error, but got nil")
		}
		if client != nil {
			t.Error("Expected a nil SSH client, but got non-nil")
		}

		assert.Contains(t, err.Error(), errMsgParams)
	})
}

func TestKnownHostCallback(t *testing.T) {
	t.Run("IgnoreHostKey", func(t *testing.T) {
		config := &SSHConfig{
			SSHInsecureIgnoreHostKey: true,
		}
		callback, err := knownHostCallback(config)
		assert.NoError(t, err)
		assert.NotNil(t, callback)
	})
}

func TestAuthenticationMethod(t *testing.T) {

	// Merge test cases here
	var testCases []string
	testCases = append(testCases, privateKeyRSA, privateKeyED25519)

	tempDir := t.TempDir()

	t.Run("PrivateKey", func(t *testing.T) {
		for _, key := range testCases {
			config := &SSHConfig{
				SSHPrivateKey: key,
			}

			authMethod, err := authenticationMethod(config)
			assert.NoError(t, err)
			assert.Len(t, authMethod, 1)
			assert.IsType(t, []ssh.AuthMethod{}, authMethod)
			assert.IsType(t, ssh.PublicKeysCallback(nil), authMethod[0])
		}
	})

	t.Run("PrivateKeyFromFile", func(t *testing.T) {
		for i, key := range testCases {

			// Write test key to temporary file
			file, err := os.CreateTemp(tempDir, fmt.Sprintf("test_file_privatekey_%d", i))
			assert.NoError(t, err)
			file.WriteString(key)
			t.Cleanup(func() {
				file.Close()
			})

			config := &SSHConfig{
				SSHPrivateKeyPath: file.Name(),
			}

			authMethod, err := authenticationMethod(config)
			assert.NoError(t, err)
			assert.Len(t, authMethod, 1)
			assert.IsType(t, []ssh.AuthMethod{}, authMethod)
			assert.IsType(t, ssh.PublicKeysCallback(nil), authMethod[0])
		}
	})

	t.Run("Password", func(t *testing.T) {
		config := &SSHConfig{
			SSHPassword: "your_password_here",
		}

		authMethod, err := authenticationMethod(config)
		assert.NoError(t, err)
		assert.Len(t, authMethod, 1)
		assert.IsType(t, []ssh.AuthMethod{}, authMethod)
		assert.IsType(t, ssh.PasswordCallback(nil), authMethod[0])
	})

	t.Run("PrivateKeyAndPassword", func(t *testing.T) {
		for _, key := range testCases {
			config := &SSHConfig{
				SSHPrivateKey: key,
				SSHPassword:   "test_password",
			}

			authMethod, err := authenticationMethod(config)
			assert.NoError(t, err)
			assert.Len(t, authMethod, 2)
			assert.IsType(t, []ssh.AuthMethod{}, authMethod)
			assert.IsType(t, ssh.PublicKeysCallback(nil), authMethod[0])
			assert.IsType(t, ssh.PasswordCallback(nil), authMethod[1])
		}
	})

	t.Run("PrivateKeyAndFileAndPassword", func(t *testing.T) {
		for i, key := range testCases {

			// Write test key to temporary file
			file, err := os.CreateTemp(tempDir, fmt.Sprintf("test_file_privatekey_%d", i))
			assert.NoError(t, err)
			file.WriteString(key)
			t.Cleanup(func() {
				file.Close()
			})

			config := &SSHConfig{
				SSHPrivateKey:     key,
				SSHPrivateKeyPath: file.Name(),
				SSHPassword:       "test_password",
			}

			authMethod, err := authenticationMethod(config)
			assert.NoError(t, err)
			assert.Len(t, authMethod, 2)
			assert.IsType(t, []ssh.AuthMethod{}, authMethod)
			assert.IsType(t, ssh.PublicKeysCallback(nil), authMethod[0])
			assert.IsType(t, ssh.PasswordCallback(nil), authMethod[1])
		}
	})
}
