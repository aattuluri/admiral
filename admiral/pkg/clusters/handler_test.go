package clusters

import (
	argo "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/gogo/protobuf/types"
	"github.com/google/go-cmp/cmp"
	"github.com/istio-ecosystem/admiral/admiral/pkg/apis/admiral/model"
	"github.com/istio-ecosystem/admiral/admiral/pkg/controller/admiral"
	"github.com/istio-ecosystem/admiral/admiral/pkg/controller/common"
	"github.com/istio-ecosystem/admiral/admiral/pkg/controller/istio"
	"github.com/istio-ecosystem/admiral/admiral/pkg/test"
	"istio.io/api/networking/v1alpha3"
	v1alpha32 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	istiofake "istio.io/client-go/pkg/clientset/versioned/fake"
	coreV1 "k8s.io/api/core/v1"
	k8sv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"testing"
	"time"
)

func TestIgnoreIstioResource(t *testing.T) {

	//Struct of test case info. Name is required.
	testCases := []struct {
		name           string
		exportTo       []string
		expectedResult bool
	}{
		{
			name:           "Should return false when exportTo is not present",
			exportTo:       nil,
			expectedResult: false,
		},
		{
			name:           "Should return false when its exported to *",
			exportTo:       []string{"*"},
			expectedResult: false,
		},
		{
			name:           "Should return true when its exported to .",
			exportTo:       []string{"."},
			expectedResult: true,
		},
		{
			name:           "Should return true when its exported to a handful of namespaces",
			exportTo:       []string{"namespace1", "namespace2"},
			expectedResult: true,
		},
	}

	//Run the test for every provided case
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			result := IgnoreIstioResource(c.exportTo)
			if result == c.expectedResult {
				//perfect
			} else {
				t.Errorf("Failed. Got %v, expected %v", result, c.expectedResult)
			}
		})
	}
}

func TestGetDestinationRule(t *testing.T) {
	//Do setup here
	outlierDetection := &v1alpha3.OutlierDetection{
		BaseEjectionTime:  &types.Duration{Seconds: 300},
		Consecutive_5XxErrors: &types.UInt32Value{Value: 10},
		Interval:          &types.Duration{Seconds: 60}}
	mTLS := &v1alpha3.TrafficPolicy{Tls: &v1alpha3.TLSSettings{Mode: v1alpha3.TLSSettings_ISTIO_MUTUAL}, OutlierDetection: outlierDetection}

	noGtpDr := v1alpha3.DestinationRule{
		Host:          "qa.myservice.global",
		TrafficPolicy: mTLS,
	}

	basicGtpDr := v1alpha3.DestinationRule{
		Host: "qa.myservice.global",
		TrafficPolicy: &v1alpha3.TrafficPolicy{
			Tls: &v1alpha3.TLSSettings{Mode: v1alpha3.TLSSettings_ISTIO_MUTUAL},
			LoadBalancer: &v1alpha3.LoadBalancerSettings{
				LbPolicy:          &v1alpha3.LoadBalancerSettings_Simple{Simple: v1alpha3.LoadBalancerSettings_ROUND_ROBIN},
				LocalityLbSetting: &v1alpha3.LocalityLoadBalancerSetting{},
			},
			OutlierDetection: outlierDetection,
		},
	}

	failoverGtpDr := v1alpha3.DestinationRule{
		Host: "qa.myservice.global",
		TrafficPolicy: &v1alpha3.TrafficPolicy{
			Tls: &v1alpha3.TLSSettings{Mode: v1alpha3.TLSSettings_ISTIO_MUTUAL},
			LoadBalancer: &v1alpha3.LoadBalancerSettings{
				LbPolicy: &v1alpha3.LoadBalancerSettings_Simple{Simple: v1alpha3.LoadBalancerSettings_ROUND_ROBIN},
				LocalityLbSetting: &v1alpha3.LocalityLoadBalancerSetting{
					Distribute: []*v1alpha3.LocalityLoadBalancerSetting_Distribute{
						{
							From: "uswest2/*",
							To:   map[string]uint32{"us-west-2": 100},
						},
					},
				},
			},
			OutlierDetection: outlierDetection,
		},
	}

	topologyGTPPolicy := &model.TrafficPolicy{
		LbType: model.TrafficPolicy_TOPOLOGY,
		Target: []*model.TrafficGroup{
			{
				Region: "us-west-2",
				Weight: 100,
			},
		},
	}

	failoverGTPPolicy := &model.TrafficPolicy{
		LbType: model.TrafficPolicy_FAILOVER,
		Target: []*model.TrafficGroup{
			{
				Region: "us-west-2",
				Weight: 100,
			},
			{
				Region: "us-east-2",
				Weight: 0,
			},
		},
	}

	//Struct of test case info. Name is required.
	testCases := []struct {
		name            string
		host            string
		locality        string
		gtpPolicy       *model.TrafficPolicy
		destinationRule *v1alpha3.DestinationRule
	}{
		{
			name:            "Should handle a nil GTP",
			host:            "qa.myservice.global",
			locality:        "uswest2",
			gtpPolicy:       nil,
			destinationRule: &noGtpDr,
		},
		{
			name:            "Should return default DR with empty locality",
			host:            "qa.myservice.global",
			locality:        "",
			gtpPolicy:       failoverGTPPolicy,
			destinationRule: &noGtpDr,
		},
		{
			name:            "Should handle a topology GTP",
			host:            "qa.myservice.global",
			locality:        "uswest2",
			gtpPolicy:       topologyGTPPolicy,
			destinationRule: &basicGtpDr,
		},
		{
			name:            "Should handle a failover GTP",
			host:            "qa.myservice.global",
			locality:        "uswest2",
			gtpPolicy:       failoverGTPPolicy,
			destinationRule: &failoverGtpDr,
		},
	}

	//Run the test for every provided case
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			result := getDestinationRule(c.host, c.locality, c.gtpPolicy)
			if !cmp.Equal(result, c.destinationRule) {
				t.Fatalf("DestinationRule Mismatch. Diff: %v", cmp.Diff(result, c.destinationRule))
			}
		})
	}
}

func TestHandleVirtualServiceEvent(t *testing.T) {
	//Do setup here
	syncNs := v1alpha32.VirtualService{}
	syncNs.Namespace = "ns"

	tooManyHosts := v1alpha32.VirtualService{
		Spec: v1alpha3.VirtualService{
			Hosts: []string{"qa.blah.global", "e2e.blah.global"},
		},
	}
	tooManyHosts.Namespace = "other-ns"

	happyPath := v1alpha32.VirtualService{
		Spec: v1alpha3.VirtualService{
			Hosts: []string{"e2e.blah.global"},
		},
	}
	happyPath.Namespace = "other-ns"
	happyPath.Name = "vs-name"

	cnameCache := common.NewMapOfMaps()
	noDependencClustersHandler := VirtualServiceHandler{
		RemoteRegistry: &RemoteRegistry{
			RemoteControllers: map[string]*RemoteController{},
			AdmiralCache: &AdmiralCache{
				CnameDependentClusterCache: cnameCache,
				SeClusterCache:             common.NewMapOfMaps(),
			},
		},
	}

	fakeIstioClient := istiofake.NewSimpleClientset()
	goodCnameCache := common.NewMapOfMaps()
	goodCnameCache.Put("e2e.blah.global", "cluster.k8s.global", "cluster.k8s.global")
	handlerEmptyClient := VirtualServiceHandler{
		RemoteRegistry: &RemoteRegistry{
			RemoteControllers: map[string]*RemoteController{
				"cluster.k8s.global": &RemoteController{
					VirtualServiceController: &istio.VirtualServiceController{
						IstioClient: fakeIstioClient,
					},
				},
			},
			AdmiralCache: &AdmiralCache{
				CnameDependentClusterCache: goodCnameCache,
				SeClusterCache:             common.NewMapOfMaps(),
			},
		},
	}

	fullFakeIstioClient := istiofake.NewSimpleClientset()
	fullFakeIstioClient.NetworkingV1alpha3().VirtualServices("ns").Create(&v1alpha32.VirtualService{
		ObjectMeta: v12.ObjectMeta{
			Name: "vs-name",
		},
		Spec: v1alpha3.VirtualService{
			Hosts: []string{"e2e.blah.global"},
		},
	})
	handlerFullClient := VirtualServiceHandler{
		RemoteRegistry: &RemoteRegistry{
			RemoteControllers: map[string]*RemoteController{
				"cluster.k8s.global": &RemoteController{
					VirtualServiceController: &istio.VirtualServiceController{
						IstioClient: fullFakeIstioClient,
					},
				},
			},
			AdmiralCache: &AdmiralCache{
				CnameDependentClusterCache: goodCnameCache,
				SeClusterCache:             common.NewMapOfMaps(),
			},
		},
	}

	//Struct of test case info. Name is required.
	testCases := []struct {
		name          string
		vs            *v1alpha32.VirtualService
		handler       *VirtualServiceHandler
		expectedError error
		event         common.Event
	}{
		{
			name:          "Virtual Service in sync namespace",
			vs:            &syncNs,
			expectedError: nil,
			handler:       &noDependencClustersHandler,
			event:         0,
		},
		{
			name:          "Virtual Service with multiple hosts",
			vs:            &tooManyHosts,
			expectedError: nil,
			handler:       &noDependencClustersHandler,
			event:         0,
		},
		{
			name:          "No dependent clusters",
			vs:            &happyPath,
			expectedError: nil,
			handler:       &noDependencClustersHandler,
			event:         0,
		},
		{
			name:          "New Virtual Service",
			vs:            &happyPath,
			expectedError: nil,
			handler:       &handlerEmptyClient,
			event:         0,
		},
		{
			name:          "Existing Virtual Service",
			vs:            &happyPath,
			expectedError: nil,
			handler:       &handlerFullClient,
			event:         1,
		},
		{
			name:          "Deleted Virtual Service",
			vs:            &happyPath,
			expectedError: nil,
			handler:       &handlerFullClient,
			event:         2,
		},
	}

	//Run the test for every provided case
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := handleVirtualServiceEvent(c.vs, c.handler, c.event, common.VirtualService)
			if err != c.expectedError {
				t.Fatalf("Error mismatch, expected %v but got %v", c.expectedError, err)
			}
		})
	}
}

func TestGetServiceForRolloutCanary(t *testing.T) {
	//Struct of test case info. Name is required.
	const NAMESPACE = "namespace"
	const SERVICENAME = "serviceName"
	const STABLESERVICENAME = "stableserviceName"
	const CANARYSERVICENAME = "canaryserviceName"
	const VS_NAME_1 = "virtualservice1"
	const VS_NAME_2 = "virtualservice2"
	const VS_NAME_3 = "virtualservice3"
	const VS_NAME_4 = "virtualservice4"
	const VS_ROUTE_PRIMARY = "primary"
	config := rest.Config{
		Host: "localhost",
	}
	stop := make(chan struct{})

	s, e := admiral.NewServiceController(stop, &test.MockServiceHandler{}, &config, time.Second*time.Duration(300))
	r, e := admiral.NewRolloutsController(stop, &test.MockRolloutHandler{}, &config, time.Second*time.Duration(300))

	fakeIstioClient := istiofake.NewSimpleClientset()

	v := &istio.VirtualServiceController{
		IstioClient: fakeIstioClient,
	}

	if e != nil {
		t.Fatalf("Inititalization failed")
	}

	rcTemp := &RemoteController{
		VirtualServiceController: v,
		ServiceController:        s,
		RolloutController:        r}

	selectorMap := make(map[string]string)
	selectorMap["app"] = "test"

	service := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
		},
	}
	service.Name = SERVICENAME
	service.Namespace = NAMESPACE
	port1 := coreV1.ServicePort{
		Port: 8080,
	}

	port2 := coreV1.ServicePort{
		Port: 8081,
	}

	ports := []coreV1.ServicePort{port1, port2}
	service.Spec.Ports = ports

	stableService := &coreV1.Service{
		ObjectMeta: v12.ObjectMeta{Name:STABLESERVICENAME, Namespace: NAMESPACE},
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
			Ports: ports,
		},
	}

	canaryService := &coreV1.Service{
		ObjectMeta: v12.ObjectMeta{Name:CANARYSERVICENAME, Namespace: NAMESPACE},
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
			Ports: ports,
		},
	}

	selectorMap1 := make(map[string]string)
	selectorMap1["app"] = "test1"
	service1 := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
		},
	}
	service1.Name = "dummy"
	service1.Namespace = "namespace1"
	port3 := coreV1.ServicePort{
		Port: 8080,
		Name: "random3",
	}

	port4 := coreV1.ServicePort{
		Port: 8081,
		Name: "random4",
	}

	ports1 := []coreV1.ServicePort{port3, port4}
	service1.Spec.Ports = ports1

	selectorMap4 := make(map[string]string)
	selectorMap4["app"] = "test"
	service4 := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap4,
		},
	}
	service4.Name = "dummy"
	service4.Namespace = "namespace4"
	port11 := coreV1.ServicePort{
		Port: 8080,
		Name: "random3",
	}

	port12 := coreV1.ServicePort{
		Port: 8081,
		Name: "random4",
	}

	ports11 := []coreV1.ServicePort{port11, port12}
	service4.Spec.Ports = ports11

	rcTemp.ServiceController.Cache.Put(service)
	rcTemp.ServiceController.Cache.Put(service1)
	rcTemp.ServiceController.Cache.Put(service4)
	rcTemp.ServiceController.Cache.Put(stableService)
	rcTemp.ServiceController.Cache.Put(canaryService)

	virtualService := &v1alpha32.VirtualService{
		ObjectMeta: v12.ObjectMeta{Name: VS_NAME_1, Namespace: NAMESPACE},
		Spec:v1alpha3.VirtualService{
			Http:                 []*v1alpha3.HTTPRoute{{Route:[]*v1alpha3.HTTPRouteDestination{
				{Destination: &v1alpha3.Destination{Host: STABLESERVICENAME}, Weight:80},
				{Destination: &v1alpha3.Destination{Host: CANARYSERVICENAME}, Weight:20},
			}}},
		},
	}

	vsMutipleRoutesWithMatch := &v1alpha32.VirtualService{
		ObjectMeta: v12.ObjectMeta{Name: VS_NAME_2, Namespace: NAMESPACE},
		Spec:v1alpha3.VirtualService{
			Http:                 []*v1alpha3.HTTPRoute{{Name:VS_ROUTE_PRIMARY, Route:[]*v1alpha3.HTTPRouteDestination{
				{Destination: &v1alpha3.Destination{Host: STABLESERVICENAME}, Weight:80},
				{Destination: &v1alpha3.Destination{Host: CANARYSERVICENAME}, Weight:20},
			}}},
		},
	}

	vsMutipleRoutesWithZeroWeight := &v1alpha32.VirtualService{
		ObjectMeta: v12.ObjectMeta{Name: VS_NAME_4, Namespace: NAMESPACE},
		Spec:v1alpha3.VirtualService{
			Http:                 []*v1alpha3.HTTPRoute{{Name:"random", Route:[]*v1alpha3.HTTPRouteDestination{
				{Destination: &v1alpha3.Destination{Host: STABLESERVICENAME}, Weight:100},
				{Destination: &v1alpha3.Destination{Host: CANARYSERVICENAME}, Weight:0},
			}}},
		},
	}

	rcTemp.VirtualServiceController.IstioClient.NetworkingV1alpha3().VirtualServices(NAMESPACE).Create(virtualService)
	rcTemp.VirtualServiceController.IstioClient.NetworkingV1alpha3().VirtualServices(NAMESPACE).Create(vsMutipleRoutesWithMatch)
	rcTemp.VirtualServiceController.IstioClient.NetworkingV1alpha3().VirtualServices(NAMESPACE).Create(vsMutipleRoutesWithZeroWeight)

	canaryRollout := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	matchLabel := make(map[string]string)
	matchLabel["app"] = "test"

	labelSelector := v12.LabelSelector{
		MatchLabels: matchLabel,
	}
	canaryRollout.Spec.Selector = &labelSelector

	canaryRollout.Namespace = NAMESPACE
	canaryRollout.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{},
	}

	canaryRolloutNS1 := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	matchLabel2 := make(map[string]string)
	matchLabel2["app"] = "test1"

	labelSelector2 := v12.LabelSelector{
		MatchLabels: matchLabel2,
	}
	canaryRolloutNS1.Spec.Selector = &labelSelector2

	canaryRolloutNS1.Namespace = "namespace1"
	canaryRolloutNS1.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{},
	}

	canaryRolloutNS4 := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	matchLabel4 := make(map[string]string)
	matchLabel4["app"] = "test"

	labelSelector4 := v12.LabelSelector{
		MatchLabels: matchLabel4,
	}
	canaryRolloutNS4.Spec.Selector = &labelSelector4

	canaryRolloutNS4.Namespace = "namespace4"
	canaryRolloutNS4.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{},
	}

	anotationsNS4Map := make(map[string]string)
	anotationsNS4Map[common.SidecarEnabledPorts] = "8080"

	canaryRolloutNS4.Spec.Template.Annotations = anotationsNS4Map

	canaryRolloutIstioVs := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	canaryRolloutIstioVs.Spec.Selector = &labelSelector

	canaryRolloutIstioVs.Namespace = NAMESPACE
	canaryRolloutIstioVs.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{
			StableService: STABLESERVICENAME,
			CanaryService: CANARYSERVICENAME,
			TrafficRouting: &argo.RolloutTrafficRouting{
				Istio: &argo.IstioTrafficRouting{
					VirtualService: argo.IstioVirtualService{Name: VS_NAME_1},
				},
			},
		},
	}

	canaryRolloutIstioVsRouteMatch := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	canaryRolloutIstioVsRouteMatch.Spec.Selector = &labelSelector

	canaryRolloutIstioVsRouteMatch.Namespace = NAMESPACE
	canaryRolloutIstioVsRouteMatch.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{
			StableService: STABLESERVICENAME,
			CanaryService: CANARYSERVICENAME,
			TrafficRouting: &argo.RolloutTrafficRouting{
				Istio: &argo.IstioTrafficRouting{
					VirtualService: argo.IstioVirtualService{Name: VS_NAME_2, Routes: []string {VS_ROUTE_PRIMARY}},
				},
			},
		},
	}

	canaryRolloutIstioVsRouteMisMatch := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	canaryRolloutIstioVsRouteMisMatch.Spec.Selector = &labelSelector

	canaryRolloutIstioVsRouteMisMatch.Namespace = NAMESPACE
	canaryRolloutIstioVsRouteMisMatch.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{
			StableService: STABLESERVICENAME,
			CanaryService: CANARYSERVICENAME,
			TrafficRouting: &argo.RolloutTrafficRouting{
				Istio: &argo.IstioTrafficRouting{
					VirtualService: argo.IstioVirtualService{Name: VS_NAME_3, Routes: []string {"random"}},
				},
			},
		},
	}

	canaryRolloutIstioVsZeroWeight := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	canaryRolloutIstioVsZeroWeight.Spec.Selector = &labelSelector

	canaryRolloutIstioVsZeroWeight.Namespace = NAMESPACE
	canaryRolloutIstioVsZeroWeight.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{
			StableService: STABLESERVICENAME,
			CanaryService: CANARYSERVICENAME,
			TrafficRouting: &argo.RolloutTrafficRouting{
				Istio: &argo.IstioTrafficRouting{
					VirtualService: argo.IstioVirtualService{Name: VS_NAME_4},
				},
			},
		},
	}

	canaryRolloutIstioVsMimatch := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	canaryRolloutIstioVsMimatch.Spec.Selector = &labelSelector

	canaryRolloutIstioVsMimatch.Namespace = NAMESPACE
	canaryRolloutIstioVsMimatch.Spec.Strategy = argo.RolloutStrategy{
		Canary: &argo.CanaryStrategy{
			StableService: STABLESERVICENAME,
			CanaryService: CANARYSERVICENAME,
			TrafficRouting: &argo.RolloutTrafficRouting{
				Istio: &argo.IstioTrafficRouting{
					VirtualService: argo.IstioVirtualService{Name: "random"},
				},
			},
		},
	}

	resultForDummy := map[string]*WeightedService {"dummy": {Weight:1, Service:service1},}

	resultForRandomMatch := map[string]*WeightedService {CANARYSERVICENAME: {Weight:1, Service:canaryService},}

	resultForStableServiceOnly := map[string]*WeightedService {STABLESERVICENAME: {Weight:1, Service:stableService},}

	resultForCanaryWithIstio := map[string]*WeightedService {STABLESERVICENAME: {Weight:80, Service:stableService},
		CANARYSERVICENAME: {Weight:20, Service:canaryService},}

	resultForCanaryWithStableService := map[string]*WeightedService {STABLESERVICENAME: {Weight:100, Service:stableService},}

	testCases := []struct {
		name    string
		rollout *argo.Rollout
		rc      *RemoteController
		result  map[string]*WeightedService
	}{
		{
			name:    "canaryRolloutHappyCaseMeshPortAnnotationOnRollout",
			rollout: &canaryRolloutNS4,
			rc:      rcTemp,
			result:  resultForDummy,
		}, {
			name:    "canaryRolloutWithoutSelectorMatch",
			rollout: &canaryRolloutNS1,
			rc:      rcTemp,
			result:  make(map[string]*WeightedService, 0),
		}, {
			name:    "canaryRolloutHappyCase",
			rollout: &canaryRollout,
			rc:      rcTemp,
			result:  resultForRandomMatch,
		}, {
			name:    "canaryRolloutWithStableService",
			rollout: &canaryRolloutIstioVsMimatch,
			rc:      rcTemp,
			result:  resultForStableServiceOnly,
		}, {
			name:    "canaryRolloutWithIstioVirtualService",
			rollout: &canaryRolloutIstioVs,
			rc:      rcTemp,
			result:  resultForCanaryWithIstio,
		}, {
			name:    "canaryRolloutWithIstioVirtualServiceZeroWeight",
			rollout: &canaryRolloutIstioVsZeroWeight,
			rc:      rcTemp,
			result:  resultForCanaryWithStableService,
		}, {
			name:    "canaryRolloutWithIstioRouteMatch",
			rollout: &canaryRolloutIstioVsRouteMatch,
			rc:      rcTemp,
			result:  resultForCanaryWithIstio,
		}, {
			name:    "canaryRolloutWithIstioRouteMisMatch",
			rollout: &canaryRolloutIstioVsRouteMisMatch,
			rc:      rcTemp,
			result:  resultForStableServiceOnly,
		},
	}

	//Run the test for every provided case
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			result := getServiceForRollout(c.rc, c.rollout)
			if len(c.result) == 0 {
				if result != nil && len(result) != 0 {
					t.Fatalf("Service expected to be nil")
				}
			} else {
				for key, wanted := range c.result {
					if got, ok := result[key]; ok {
						if !cmp.Equal(got.Service.Name, wanted.Service.Name) {
							t.Fatalf("Service Mismatch. Diff: %v", cmp.Diff(got.Service.Name, wanted.Service.Name))
						}
						if !cmp.Equal(got.Weight, wanted.Weight) {
							t.Fatalf("Service Weight Mismatch. Diff: %v", cmp.Diff(got.Weight, wanted.Weight))
						}
					} else {
						t.Fatalf("Expected a service with name=%s but none returned", key)
					}
				}
			}
		})
	}
}

func TestGetServiceForRolloutBlueGreen(t *testing.T) {
	//Struct of test case info. Name is required.
	const NAMESPACE = "namespace"
	const SERVICENAME = "serviceNameActive"
	const ROLLOUT_POD_HASH_LABEL string = "rollouts-pod-template-hash"

	config := rest.Config{
		Host: "localhost",
	}
	stop := make(chan struct{})

	s, e := admiral.NewServiceController(stop, &test.MockServiceHandler{}, &config, time.Second*time.Duration(300))
	r, e := admiral.NewRolloutsController(stop, &test.MockRolloutHandler{}, &config, time.Second*time.Duration(300))

	emptyCacheService, e := admiral.NewServiceController(stop, &test.MockServiceHandler{}, &config, time.Second*time.Duration(300))

	if e != nil {
		t.Fatalf("Inititalization failed")
	}

	rc := &RemoteController{
		VirtualServiceController: &istio.VirtualServiceController{},
		ServiceController:        s,
		RolloutController:        r}

	bgRollout := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}

	matchLabel := make(map[string]string)
	matchLabel["app"] = "test"

	labelSelector := v12.LabelSelector{
		MatchLabels: matchLabel,
	}
	bgRollout.Spec.Selector = &labelSelector

	bgRollout.Namespace = NAMESPACE
	bgRollout.Spec.Strategy = argo.RolloutStrategy{
		BlueGreen: &argo.BlueGreenStrategy{
			ActiveService:  SERVICENAME,
			PreviewService: "previewService",
		},
	}

	selectorMap := make(map[string]string)
	selectorMap["app"] = "test"
	selectorMap[ROLLOUT_POD_HASH_LABEL] = "hash"

	activeService := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
		},
	}
	activeService.Name = SERVICENAME
	activeService.Namespace = NAMESPACE
	port1 := coreV1.ServicePort{
		Port: 8080,
		Name: "random1",
	}

	port2 := coreV1.ServicePort{
		Port: 8081,
		Name: "random2",
	}

	ports := []coreV1.ServicePort{port1, port2}
	activeService.Spec.Ports = ports

	selectorMap1 := make(map[string]string)
	selectorMap1["app"] = "test1"

	service1 := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
		},
	}
	service1.Name = "dummy"
	service1.Namespace = NAMESPACE
	port3 := coreV1.ServicePort{
		Port: 8080,
		Name: "random3",
	}

	port4 := coreV1.ServicePort{
		Port: 8081,
		Name: "random4",
	}

	ports1 := []coreV1.ServicePort{port3, port4}
	service1.Spec.Ports = ports1

	selectorMap2 := make(map[string]string)
	selectorMap2["app"] = "test"
	selectorMap2[ROLLOUT_POD_HASH_LABEL] = "hash"
	previewService := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
		},
	}
	previewService.Name = "previewService"
	previewService.Namespace = NAMESPACE
	port5 := coreV1.ServicePort{
		Port: 8080,
		Name: "random3",
	}

	port6 := coreV1.ServicePort{
		Port: 8081,
		Name: "random4",
	}

	ports2 := []coreV1.ServicePort{port5, port6}

	previewService.Spec.Ports = ports2

	serviceNS1 := &coreV1.Service{
		Spec: coreV1.ServiceSpec{
			Selector: selectorMap,
		},
	}
	serviceNS1.Name = "dummy"
	serviceNS1.Namespace = "namespace1"
	port8 := coreV1.ServicePort{
		Port: 8080,
		Name: "random3",
	}

	port9 := coreV1.ServicePort{
		Port: 8081,
		Name: "random4",
	}

	ports12 := []coreV1.ServicePort{port8, port9}
	serviceNS1.Spec.Ports = ports12

	rc.ServiceController.Cache.Put(service1)
	rc.ServiceController.Cache.Put(previewService)
	rc.ServiceController.Cache.Put(activeService)
	rc.ServiceController.Cache.Put(serviceNS1)

	noStratergyRollout := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}
	noStratergyRollout.Namespace = NAMESPACE

	noStratergyRollout.Spec.Strategy = argo.RolloutStrategy{}

	bgRolloutNs1 := argo.Rollout{
		Spec: argo.RolloutSpec{Template: coreV1.PodTemplateSpec{
			ObjectMeta: k8sv1.ObjectMeta{Annotations: map[string]string{}},
		}}}

	matchLabel1 := make(map[string]string)
	matchLabel1["app"] = "test"

	labelSelector1 := v12.LabelSelector{
		MatchLabels: matchLabel,
	}
	bgRolloutNs1.Spec.Selector = &labelSelector1

	bgRolloutNs1.Namespace = "namespace1"
	bgRolloutNs1.Spec.Strategy = argo.RolloutStrategy{
		BlueGreen: &argo.BlueGreenStrategy{
			ActiveService:  SERVICENAME,
			PreviewService: "previewService",
		},
	}

	resultForBlueGreen := map[string]*WeightedService {SERVICENAME: {Weight:1, Service:activeService},}

	testCases := []struct {
		name    string
		rollout *argo.Rollout
		rc      *RemoteController
		result  map[string]*WeightedService
	}{
		{
			name:    "canaryRolloutNoLabelMatch",
			rollout: &bgRolloutNs1,
			rc:      rc,
			result:  make(map[string]*WeightedService, 0),
		}, {
			name:    "canaryRolloutNoStratergy",
			rollout: &noStratergyRollout,
			rc:      rc,
			result:  make(map[string]*WeightedService, 0),
		}, {
			name:    "canaryRolloutHappyCase",
			rollout: &bgRollout,
			rc:      rc,
			result:  resultForBlueGreen ,
		},
		{
			name:    "canaryRolloutNilRollout",
			rollout: nil,
			rc:      rc,
			result:  make(map[string]*WeightedService, 0),
		},
		{
			name:    "canaryRolloutEmptyServiceCache",
			rollout: &bgRollout,
			rc: &RemoteController{
				ServiceController: emptyCacheService,
			},
			result: make(map[string]*WeightedService, 0),
		},
	}

	//Run the test for every provided case
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			result := getServiceForRollout(c.rc, c.rollout)
			if len(c.result) == 0 {
				if result != nil && len(result) > 0 {
					t.Fatalf("Service expected to be nil")
				}
			} else {
				for key, service := range c.result {
					if val, ok := result[key]; ok {
						if !cmp.Equal(val.Service.Name, service.Service.Name) {
							t.Fatalf("Service Mismatch. Diff: %v", cmp.Diff(val.Service.Name, service.Service.Name))
						}
					} else {
						t.Fatalf("Expected a service with name=%s but none returned", key)
					}
				}
			}
		})
	}
}