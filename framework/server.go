package framework

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	normanclientbase "github.com/rancher/norman/clientbase"
	normantypes "github.com/rancher/norman/types"
	rclusterv3 "github.com/rancher/types/client/cluster/v3"
	rmgmtv3 "github.com/rancher/types/client/management/v3"
	rprojectv3 "github.com/rancher/types/client/project/v3"
)

var (
	insecureClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)

// RancherServer is used to hold the information to a
// single Rancher Server installation. This is pointing to
// the IP address of the install, which means it's agnostic to
// HA installation.
type RancherServer struct {
	URL                  string
	AccessKey            string
	SecretKey            string
	TokenKey             string
	ClusterName          string
	APIEndPoint          string
	ManagementClient     *rmgmtv3.Client
	Cluster              map[string]*rclusterv3.Client
	DefaultClusterClient *rclusterv3.Client
	DefaultCluster       *rmgmtv3.Cluster
}

// NewRancherServerFromEnvVars creates a RancherServer struct
// by reading the information from environment variables
func NewRancherServerFromEnvVars() (*RancherServer, error) {
	var err error
	var url, accessKey, secretKey, tokenKey, clusterName string
	var apiEndpoint string
	var mgmtClient *rmgmtv3.Client
	var clusterClient *rclusterv3.Client
	var rs *RancherServer

	if url = os.Getenv("RANCHER_SERVER_URL"); url == "" {
		return nil, fmt.Errorf("RANCHER_SERVER_URL not specified")
	}
	apiEndpoint = url + "/v3"

	accessKey = os.Getenv("RANCHER_ACCESS_KEY")
	secretKey = os.Getenv("RANCHER_SECRET_KEY")
	tokenKey = os.Getenv("RANCHER_TOKEN")

	if accessKey == "" && secretKey == "" && tokenKey == "" {
		return rs, fmt.Errorf("either access/secret key or token needs to be specified")
	}
	clusterName = os.Getenv("RANCHER_DEFAULT_CLUSTER_NAME")

	mgmtClientOpts := normanclientbase.ClientOpts{
		URL:        apiEndpoint,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		TokenKey:   tokenKey,
		HTTPClient: insecureClient,
	}
	mgmtClient, err = rmgmtv3.NewClient(&mgmtClientOpts)
	if err != nil {
		return rs, fmt.Errorf("error creating managment client: %v", err)
	}

	clusterListOpts := &normantypes.ListOpts{}
	if clusterName != "" {
		clusterListOpts.Filters = map[string]interface{}{"name": clusterName}
	}

	clusterCollection, err := mgmtClient.Cluster.List(clusterListOpts)
	if err != nil {
		return rs, fmt.Errorf("error fetching cluster list: %v", err)
	}
	if clusterName == "" && len(clusterCollection.Data) != 1 {
		return rs, fmt.Errorf("found %v clusters, expected either to find one cluster or default cluster name to be specified", len(clusterCollection.Data))
	}

	defaultCluster := clusterCollection.Data[0]

	clusterClientOpts := normanclientbase.ClientOpts{
		URL:        defaultCluster.Links["self"],
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		TokenKey:   tokenKey,
		HTTPClient: insecureClient,
	}
	clusterClient, err = rclusterv3.NewClient(&clusterClientOpts)
	if err != nil {
		return rs, fmt.Errorf("error fetching cluster client: %v", err)
	}

	rs = &RancherServer{
		URL:                  url,
		AccessKey:            accessKey,
		SecretKey:            secretKey,
		ClusterName:          clusterName,
		TokenKey:             tokenKey,
		ManagementClient:     mgmtClient,
		DefaultClusterClient: clusterClient,
		DefaultCluster:       &defaultCluster,

		APIEndPoint: apiEndpoint,
	}

	return rs, nil
}

func (rs *RancherServer) GetProjectAPIEndpointByID(projectID string) string {
	return rs.APIEndPoint + "/projects/" + projectID
}

func (rs *RancherServer) GetProjectClientByID(projectID string) (*rprojectv3.Client, error) {
	projectClientOpts := normanclientbase.ClientOpts{
		URL:        rs.GetProjectAPIEndpointByID(projectID),
		AccessKey:  rs.AccessKey,
		SecretKey:  rs.SecretKey,
		TokenKey:   rs.TokenKey,
		HTTPClient: insecureClient,
	}
	return rprojectv3.NewClient(&projectClientOpts)
}
