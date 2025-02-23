package cli

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/errata-ai/vale/v2/internal/core"
)

// Style represents an externally-hosted style.
type Style struct {
	// User-provided fields.
	Author      string `json:"author"`
	Description string `json:"description"`
	Deps        string `json:"deps"`
	Feed        string `json:"feed"`
	Homepage    string `json:"homepage"`
	Name        string `json:"name"`
	URL         string `json:"url"`

	// Generated fields.
	HasUpdate bool `json:"has_update"`
	InLibrary bool `json:"in_library"`
	Installed bool `json:"installed"`
	Addon     bool `json:"addon"`
}

// Meta represents an installed style's meta data.
type Meta struct {
	Author      string   `json:"author"`
	Coverage    float64  `json:"coverage"`
	Description string   `json:"description"`
	Email       string   `json:"email"`
	Feed        string   `json:"feed"`
	Lang        string   `json:"lang"`
	License     string   `json:"license"`
	Name        string   `json:"name"`
	Sources     []string `json:"sources"`
	URL         string   `json:"url"`
	Vale        string   `json:"vale_version"`
	Version     string   `json:"version"`
}

func init() {
	flag.BoolVar(&Flags.Remote, "mode-rev-compat", false,
		"prioritize local Vale configurations")
	flag.StringVar(&Flags.Built, "built", "", "post-processed file path")

	Actions["install"] = install
}

func fetch(src, dst string) error {
	// Fetch the resource from the web:
	resp, err := http.Get(src)
	if err != nil {
		return err
	}

	// Create a temp file to represent the archive locally:
	tmpfile, err := ioutil.TempFile("", "temp.*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Write to the  local archive:
	_, err = io.Copy(tmpfile, resp.Body)
	if err != nil {
		return err
	} else if err = tmpfile.Close(); err != nil {
		return err
	}
	return core.Unzip(tmpfile.Name(), dst)
}

func install(args []string, cfg *core.Config) error {
	style := filepath.Join(cfg.StylesPath, args[0])
	if core.IsDir(style) {
		os.RemoveAll(style) // Remove existing version
	}

	err := fetch(args[1], cfg.StylesPath)
	if err != nil {
		return sendResponse(
			fmt.Sprintf("Failed to install '%s'", args[1]),
			err)
	}

	return sendResponse(fmt.Sprintf(
		"Successfully installed '%s'", args[1]), nil)
}
