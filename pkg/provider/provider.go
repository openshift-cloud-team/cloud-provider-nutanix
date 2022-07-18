package provider

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

type NtnxCloud struct {
	name string

	client      clientset.Interface
	config      Config
	manager     *nutanixManager
	instancesV2 cloudprovider.InstancesV2
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName,
		func(config io.Reader) (cloudprovider.Interface, error) {
			return newNtnxCloud(config)
		})
}

func newNtnxCloud(config io.Reader) (cloudprovider.Interface, error) {
	bytes, err := ioutil.ReadAll(config)
	if err != nil {
		klog.Infof("Error in initializing %s cloudprovid config %q\n", ProviderName, err)
		return nil, err
	}
	klog.Infoln(string(bytes))

	nutanixConfig := Config{}
	err = json.Unmarshal(bytes, &nutanixConfig)
	if err != nil {
		return nil, err
	}
	nutanixManager, err := newNutanixManager(nutanixConfig)
	if err != nil {
		return nil, err
	}

	ntnx := &NtnxCloud{
		name:        ProviderName,
		config:      nutanixConfig,
		manager:     nutanixManager,
		instancesV2: newInstancesV2(nutanixManager),
	}

	return ntnx, err
}

// Initialize cloudprovider
func (nc *NtnxCloud) Initialize(clientBuilder cloudprovider.ControllerClientBuilder,
	stopCh <-chan struct{},
) {
	klog.Info("Initializing client ...")
	nc.addKubernetesClient(clientBuilder.ClientOrDie("cloud-provider-nutanix"))
	klog.Infof("Client initialized")
}

// SetInformers sets the informer on the cloud object. Implements cloudprovider.InformerUser
func (nc *NtnxCloud) SetInformers(informerFactory informers.SharedInformerFactory) {
	klog.Info("SetInformers")
	nc.manager.setInformers(informerFactory)
}

func (nc *NtnxCloud) addKubernetesClient(kclient clientset.Interface) {
	nc.client = kclient
	nc.manager.setKubernetesClient(kclient)
}

// ProviderName returns the cloud provider ID.
func (nc *NtnxCloud) ProviderName() string {
	return nc.name
}

// HasClusterID returns true if the cluster has a clusterID
func (nc *NtnxCloud) HasClusterID() bool {
	return true
}

func (nc *NtnxCloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nc, false
}

func (nc *NtnxCloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (nc *NtnxCloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (nc *NtnxCloud) Zones() (cloudprovider.Zones, bool) {
	klog.Info("Zones [DEPRECATED]")
	return nil, false
}

func (nc *NtnxCloud) Instances() (cloudprovider.Instances, bool) {
	klog.Info("Instances [DEPRECATED]")
	return nil, false
}

func (nc *NtnxCloud) InstancesV2() (cloudprovider.InstancesV2, bool) {
	klog.Info("InstancesV2")
	return nc.instancesV2, true
}
