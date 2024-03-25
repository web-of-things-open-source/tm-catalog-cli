package repos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/viper"
	"github.com/wot-oss/tmc/internal/config"
	"github.com/wot-oss/tmc/internal/model"
	"github.com/wot-oss/tmc/internal/utils"
)

const (
	KeyRepos       = "repos"
	keyRemotes     = "remotes" // left for compatibility
	KeyRepoType    = "type"
	KeyRepoLoc     = "loc"
	KeyRepoAuth    = "auth"
	KeyRepoEnabled = "enabled"

	RepoTypeFile             = "file"
	RepoTypeHttp             = "http"
	RepoTypeTmc              = "tmc"
	CompletionKindNames      = "names"
	CompletionKindFetchNames = "fetchNames"
	RepoConfDir              = ".tmc"
	IndexFilename            = "tm-catalog.toc.json"
	TmNamesFile              = "tmnames.txt"
)

var ValidRepoNameRegex = regexp.MustCompile("^[a-zA-Z0-9][\\w\\-_:]*$")

type Config map[string]map[string]any

var SupportedTypes = []string{RepoTypeFile, RepoTypeHttp, RepoTypeTmc}

//go:generate mockery --name Repo --outpkg mocks --output mocks
type Repo interface {
	// Push writes the Thing Model file into the path under root that corresponds to id.
	// Returns ErrTMIDConflict if the same file is already stored with a different timestamp or
	// there is a file with the same semantic version and timestamp but different content
	Push(ctx context.Context, id model.TMID, raw []byte) error
	// Fetch retrieves the Thing Model file from repo
	// Returns the actual id of the retrieved Thing Model (it may differ in the timestamp from the id requested), the file contents, and an error
	Fetch(ctx context.Context, id string) (string, []byte, error)
	// Index updates repository's index file with data from given TM files. For ids that refer to non-existing files,
	// removes those from index. Performs a full update if no updatedIds given
	Index(ctx context.Context, updatedIds ...string) error
	// List searches the catalog for TMs matching search parameters
	List(ctx context.Context, search *model.SearchParams) (model.SearchResult, error)
	// Versions lists versions of a TM with given name
	Versions(ctx context.Context, name string) ([]model.FoundVersion, error)
	// Spec returns the spec this Repo has been created from
	Spec() model.RepoSpec
	// Delete deletes the TM with given id from repo. Returns ErrTmNotFound if TM does not exist
	Delete(ctx context.Context, id string) error

	ListCompletions(ctx context.Context, kind string, toComplete string) ([]string, error)
}

var Get = func(spec model.RepoSpec) (Repo, error) {
	if spec.Dir() != "" {
		if spec.RepoName() != "" {
			return nil, model.ErrInvalidSpec
		}
		return NewFileRepo(map[string]any{KeyRepoType: "file", KeyRepoLoc: spec.Dir()}, spec)
	}
	repos, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	repos = filterEnabled(repos)
	rc, ok := repos[spec.RepoName()]
	if spec.RepoName() == "" {
		switch len(repos) {
		case 0:
			return nil, ErrRepoNotFound
		case 1:
			for n, v := range repos {
				rc = v
				spec = model.NewRepoSpec(n)
			}
		default:
			return nil, ErrAmbiguous
		}
	} else {
		if !ok {
			return nil, ErrRepoNotFound
		}
	}
	return createRepo(rc, spec)
}

func filterEnabled(repos Config) Config {
	res := make(Config)
	for n, rc := range repos {
		enabled := utils.JsGetBool(rc, KeyRepoEnabled)
		if enabled != nil && !*enabled {
			continue
		}
		res[n] = rc
	}
	return res
}

func createRepo(rc map[string]any, spec model.RepoSpec) (Repo, error) {
	switch t := rc[KeyRepoType]; t {
	case RepoTypeFile:
		return NewFileRepo(rc, spec)
	case RepoTypeHttp:
		return NewHttpRepo(rc, spec)
	case RepoTypeTmc:
		return NewTmcRepo(rc, spec)
	default:
		return nil, fmt.Errorf("unsupported repo type: %v. Supported types are %v", t, SupportedTypes)
	}
}

var All = func() ([]Repo, error) {
	conf, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	var rs []Repo

	for n, rc := range conf {
		en := utils.JsGetBool(rc, KeyRepoEnabled)
		if en != nil && !*en {
			continue
		}
		r, err := createRepo(rc, model.NewRepoSpec(n))
		if err != nil {
			return rs, err
		}
		rs = append(rs, r)
	}
	return rs, err
}

func ReadConfig() (Config, error) {
	reposConfig := viper.Get(KeyRepos)
	if reposConfig == nil {
		remotesConfig := viper.Get(keyRemotes) // attempt to find obsolete key and convert config to new key
		if remotes, ok := remotesConfig.(map[string]any); ok {
			conf, err := mapToConfig(remotes)
			if err != nil {
				return nil, err
			}
			err = saveConfig(conf)
			if err != nil {
				return nil, err
			}
			reposConfig = remotesConfig
		}
	}
	repos, ok := reposConfig.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid repo config")
	}
	return mapToConfig(repos)
}

func mapToConfig(repos map[string]any) (Config, error) {
	cp := map[string]map[string]any{}
	for k, v := range repos {
		if cfg, ok := v.(map[string]any); ok {
			cp[k] = cfg
		} else {
			return nil, fmt.Errorf("invalid repo config: %s", k)
		}
	}
	return cp, nil
}

func ToggleEnabled(name string) error {
	conf, err := ReadConfig()
	if err != nil {
		return err
	}
	c, ok := conf[name]
	if !ok {
		return ErrRepoNotFound
	}
	if enabled, ok := c[KeyRepoEnabled]; ok {
		if eb, ok := enabled.(bool); ok && !eb {
			delete(c, KeyRepoEnabled)
		} else {
			c[KeyRepoEnabled] = false
		}
	} else {
		c[KeyRepoEnabled] = false
	}
	conf[name] = c
	return saveConfig(conf)
}

func Remove(name string) error {
	conf, err := ReadConfig()
	if err != nil {
		return err
	}
	if _, ok := conf[name]; !ok {
		return ErrRepoNotFound
	}
	delete(conf, name)
	return saveConfig(conf)
}

func Add(name, typ, confStr string, confFile []byte) error {
	_, err := Get(model.NewRepoSpec(name))
	if err == nil || !errors.Is(err, ErrRepoNotFound) {
		return ErrRepoExists
	}

	return setRepoConfig(name, typ, confStr, confFile, err)
}

func SetConfig(name, typ, confStr string, confFile []byte) error {
	_, err := Get(model.NewRepoSpec(name))
	if err != nil && errors.Is(err, ErrRepoNotFound) {
		return ErrRepoNotFound
	}

	return setRepoConfig(name, typ, confStr, confFile, err)
}

func setRepoConfig(name string, typ string, confStr string, confFile []byte, err error) error {
	var rc map[string]any
	switch typ {
	case RepoTypeFile:
		rc, err = createFileRepoConfig(confStr, confFile)
		if err != nil {
			return err
		}
	case RepoTypeHttp:
		rc, err = createHttpRepoConfig(confStr, confFile)
		if err != nil {
			return err
		}
	case RepoTypeTmc:
		rc, err = createTmcRepoConfig(confStr, confFile)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported repo type: %v. Supported types are %v", typ, SupportedTypes)
	}

	conf, err := ReadConfig()
	if err != nil {
		return err
	}

	conf[name] = rc

	return saveConfig(conf)
}

func Rename(oldName, newName string) error {
	if !ValidRepoNameRegex.MatchString(newName) {
		return ErrInvalidRepoName
	}
	conf, err := ReadConfig()
	if err != nil {
		return err
	}
	if rc, ok := conf[oldName]; ok {
		conf[newName] = rc
		delete(conf, oldName)
		return saveConfig(conf)
	} else {
		return ErrRepoNotFound
	}
}
func saveConfig(conf Config) error {
	viper.Set(KeyRepos, conf)
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		configFile = filepath.Join(config.DefaultConfigDir, "config.json")
	}
	err := os.MkdirAll(config.DefaultConfigDir, 0770)
	if err != nil {
		return err
	}

	b, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(b) == 0 {
		b = []byte("{}")
	}
	var j map[string]any
	err = json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	j[KeyRepos] = conf
	delete(j, keyRemotes)

	w, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}
	return utils.AtomicWriteFile(configFile, w, 0660)
}

func AsRepoConfig(bytes []byte) (map[string]any, error) {
	var js any
	err := json.Unmarshal(bytes, &js)
	if err != nil {
		return nil, fmt.Errorf("invalid json config: %w", err)
	}
	rc, ok := js.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid json config. must be a map")
	}
	return rc, nil
}

// GetSpecdOrAll returns a union containing the repo specified by spec, or union of all repo, if the spec is empty
func GetSpecdOrAll(spec model.RepoSpec) (*Union, error) {
	if spec.RepoName() != "" || spec.Dir() != "" {
		repo, err := Get(spec)
		if err != nil {
			return nil, err
		}
		return NewUnion(repo), nil
	}
	all, err := All()
	if err != nil {
		return nil, err
	}
	return NewUnion(all...), nil
}
