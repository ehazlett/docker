package registry

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Api struct {
	host    string
	version string
	secure  bool
}

type ImageInfo struct {
	Image          string
	Revision       int64
	ParentRevision int64
	TarballUrl     string
	State          string
}

func New(host string, version string, secure bool) *Api {
	if version == "" {
		version = "v1"
	}

	return &Api{
		host:    host,
		version: version,
		secure:  secure,
	}
}

func (api *Api) requestUrl(pth string) (*url.URL, error) {
	u, err := url.Parse("/" + api.version + pth)
	if err != nil {
		return nil, err
	}
	u.Host = api.host
	if api.secure {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	return u, nil
}

func (api *Api) ListBranches(user string) ([]string, error) {
	var resp *http.Response
	user = path.Clean(user)
	u, err := api.requestUrl("/users/" + user + "/images")
	if err != nil {
		return nil, err
	}
	if resp, err = http.Get(u.String()); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("Got HTTP error code: " + resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var lst []string
	if err = dec.Decode(&lst); err != nil && err != io.EOF {
		return nil, err
	}
	return lst, nil
}

func (api *Api) ListImages(user, branch string) ([]*ImageInfo, error) {
	var resp *http.Response
	user = path.Clean(user)
	branch = path.Clean(branch)
	u, err := api.requestUrl("/users/" + user + "/images/" + branch)
	if err != nil {
		return nil, err
	}
	if resp, err = http.Get(u.String()); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("Got HTTP error code: " + resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	var lst []*ImageInfo
	if err = dec.Decode(&lst); err != nil && err != io.EOF {
		return nil, err
	}
	return lst, nil
}

func (api *Api) ImageInfo(user, branch, id string) (*ImageInfo, error) {
	var resp *http.Response
	user = path.Clean(user)
	branch = path.Clean(branch)
	id = url.QueryEscape(id)
	u, err := api.requestUrl("/users/" + user + "/images/" + branch)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("id", id)
	u.RawQuery = q.Encode()

	if resp, err = http.Get(u.String()); err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("Got HTTP error code: " + resp.Status)
	}
	dec := json.NewDecoder(resp.Body)
	var img *ImageInfo
	if err = dec.Decode(&img); err != nil && err != io.EOF {
		return nil, err
	}
	return img, nil
}
