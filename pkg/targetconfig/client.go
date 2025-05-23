package targetconfig

import (
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"

	"github.com/kyverno/policy-reporter/pkg/crd/api/targetconfig/v1alpha1"
	tcv1alpha1 "github.com/kyverno/policy-reporter/pkg/crd/client/targetconfig/clientset/versioned"
	tcinformer "github.com/kyverno/policy-reporter/pkg/crd/client/targetconfig/informers/externalversions"
	"github.com/kyverno/policy-reporter/pkg/target"
)

type Client struct {
	targetFactory target.Factory
	collection    *target.Collection
	informer      cache.SharedIndexInformer
}

func (c *Client) ConfigureInformer() {
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			tc := obj.(*v1alpha1.TargetConfig)
			zap.L().Info("new target", zap.String("name", tc.Name))

			t, err := c.targetFactory.CreateSingleClient(tc)
			if err != nil {
				zap.L().Error("unable to create target from TargetConfig", zap.String("name", tc.Name), zap.Error(err))
				return
			} else if t == nil {
				zap.L().Error("provided TargetConfig is invalid", zap.String("name", tc.Name))
				return
			}

			c.collection.AddTarget(tc.Name, t)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			tc := newObj.(*v1alpha1.TargetConfig)
			zap.L().Info("update target", zap.String("name", tc.Name))

			t, err := c.targetFactory.CreateSingleClient(tc)
			if err != nil {
				zap.L().Error("unable to create target from TargetConfig", zap.String("name", tc.Name), zap.Error(err))
				return
			} else if t == nil {
				zap.L().Error("provided TargetConfig is invalid", zap.String("name", tc.Name))
				return
			}

			c.collection.AddTarget(tc.Name, t)
		},
		DeleteFunc: func(obj interface{}) {
			tc := obj.(*v1alpha1.TargetConfig)
			zap.L().Info("delete target", zap.String("name", tc.Name))

			c.collection.RemoveTarget(tc.Name)
		},
	})
}

func (c *Client) Run(stopChan chan struct{}) {
	go c.informer.Run(stopChan)

	if !cache.WaitForCacheSync(stopChan, c.informer.HasSynced) {
		zap.L().Error("Failed to sync target config cache")
		return
	}

	zap.L().Info("target config cache synced")
}

func NewClient(tcClient tcv1alpha1.Interface, f target.Factory, targets *target.Collection) *Client {
	tcInformer := tcinformer.NewSharedInformerFactory(tcClient, 0)

	return &Client{
		informer:      tcInformer.Policyreporter().V1alpha1().TargetConfigs().Informer(),
		targetFactory: f,
		collection:    targets,
	}
}
