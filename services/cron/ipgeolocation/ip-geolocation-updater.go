// Package ipgeolocation implements functions to update public IP and geolocation details on the edge node
package ipgeolocation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/decentralized-cloud/edge-core/services/configuration"
	cronContract "github.com/decentralized-cloud/edge-core/services/cron"
	commonErrors "github.com/micro-business/go-core/system/errors"
	cron "github.com/robfig/cron/v3"
	"github.com/shengdoushi/base58"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type cronService struct {
	logger            *zap.Logger
	cronSpec          string
	ipinfoUrl         string
	ipinfoAccessToken string
	cron              *cron.Cron
	clientset         *kubernetes.Clientset
	runningNodeName   string
	clusterType       configuration.ClusterType
}

type ipinfoResponse struct {
	Ip       string
	Hostname string
	City     string
	Region   string
	Country  string
	Loc      string
	Org      string
	Postal   string
	Timezone string
}

var Live bool
var Ready bool

func init() {
	Live = false
	Ready = false
}

var acceptedCharactersForLabels = base58.NewAlphabet("ABCDEFGHJKLMNPQRSTUVWXYZ123456789abcdefghijkmnopqrstuvwxyz")

// NewCronService creates new instance of the cronService, setting up all dependencies and returns the instance
// logger: Mandatory. Reference to the logger service
// configurationService: Mandatory. Reference to the service that provides required configurations
// Returns the new service or error if something goes wrong
func NewCronService(
	logger *zap.Logger,
	configurationService configuration.ConfigurationContract) (cronContract.CronContract, error) {
	if logger == nil {
		return nil, commonErrors.NewArgumentNilError("logger", "logger is required")
	}

	if configurationService == nil {
		return nil, commonErrors.NewArgumentNilError("configurationService", "configurationService is required")
	}

	clusterType, err := configurationService.GetEdgeClusterType()
	if err != nil {
		return nil, err
	}

	if clusterType != configuration.K3S {
		return nil, commonErrors.NewUnknownError(fmt.Sprintf("clusterType %v is not supported", clusterType))
	}

	cronSpec, err := configurationService.GetGeolocationUpdaterCronSpec()
	if err != nil {
		return nil, err
	}

	ipinfoUrl, err := configurationService.GetIpinfoUrl()
	if err != nil {
		return nil, err
	}

	ipinfoAccessToken, err := configurationService.GetIpinfoAccessToken()
	if err != nil {
		return nil, err
	}

	runningNodeName, err := configurationService.GetRunningNodeName()
	if err != nil {
		return nil, err
	}

	k8sRestConfig, err := getRestConfig(logger)
	if err != nil {
		return nil, err
	}

	var clientset *kubernetes.Clientset
	if clientset, err = kubernetes.NewForConfig(k8sRestConfig); err != nil {
		return nil, commonErrors.NewUnknownErrorWithError("Failed to create client set", err)
	}

	return &cronService{
		logger:            logger,
		cronSpec:          cronSpec,
		ipinfoUrl:         ipinfoUrl,
		ipinfoAccessToken: ipinfoAccessToken,
		cron:              cron.New(),
		clientset:         clientset,
		runningNodeName:   runningNodeName,
		clusterType:       clusterType,
	}, nil
}

// Start starts the Geolocation Updater service
// Returns error if something goes wrong
func (service *cronService) Start() error {
	service.logger.Info("Geolocation Updater service started")

	_, err := service.cron.AddFunc(service.cronSpec, service.updateGeolocation)
	if err != nil {
		return err
	}

	service.cron.Start()

	go service.updateGeolocation()

	Live = true
	Ready = true

	return nil
}

// Stop stops the Geolocation Updater service
// Returns error if something goes wrong
func (service *cronService) Stop() error {
	Live = false
	Ready = false

	service.cron.Stop()

	return nil
}

// Stop stops the Geolocation Updater service
// Returns error if something goes wrong
func (service *cronService) updateGeolocation() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)

	defer cancelFunc()

	shouldUpdate, err := service.shouldUpdateGeolocation(ctx)
	if err != nil {
		return
	}

	if !shouldUpdate {
		service.logger.Debug("Manual update is set. Skipping geolocation update.")

		return
	}

	service.logger.Info("Updating geolocation...")

	ipinfoResponse, err := service.getGeolocationDetails(ctx)
	if err != nil {
		return
	}

	err = service.updateNode(ctx, ipinfoResponse)
	if err != nil {
		return
	}

	service.logger.Info("Finished updating geolocation details.")
}

func getRestConfig(logger *zap.Logger) (*rest.Config, error) {
	if kubeConfig := os.Getenv("KUBECONFIG"); kubeConfig != "" {
		logger.Info("path ", zap.String("KUBECONFIG", kubeConfig))

		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	logger.Info("homePath ", zap.String("path", homePath))

	kubeConfigFilePath := filepath.Join(homePath, ".kube", "config")

	_, err = os.Stat(kubeConfigFilePath)
	if !os.IsNotExist(err) {
		return clientcmd.BuildConfigFromFlags("", kubeConfigFilePath)
	}

	return rest.InClusterConfig()
}

func (service *cronService) shouldUpdateGeolocation(ctx context.Context) (bool, error) {
	node, err := service.clientset.CoreV1().Nodes().Get(ctx, service.runningNodeName, metav1.GetOptions{})
	if err != nil {
		service.logger.Error(
			"Failed to retrieve node information",
			zap.String("runningNodeName", service.runningNodeName),
			zap.Error(err))

		return false, err
	}

	if value, ok := node.Labels["edgecloud9.geolocation.manual"]; ok {
		if value == "false" {
			return false, nil
		}

		return true, nil
	}

	return true, nil
}

func (service *cronService) getGeolocationDetails(ctx context.Context) (*ipinfoResponse, error) {
	httpClient := &http.Client{}
	request, err := http.NewRequestWithContext(ctx, "GET", service.ipinfoUrl, nil)
	if err != nil {
		service.logger.Error(
			"Failed to create a new request to Ipinfo",
			zap.String("ipinfoUrl", service.ipinfoUrl),
			zap.Error(err))

		return nil, err
	}

	if strings.Trim(service.ipinfoAccessToken, " ") != "" {
		request.Header.Set("Authorization", "Bearer "+service.ipinfoAccessToken)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		service.logger.Error("Failed to send request to Ipinfo", zap.String("ipinfoUrl", service.ipinfoUrl), zap.Error(err))

		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		service.logger.Error(
			"Failed to read Ipinfo reponse body",
			zap.String("ipinfoUrl", service.ipinfoUrl),
			zap.String("response", string(body)),
			zap.Error(err))

		return nil, err
	}

	var ipinfoResponse ipinfoResponse

	err = json.Unmarshal(body, &ipinfoResponse)
	if err != nil {
		service.logger.Error(
			"Can't deserialize Ipinfo response",
			zap.String("ipinfoUrl", service.ipinfoUrl),
			zap.String("response", string(body)),
			zap.Error(err))

		return nil, err
	}

	return &ipinfoResponse, nil
}

func (service *cronService) updateNode(ctx context.Context, ipinfoResponse *ipinfoResponse) error {
	currentTime := base58.Encode([]byte(time.Now().Format(time.RFC3339Nano)), acceptedCharactersForLabels)

	patch := struct {
		Metadata struct {
			Labels map[string]string `json:"labels"`
		} `json:"metadata"`
	}{}

	patch.Metadata.Labels = map[string]string{}
	patch.Metadata.Labels["k3s.io/external-ip"] = ipinfoResponse.Ip
	patch.Metadata.Labels["edgecloud9.public.lastUpdatedTime"] = currentTime
	patch.Metadata.Labels["edgecloud9.public.ip"] = ipinfoResponse.Ip
	patch.Metadata.Labels["edgecloud9.public.hostname"] = base58.Encode([]byte(ipinfoResponse.Hostname), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.lastUpdatedTime"] = currentTime
	patch.Metadata.Labels["edgecloud9.geolocation.loc"] = base58.Encode([]byte(ipinfoResponse.Loc), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.city"] = base58.Encode([]byte(ipinfoResponse.City), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.region"] = base58.Encode([]byte(ipinfoResponse.Region), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.country"] = base58.Encode([]byte(ipinfoResponse.Country), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.org"] = base58.Encode([]byte(ipinfoResponse.Org), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.postal"] = base58.Encode([]byte(ipinfoResponse.Postal), acceptedCharactersForLabels)
	patch.Metadata.Labels["edgecloud9.geolocation.timezone"] = base58.Encode([]byte(ipinfoResponse.Timezone), acceptedCharactersForLabels)

	patchJson, err := json.Marshal(patch)
	if err != nil {
		service.logger.Error(
			"Failed to serialize geolocations details",
			zap.String("runningNodeName", service.runningNodeName),
			zap.Error(err))

		return err
	}

	if _, err = service.clientset.CoreV1().Nodes().Patch(ctx, service.runningNodeName, types.MergePatchType, patchJson, metav1.PatchOptions{}); err != nil {
		service.logger.Error(
			"Failed to retrieve node information",
			zap.String("runningNodeName", service.runningNodeName),
			zap.Error(err))

		return err
	}

	return nil
}
