// Package updater implements `proxyx update`: it checks GitHub for the latest
// release, downloads the archive for the current OS/arch, verifies its SHA-256
// against the published checksums, and atomically replaces the running binary.
package updater

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Repo is the GitHub repository releases are pulled from.
const Repo = "otabek05/ProxyX"

const (
	apiLatest = "https://api.github.com/repos/" + Repo + "/releases/latest"
	userAgent = "proxyx-self-update" // generic UA; carries no user info
	maxBytes  = 64 << 20             // 64 MiB cap on any download
)

type asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

type release struct {
	Tag    string  `json:"tag_name"`
	Assets []asset `json:"assets"`
}

// Result describes what Update did.
type Result struct {
	Updated bool
	From    string
	To      string
}

func client() *http.Client { return &http.Client{Timeout: 2 * time.Minute} }

func get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client().Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("GET %s: HTTP %d", url, res.StatusCode)
	}
	return res, nil
}

func latest() (*release, error) {
	res, err := get(apiLatest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var r release
	if err := json.NewDecoder(io.LimitReader(res.Body, maxBytes)).Decode(&r); err != nil {
		return nil, fmt.Errorf("parsing release info: %w", err)
	}
	if r.Tag == "" {
		return nil, fmt.Errorf("no latest release found for %s", Repo)
	}
	return &r, nil
}

func download(url string) ([]byte, error) {
	res, err := get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(io.LimitReader(res.Body, maxBytes))
}

// NormalizeVersion strips a leading "v" and surrounding space.
func NormalizeVersion(v string) string {
	return strings.TrimPrefix(strings.TrimSpace(v), "v")
}

// BinaryName is the executable name inside the archive for an OS.
func BinaryName(goos string) string {
	if goos == "windows" {
		return "proxyx.exe"
	}
	return "proxyx"
}

// AssetName is the archive asset name for a version and platform, matching the
// GoReleaser naming (proxyx_<version>_<os>_<arch>.<ext>).
func AssetName(version, goos, goarch string) string {
	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("proxyx_%s_%s_%s.%s", version, goos, goarch, ext)
}

// ChecksumsName is the checksums asset name for a version.
func ChecksumsName(version string) string {
	return fmt.Sprintf("proxyx_%s_checksums.txt", version)
}

func findAsset(r *release, name string) (string, bool) {
	for _, a := range r.Assets {
		if a.Name == name {
			return a.URL, true
		}
	}
	return "", false
}

// VerifyChecksum checks data's SHA-256 against the entry for name in a
// GoReleaser checksums file ("<hex>  <filename>" per line).
func VerifyChecksum(data []byte, name, checksums string) error {
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	for _, line := range strings.Split(checksums, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == name {
			if strings.EqualFold(fields[0], got) {
				return nil
			}
			return fmt.Errorf("checksum mismatch for %s", name)
		}
	}
	return fmt.Errorf("no checksum listed for %s", name)
}

// extractBinary pulls the named binary out of a .tar.gz archive.
func extractBinary(targz []byte, binName string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(string(targz)))
	if err != nil {
		return nil, fmt.Errorf("reading archive: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading archive: %w", err)
		}
		if hdr.Typeflag == tar.TypeReg && filepath.Base(hdr.Name) == binName {
			return io.ReadAll(io.LimitReader(tr, maxBytes))
		}
	}
	return nil, fmt.Errorf("%s not found in archive", binName)
}

// replaceExecutable atomically replaces the file at path with data, preserving
// an executable mode. It writes a temp file in the same directory and renames.
func replaceExecutable(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".proxyx-update-*")
	if err != nil {
		return fmt.Errorf("cannot write to %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op after a successful rename
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replacing %s: %w", path, err)
	}
	return nil
}

// Update checks for a newer release and, if found, installs it in place.
// currentVersion is the running binary's version. out receives progress lines.
func Update(currentVersion string, out io.Writer) (*Result, error) {
	logf := func(format string, a ...any) { fmt.Fprintf(out, format+"\n", a...) }

	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("self-update is not supported on Windows; download the latest .zip from https://github.com/%s/releases/latest", Repo)
	}

	logf("🔎 Checking for the latest release…")
	rel, err := latest()
	if err != nil {
		return nil, err
	}
	latestVer := NormalizeVersion(rel.Tag)
	curVer := NormalizeVersion(currentVersion)

	if curVer == latestVer {
		return &Result{Updated: false, From: curVer, To: latestVer}, nil
	}
	logf("⬇️  Updating %s → %s", curVer, latestVer)

	assetName := AssetName(latestVer, runtime.GOOS, runtime.GOARCH)
	assetURL, ok := findAsset(rel, assetName)
	if !ok {
		return nil, fmt.Errorf("release %s has no asset %s", rel.Tag, assetName)
	}
	csURL, ok := findAsset(rel, ChecksumsName(latestVer))
	if !ok {
		return nil, fmt.Errorf("release %s has no checksums file", rel.Tag)
	}

	logf("   downloading %s", assetName)
	archive, err := download(assetURL)
	if err != nil {
		return nil, err
	}
	checksums, err := download(csURL)
	if err != nil {
		return nil, err
	}

	logf("   verifying checksum")
	if err := VerifyChecksum(archive, assetName, string(checksums)); err != nil {
		return nil, err
	}

	bin, err := extractBinary(archive, BinaryName(runtime.GOOS))
	if err != nil {
		return nil, err
	}

	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}

	logf("   installing to %s", exe)
	if err := replaceExecutable(exe, bin); err != nil {
		return nil, err
	}
	return &Result{Updated: true, From: curVer, To: latestVer}, nil
}
