package controller

import (
	"github.com/golang/glog"
	"github.com/kolibox/koli/pkg/spec"
	"k8s.io/kubernetes/pkg/api"
	apps "k8s.io/kubernetes/pkg/apis/apps"
	extensions "k8s.io/kubernetes/pkg/apis/extensions"
)

type queue struct {
	addonch chan *spec.Addon
	spch    chan *spec.ServicePlan
	dpch    chan *extensions.Deployment
	psch    chan *apps.StatefulSet
	nsch    chan *api.Namespace
}

func newQueue(size int) *queue {
	return &queue{
		addonch: make(chan *spec.Addon, size),
		spch:    make(chan *spec.ServicePlan, size),
		dpch:    make(chan *extensions.Deployment),
		psch:    make(chan *apps.StatefulSet),
		nsch:    make(chan *api.Namespace, size),
	}
}

func (q *queue) add(o interface{}) {
	switch obj := o.(type) {
	case *spec.Addon:
		q.addonch <- o.(*spec.Addon)
	case *spec.ServicePlan:
		q.spch <- o.(*spec.ServicePlan)
	case *extensions.Deployment:
		q.dpch <- o.(*extensions.Deployment)
	case *apps.StatefulSet:
		q.psch <- o.(*apps.StatefulSet)
	case *api.Namespace:
		q.nsch <- o.(*api.Namespace)
	default:
		glog.Infof("add: unknown type (%T)", obj)
	}
}
func (q *queue) close() {
	close(q.addonch)
	close(q.spch)
	close(q.dpch)
	close(q.psch)
	close(q.nsch)
}

func (q *queue) pop(o interface{}) (interface{}, bool) {
	switch t := o.(type) {
	case *spec.Addon:
		obj, ok := <-q.addonch
		return obj, ok
	case *spec.ServicePlan:
		obj, ok := <-q.spch
		return obj, ok
	case *extensions.Deployment:
		obj, ok := <-q.dpch
		return obj, ok
	case *apps.StatefulSet:
		obj, ok := <-q.psch
		return obj, ok
	case *api.Namespace:
		obj, ok := <-q.nsch
		return obj, ok
	default:
		glog.Warningf("pop: unknown type (%T)", t)
		return nil, false
	}
}