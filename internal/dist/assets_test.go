package dist_test

import (
	"bytes"
	"encoding/json"
	"github.com/buildpacks/pack/internal/dist"
	"github.com/buildpacks/pack/testhelpers"
	"github.com/sclevine/spec"
	"strings"
	"testing"
)

func TestAssets(t *testing.T) {
	spec.Run(t, "TestAssets", testAssets)
}

func testAssets(t *testing.T, when spec.G,it  spec.S) {
	var assert = testhelpers.NewAssertionManager(t)
	when("Assets to json", func() {
		var (
			expectedAssetsJson string
			firstAsset dist.Asset
			secondAsset dist.Asset
		)

		it.Before(func(){
			expectedAssetsJson = `{
    "first-asset-sha256": {
        "name": "First Asset",
        "id": "first-asset-id",
        "version": "1.1.1",
        "layerDiffId": "first asset layerDiffID",
        "uri": "first asset uri",
        "licenses": [
            "first asset license"
        ],
        "description": "first asset description",
        "homepage": "first-homepage-link",
        "stacks": [
            "stack1",
            "stack2"
        ],
        "metadata": {
            "cool": "bean"
        }
    },
    "second-asset-sha256": {
        "name": "Second Asset",
        "id": "second-asset-id",
        "version": "2.2.2",
        "layerDiffId": "second asset layerDiffID",
        "uri": "second asset uri",
        "licenses": [
            "second asset license"
        ],
        "description": "second asset description",
        "homepage": "second-homepage-link",
        "stacks": [
            "stack1",
            "stack2"
        ],
        "metadata": {
            "cooler": "bean"
        }
    }
}`
			firstAsset = dist.Asset{
				Sha256:      "first-asset-sha256",
				Name:        "First Asset",
				ID:          "first-asset-id",
				Version:     "1.1.1",
				LayerDiffID: "first asset layerDiffID",
				URI:         "first asset uri",
				Licenses:    []string{"first asset license"},
				Description: "first asset description",
				Homepage:    "first-homepage-link",
				Stacks:      []string{"stack1", "stack2"},
				Metadata:    map[string]interface{}{"cool": "bean"},
			}
			secondAsset = dist.Asset{
				Sha256:      "second-asset-sha256",
				Name:        "Second Asset",
				ID:          "second-asset-id",
				Version:     "2.2.2",
				LayerDiffID: "second asset layerDiffID",
				URI:         "second asset uri",
				Licenses:    []string{"second asset license"},
				Description: "second asset description",
				Homepage:    "second-homepage-link",
				Stacks:      []string{"stack1", "stack2"},
				Metadata:    map[string]interface{}{"cooler": "bean"},
			}
		})
		when("#MarshalJSON", func() {
			it("marshal Assets into a json map of sha256 -> asset", func() {
				assets := dist.Assets{firstAsset, secondAsset}

				buf := bytes.NewBuffer(nil)
				enc := json.NewEncoder(buf)
				enc.SetIndent("","    ")
				err := enc.Encode(assets)
				assert.Nil(err)

				assert.EqualJSON(buf.String(),expectedAssetsJson)
			})
		})
		when("#UnmarshalJSON", func() {
			it("unmarshall map of sha256 -> asset to Assets", func() {
				reader := strings.NewReader(expectedAssetsJson)
				dec := json.NewDecoder(reader)
				var assets dist.Assets
				err := dec.Decode(&assets)
				assert.Nil(err)

				assert.Equal(assets, dist.Assets{firstAsset,secondAsset})
			})
		})
	})
}