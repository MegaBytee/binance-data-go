package binance_data

import (
	"archive/zip"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const outDir = "binance_data" // output base directory

func localPathForURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	// extract prefix=... if present
	if q := parsed.RawQuery; q != "" {
		for _, part := range strings.Split(q, "&") {
			if strings.HasPrefix(part, "prefix=") {
				val := strings.TrimPrefix(part, "prefix=")
				val = strings.TrimPrefix(val, "data/")
				return filepath.Join(outDir, filepath.FromSlash(val)), nil
			}
		}
	}
	// fallback to path
	p := strings.TrimPrefix(parsed.Path, "/")
	p = strings.TrimPrefix(p, "data/")
	return filepath.Join(outDir, filepath.FromSlash(p)), nil
}

func ensureDir(p string) error { return os.MkdirAll(p, 0o755) }

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func extractZipFile(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	destDir := strings.TrimSuffix(zipPath, ".zip")
	if err := ensureDir(destDir); err != nil {
		return err
	}
	for _, f := range r.File {
		fp := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			ensureDir(fp)
			continue
		}
		if err := ensureDir(filepath.Dir(fp)); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(fp)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return err
		}
	}
	fmt.Printf("Extracted: %s -> %s/\n", zipPath, destDir)
	return nil
}
