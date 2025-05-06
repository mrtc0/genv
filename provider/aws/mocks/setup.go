package mocks

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func InitSessionTestEnv(t *testing.T) {
	StashEnv(t)
	t.Setenv("AWS_CONFIG_FILE", "file_not_exists")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", "file_not_exists")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
}

func SsoTestSetup(t *testing.T, ssoKey string) (err error) {
	t.Helper()

	dir := t.TempDir()

	cacheDir := filepath.Join(dir, ".aws", "sso", "cache")
	err = os.MkdirAll(cacheDir, 0750)
	if err != nil {
		return err
	}

	hash := sha1.New()
	if _, err := hash.Write([]byte(ssoKey)); err != nil {
		return err
	}

	cacheFilename := strings.ToLower(hex.EncodeToString(hash.Sum(nil))) + ".json"

	tokenFile, err := os.Create(filepath.Join(cacheDir, cacheFilename))
	if err != nil {
		return err
	}

	defer func() {
		closeErr := tokenFile.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("close error: %v, original error: %w", closeErr, err)
		}
	}()

	if _, err = tokenFile.WriteString(
		fmt.Sprintf(
			`{"accessToken": "ssoAccessToken", "expiresAt": "%s"}`,
			time.Now().Add(15*time.Minute).Format(time.RFC3339),
		),
	); err != nil {
		return err
	}

	t.Setenv("HOME", dir)

	return nil
}

func StashEnv(t *testing.T) {
	env := os.Environ()
	os.Clearenv()

	t.Cleanup(func() {
		PopEnv(env)
	})
}

func PopEnv(env []string) {
	os.Clearenv()

	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		k, v := p[0], ""
		if len(p) > 1 {
			v = p[1]
		}
		os.Setenv(k, v)
	}
}
