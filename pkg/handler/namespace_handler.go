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
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

//Namespace handler add stuff into the namespace

type NamespaceHandler struct {
	client     kubernetes.Interface
	configmaps []api_v1.ConfigMap
}

// NewNameSpaceHandler return a handler that help us to manage Namespace
func NewNameSpaceHandler() {

	// Load in all the configmaps from kube-system with the tag of the new namespace from kube-system

	// Generate the client, pass in client
}

func (nh *NamespaceHandler) ProcessNamespace() error {

	// Make sure the namespace have the kamaji-resource-controller labels, return nil and do nothing if they don't

	// go through each and every configmaps , fetch out the yaml files one by one, then install them into the new namespace
	// based on their types, right now we support configmaps and secrets
}

func getConfigMaps(labels []string, namespace string, client kubernetes.Interface) {
}
