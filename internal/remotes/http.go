package remotes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/web-of-things-open-source/tm-catalog-cli/internal/model"
)

var ErrNotSupported = errors.New("method not supported")

type HttpRemote struct {
	root *url.URL
	name string
}

func NewHttpRemote(config map[string]any, name string) (*HttpRemote, error) {
	loc := config[KeyRemoteLoc]
	locString, ok := loc.(string)
	if !ok {
		return nil, fmt.Errorf("invalid http remote config. loc is either not found or not a string: %v", loc)
	}
	u, err := url.Parse(locString)
	if err != nil {
		return nil, err
	}
	return &HttpRemote{root: u, name: name}, nil
}

func (h *HttpRemote) Push(id model.TMID, raw []byte) error {
	return ErrNotSupported
}

func (h *HttpRemote) Fetch(id model.TMID) ([]byte, error) {
	reqUrl := h.root.JoinPath(id.String())
	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func (h *HttpRemote) CreateToC() error {
	return ErrNotSupported
}
func (h *HttpRemote) Name() string {
	return h.name
}

func (h *HttpRemote) List(filter string) (model.TOC, error) {
	reqUrl := h.root.JoinPath(TOCFilename)
	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return model.TOC{}, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.TOC{}, err
	}

	var toc model.TOC
	err = json.Unmarshal(data, &toc)
	toc.Filter(filter)
	if err != nil {
		return model.TOC{}, err
	}
	return toc, nil
}

func (h *HttpRemote) Versions(name string) (model.TOCEntry, error) {
	log := slog.Default()
	if len(name) == 0 {
		log.Error("Please specify a name to show the TM.")
		return model.TOCEntry{}, errors.New("please specify a name to show the TM")
	}
	toc, err := h.List("")
	if err != nil {
		return model.TOCEntry{}, err
	}
	name = strings.TrimSpace(name)

	tocThing := toc.FindByName(name)
	if tocThing == nil {
		msg := fmt.Sprintf("No thing model found for name: %s", name)
		log.Error(msg)
		return model.TOCEntry{}, errors.New(msg)
	}

	return *tocThing, nil
}

func createHttpRemoteConfig(root string, bytes []byte) (map[string]any, error) {
	if root != "" {
		return map[string]any{
			KeyRemoteType: RemoteTypeHttp,
			KeyRemoteLoc:  root,
		}, nil
	} else {
		rc, err := AsRemoteConfig(bytes)
		if err != nil {
			return nil, err
		}
		if rType, ok := rc[KeyRemoteType]; ok {
			if rType != RemoteTypeHttp {
				return nil, fmt.Errorf("invalid json config. type must be \"http\" or absent")
			}
		}
		rc[KeyRemoteType] = RemoteTypeHttp
		l, ok := rc[KeyRemoteLoc]
		if !ok {
			return nil, fmt.Errorf("invalid json config. must have key \"loc\"")
		}
		ls, ok := l.(string)
		if !ok {
			return nil, fmt.Errorf("invalid json config. loc must be a string")
		}
		rc[KeyRemoteLoc] = ls
		return rc, nil
	}
}