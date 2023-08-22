package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/xmirrorsecurity/opensca-cli/util/args"
	"github.com/xmirrorsecurity/opensca-cli/util/enum/language"
	"github.com/xmirrorsecurity/opensca-cli/util/logs"
	"github.com/xmirrorsecurity/opensca-cli/util/model"
)

var cacheDir string

func init() {
	// create cache dir
	cacheDir = ".cache"
	if pwd, err := os.Executable(); err == nil {
		pwd = path.Dir(strings.ReplaceAll(pwd, `\`, `/`))
		cacheDir = path.Join(pwd, ".cache")
	}
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		logs.Error(err)
	}
}

// save save cache file
func save(filepath string, data []byte) {
	if args.Config.Cache {
		if err := os.MkdirAll(path.Join(cacheDir, path.Dir(filepath)), 0755); err == nil {
			if f, err := os.Create(path.Join(cacheDir, filepath)); err == nil {
				defer f.Close()
				f.Write(data)
			}
		}
	}
}

// load load cache file
func load(filepath string) []byte {
	if args.Config.Cache {
		if data, err := ioutil.ReadFile(path.Join(cacheDir, filepath)); err == nil {
			return data
		} else {
			return nil
		}
	}
	return []byte{}
}

func filepath(dep model.Dependency) string {
	switch dep.Language {
	case language.Java:
		return path.Join("maven", dep.Vendor, dep.Name, dep.Version.Org, fmt.Sprintf("%s-%s.pom", dep.Name, dep.Version.Org))
	case language.JavaScript:
		return path.Join("npm", fmt.Sprintf("%s.json", dep.Name))
	case language.Php:
		return path.Join("composer", fmt.Sprintf("%s.json", dep.Name))
	default:
		return path.Join("none", fmt.Sprintf("%s-%s-%s", dep.Vendor, dep.Name, dep.Version.Org))
	}
}

// SaveCache save cache file
func SaveCache(dep model.Dependency, data []byte) {
	save(filepath(dep), data)
}

// LoadCache load cache file
func LoadCache(dep model.Dependency) []byte {
	return load(filepath(dep))
}
