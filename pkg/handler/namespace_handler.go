package handler

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/huydinhle/kube-controller-demo/common"
	api_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

//Namespace handler add stuff into the namespace

type NamespaceHandler struct {
	client     kubernetes.Interface
	configmaps *api_v1.ConfigMapList
	namespaces *api_v1.NamespaceList
}

// NewNameSpaceHandler return a handler that help us to manage Namespace
func NewNameSpaceHandler(label, sourceNamespace string, client kubernetes.Interface) (NamespaceHandler, error) {
	configmaps := getConfigMaps(label, sourceNamespace, client)

	//get namespace based on label
	namespaces := nil

	// Load in all the configmaps from kube-system with the tag of the new namespace from kube-system

	// Generate the client, pass in client
	return NameSpaceHandler{
		client:     client,
		configmaps: configmaps,
		namespaces: namespaces,
	}
}

func (nh *NamespaceHandler) ProcessNamespace() {

	// Make sure the namespace have the kamaji-resource-controller labels, return nil and do nothing if they don't

	// go through each and every configmaps , fetch out the yaml files one by one, then install them into the new namespace
	// based on their types, right now we support configmaps and secrets

	for _, namespace := range nh.namespaces.Items {
		for _, configmap := range nh.configmaps.Items {
			err := applyConfigMapToNameSpace(namespace.Name(), client, configmap)
			if err != nil {
				// log out more info
				glog.Info("can't apply configmap cluster ")
			}
			glog.Info("succesfully apply the configmap to the cluster")
		}
	}
}

func applyConfigMapToNameSpace(namespace string, client kubernetes.Interface, cm api_v1.ConfigMap) err {
	for _, f := range cm.Data {

		// parse the yaml into an object

		// have a switch statement.
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, groupVersionKind, err := decode([]byte(f), nil, nil)

		if err != nil {
			log.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
		}

		// now use switch over the type of the object
		// and match each type-case
		switch o := obj.(type) {
		case *api_v1.ConfigMap:
			cmClient := client.CoreV1().ConfigMaps(namespace)
			_, err := cmClient.Create(o)
			if err != nil {
			}
		case *api_v1.Secret:
			cmClient := client.CoreV1().ConfigMaps(namespace)
		default:
			//o is unknown for us
		}
	}

	return nil
}

// get all the configmaps that has the labels list
func getConfigMaps(label string, sourceNamespace string, client kubernetes.Interface) (*api_v1.ConfigMapList, error) {

	cmClient := client.CoreV1().ConfigMaps(sourceNamespace)

	cmList, err := cmClient.List(meta_v1.ListOptions{
		LabelSelecto: label,
	})
	if err != nil {
		return nil, err
	}

	return cmList
}
