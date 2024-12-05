package exec

import (
	"os"
	ose "os/exec"
	"path/filepath"
	"runtime"

	"github.com/jolt9dev/jolt9/internal/fs"
	"github.com/jolt9dev/jolt9/pkg/os/env"
	"github.com/jolt9dev/jolt9/pkg/strings"
)

var (
	whichCache = make(map[string]string)
)

type WhichOptions struct {
	UseCache     bool
	PrependPaths []string
}

func Which(command string) (string, bool) {
	return WhichFirst(command, nil)
}

func WhichFirst(command string, options *WhichOptions) (string, bool) {
	if command == "" {
		return "", false
	}

	if options == nil {
		options = &WhichOptions{}
	}

	base := filepath.Base(command)
	ext := filepath.Ext(command)
	name := base[0 : len(base)-len(ext)]
	if options.UseCache {
		path, ok := whichCache[name]
		if ok {
			return path, true
		}
	}

	if filepath.IsAbs(command) {
		fi, err := os.Lstat(command)

		if err != nil {
			return "", false
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			path, err := ose.LookPath(command)
			if err != nil {
				return "", false
			}

			if options.UseCache {
				whichCache[name] = path
			}

			return path, true
		}
	}

	pathSegments := []string{}
	if len(options.PrependPaths) > 0 {
		pathSegments = append(pathSegments, options.PrependPaths...)
	}

	pathSegments = append(pathSegments, env.SplitPath()...)

	for i, path := range pathSegments {
		value := env.ExpandSafe(path)
		if value == "" {
			continue
		}

		pathSegments[i] = value
	}

	for _, path := range pathSegments {
		if strings.IsEmptySpace(path) || !fs.Exists(path) {
			continue
		}

		if runtime.GOOS == "windows" {
			pathExt := env.Get("PATHEXT")
			if strings.IsEmptySpace(pathExt) {
				pathExt = ".com;.exe;.bat;.cmd;.vbs;.vbe;.js;.jse;.wsf;.wsh"
			} else {
				pathExt = strings.ToLower(pathExt)
			}

			extSegments := strings.Split(pathExt, ";")

			entries, err := os.ReadDir(path)
			if err != nil {
				// TODO: debug/trace this erro
				continue
			}

			hasExt := false
			for _, n := range extSegments {
				if strings.EqualFold(n, ext) {
					hasExt = true
					break
				}
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				if hasExt {
					if strings.EqualFold(entry.Name(), command) {
						fp := filepath.Join(path, entry.Name())
						whichCache[name] = fp
						return fp, true
					}

					continue
				}

				entryName := entry.Name()
				entryExt := filepath.Ext(entryName)
				for _, n := range extSegments {
					if strings.EqualFold(n, entryExt) {
						fp := filepath.Join(path, entryName)
						whichCache[name] = fp
						return fp, true
					}
				}
			}
		} else {
			entries, err := os.ReadDir(path)
			if err != nil {
				// TODO: debug/trace this erro
				continue
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				if strings.EqualFold(entry.Name(), name) {
					fp := filepath.Join(path, entry.Name())
					whichCache[name] = fp
					return fp, true
				}
			}
		}
	}

	return "", false
}
