package dist

import (
	"encoding/json"
	"sort"
)

// TODO -Dan- add this as inspection output
type Asset struct {
	Sha256      string                 `toml:"sha256" json:"sha256,omitempty"`
	ID          string                 `toml:"id" json:"id"`
	Version     string                 `toml:"version" json:"version"`
	Name        string                 `toml:"name" json:"name,omitempty"`
	LayerDiffID string                 `json:"layerDiffId,omitempty"`
	URI         string                 `toml:"uri" json:"uri,omitempty"`
	Licenses    []string               `toml:"licenses" json:"licenses,omitempty"`
	Description string                 `toml:"description" json:"description,omitempty"`
	Homepage    string                 `toml:"homepage" json:"homepage,omitempty"`
	Stacks      []string               `toml:"stacks" json:"stacks"`
	Metadata    map[string]interface{} `toml:"metadata" json:"metadata,omitempty"`
}

type Assets []Asset

func (a Assets) MarshalJSON() ([]byte, error) {
	m := map[string]Asset{}
	aCopy := make(Assets, len(a))
	copy(aCopy, a)
	for _, asset := range aCopy {
		sha := asset.Sha256
		asset.Sha256 = ""
		m[sha] = asset
	}

	return json.Marshal(m)
}

func (a *Assets) UnmarshalJSON(b []byte) error {
	var m map[string]Asset
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	//
	for sha, asset := range m {
		asset.Sha256 = sha
		*a = append(*a, asset)
	}

	// TODO -Dan- validate ordering
	aSlice := (*a)[:]
	sort.Slice(aSlice, func(i, j int) bool {
		return aSlice[i].ID < aSlice[j].ID
	})

	return nil
}
