package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-logr/logr"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apixv1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	crdsPath          = "./config/crd/kic"
	decoderBufferSize = 100
)

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

func createKICCustomResourceDefinitions(log logr.Logger, mgr manager.Manager) error {
	// Create CRDs with a different client (apiextensions)
	apixClient, err := apixv1client.NewForConfig(mgr.GetConfig())
	if err != nil {
		log.Error(err, "unable to create client for CRD registration")
		return err
	}

	crds, err := kicCRDs()
	if err != nil {
		return err
	}

	crdsClient := apixClient.CustomResourceDefinitions()
	for _, crd := range crds {
		oldCRD, err := crdsClient.Get(context.TODO(), crd.Name, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.Info(fmt.Sprintf("no previous CRD %v found, creating a new one.", crd.Name))
				_, err = crdsClient.Create(context.TODO(), crd, metav1.CreateOptions{})
				if err != nil {
					return fmt.Errorf("error creating CRD %v: %v", crd.Name, err)
				}
			} else {
				return fmt.Errorf("error getting CRD %v: %v", crd.Name, err)
			}
		} else {
			// Update CRDs if they already exist
			log.Info(fmt.Sprintf("previous CRD %v found, updating.", crd.Name))
			oldCRD.Spec = crd.Spec
			_, err = crdsClient.Update(context.TODO(), oldCRD, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("error updating CRD %v: %v", crd.Name, err)
			}
		}
	}

	return nil
}
