package mock

import (
	"fmt"

	prismClientV3 "github.com/nutanix-cloud-native/prism-go-client/v3"
)

type MockPrism struct {
	mockEnvironment MockEnvironment
}

func (mp *MockPrism) GetVM(vmUUID string) (*prismClientV3.VMIntentResponse, error) {
	if v, ok := mp.mockEnvironment.managedMockMachines[vmUUID]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf(entityNotFoundError)
	}
}

func (mp *MockPrism) GetCluster(clusterUUID string) (*prismClientV3.ClusterIntentResponse, error) {
	return mp.mockEnvironment.managedMockClusters[clusterUUID], nil
}

func (mp *MockPrism) ListAllCluster(fitler string) (*prismClientV3.ClusterListIntentResponse, error) {
	entities := make([]*prismClientV3.ClusterIntentResponse, 0)

	for _, e := range mp.mockEnvironment.managedMockClusters {
		entities = append(entities, e)
	}
	return &prismClientV3.ClusterListIntentResponse{
		Entities: entities,
	}, nil
}
