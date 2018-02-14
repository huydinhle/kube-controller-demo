package handler

import (
	"fmt"

	"github.com/golang/glog"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

type NamespaceHandler struct {
	client     kubernetes.Interface
	configmaps *api_v1.ConfigMapList
	namespaces *api_v1.NamespaceList
}

// NewNameSpaceHandler return a handler that help us to manage Namespace
func NewNamespaceHandler(label, sourceNamespace string, client kubernetes.Interface) (*NamespaceHandler, error) {
	// Load in all the configmaps from kube-system with the tag of the new namespace from kube-system
	configmaps, err := getConfigMaps(label, sourceNamespace, client)
	if err != nil {
		// handle error here
		return nil, err
	}

	//get namespace based on label
	nsClient := client.CoreV1().Namespaces()
	namespaces, err := nsClient.List(meta_v1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		//handle error here
		return nil, err
	}

	// Generate the Handler, using the client that got passed in
	return &NamespaceHandler{
		client:     client,
		configmaps: configmaps,
		namespaces: namespaces,
	}, nil
}

func (nh *NamespaceHandler) ProcessNamespace(namespace string) error {

	// Make sure the namespace have the kamaji-resource-controller labels, return nil and do nothing if they don't

	// go through each and every configmaps , fetch out the yaml files one by one, then install them into the new namespace
	// based on their types, right now we support configmaps and secrets

	var err error
	for _, configmap := range nh.configmaps.Items {
		err = applyConfigMapToNameSpace(namespace, nh.client, configmap)
		if err != nil {
			// log out more info
			glog.Infof("can't apply configmap %s cluster . The reason is  %v", configmap.Name, err)
		}
		glog.Infof("succesfully apply the configmap %s to the cluster", configmap.Name)
	}
	return err
}

func applyConfigMapToNameSpace(namespace string, client kubernetes.Interface, cm api_v1.ConfigMap) error {
	for _, f := range cm.Data {

		// parse the yaml into an object

		// have a switch statement.
		decode := scheme.Codecs.UniversalDeserializer().Decode
		glog.Infof("file is %s", string(f))
		obj, groupVersionKind, err := decode([]byte(f), nil, nil)

		if err != nil {
			glog.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
		}

		switch groupVersionKind.Kind {
		case "ConfigMap":
			glog.Infof("Let's create some configmap yoooo")
			cm := obj.(*api_v1.ConfigMap)

			cmClient := client.CoreV1().ConfigMaps(namespace)
			_, err = cmClient.Create(cm)
			// _, err := cmClient.Create(*o)
			if err != nil {
				glog.Infof("Creating configmap failed. Reason is %s", err)
			}
		default:
			glog.Infof("Failed to read the damn configmap")
		}
	}

	return nil
}

// get all the configmaps that has the labels list
func getConfigMaps(label string, sourceNamespace string, client kubernetes.Interface) (*api_v1.ConfigMapList, error) {
	cmClient := client.CoreV1().ConfigMaps(sourceNamespace)

	cmList, err := cmClient.List(meta_v1.ListOptions{
		LabelSelector: label,
	})
	if err != nil {
		return nil, err
	}

	return cmList, nil
}
