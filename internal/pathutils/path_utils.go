package pathutils

import (
	"os"
	"path/filepath"

	"github.com/saeidalz13/gurl/internal/errutils"
)

func MustMakeIpCacheDir() string {
	homeDir, err := os.UserHomeDir()
	errutils.CheckErr(err)

	ipCacheDir := filepath.Join(homeDir, ".gurl", "ipcache")

	os.MkdirAll(ipCacheDir, 0o600)

	return ipCacheDir
}
