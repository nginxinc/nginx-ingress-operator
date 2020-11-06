package nginxingresscontroller

import (
	"fmt"
	"io/ioutil"
	"os"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const crdsPath = "/kic_crds"
const decoderBufferSize = 100

func getCRDsManifests() ([]string, error) {
	files, err := ioutil.ReadDir(crdsPath)
	if err != nil {
		return nil, err
	}

	var manifests []string
	for _, f := range files {
		manifests = append(manifests, fmt.Sprintf("%v/%v", crdsPath, f.Name()))
	}

	return manifests, nil
}

func kicCRDs() ([]*v1.CustomResourceDefinition, error) {
	manifests, err := getCRDsManifests()
	if err != nil {
		return nil, err
	}

	var crds []*v1.CustomResourceDefinition
	for _, path := range manifests {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open the CRD manifest %v: %v", path, err)
		}

		var crd v1.CustomResourceDefinition

		err = yaml.NewYAMLOrJSONDecoder(f, decoderBufferSize).Decode(&crd)

		if err != nil {
			return nil, fmt.Errorf("failed to parse the CRD manifest %v: %v", path, err)
		}

		err = f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close the CRD manifest %v: %v", path, err)
		}

		crds = append(crds, &crd)
	}

	return crds, nil
}
