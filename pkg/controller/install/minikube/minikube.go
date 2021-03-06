package minikube

import (
	"context"

	mf "github.com/jcrossley3/manifestival"
	"github.com/openshift-knative/knative-serving-operator/pkg/controller/install/common"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("minikube")

// Configure minikube if we're soaking in it
func Configure(c client.Client, _ *runtime.Scheme) []mf.Transformer {
	node := &v1.Node{}
	if err := c.Get(context.TODO(), types.NamespacedName{Name: "minikube"}, node); err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "Unable to query for minikube node")
		}
		return nil // not running on minikube
	}
	return []mf.Transformer{egress}
}

func egress(u *unstructured.Unstructured) *unstructured.Unstructured {
	if u.GetKind() == "ConfigMap" && u.GetName() == "config-network" {
		data := map[string]string{
			"istio.sidecar.includeOutboundIPRanges": "10.0.0.1/24",
		}
		common.UpdateConfigMap(u, data, log)
	}
	return u
}
