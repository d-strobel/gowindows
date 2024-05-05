package parsing

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/text/encoding/unicode"
)

// EncodePwshCmd encodes a powershell command to be executed on a Windows machine.
// It encodes the command to UTF-16-LE and then to base64.
// It returns a valid powershell.exe command with no profile and the encoded command.
func EncodePwshCmd(cmd string) (string, error) {
	// Disable unnecessary progress bars which is considered as stderr.
	cmd = fmt.Sprintf("$ProgressPreference = 'SilentlyContinue'; %s", cmd)

	// Encode string to UTF16-LE.
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	encoded, err := encoder.String(cmd)
	if err != nil {
		return "", fmt.Errorf("parsing.EncodePwshCmd: %w", err)
	}

	// Finally make it base64 encoded which is required for powershell.
	cmd = base64.StdEncoding.EncodeToString([]byte(encoded))

	// Specify powershell.exe to run encoded command.
	return fmt.Sprintf("powershell.exe -NoProfile -EncodedCommand %s", cmd), nil
}
