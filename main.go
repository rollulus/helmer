package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func (h *helmer) install(releaseName string, values map[string]interface{}) error {
	client := action.NewInstall(h.actionConfig)
	client.ReleaseName = releaseName

	_, err := client.Run(h.chart, values)
	return errors.Wrap(err, "client.Run")
}

func (h *helmer) delete(releaseName string) error {
	client := action.NewUninstall(h.actionConfig)
	_, err := client.Run(releaseName)
	return errors.Wrap(err, "client.Run")
}

func (h *helmer) exists(releaseName string) (bool, error) {
	client := action.NewGet(h.actionConfig)
	_, err := client.Run(releaseName)
	if err == nil {
		return true, nil
	}
	if err == driver.ErrReleaseNotFound {
		return false, nil
	}
	return false, errors.Wrap(err, "client.Run")
}

type helmer struct {
	namespace    string
	chart        *chart.Chart
	actionConfig *action.Configuration
}

func newHelmer(chartName, namespace string) (*helmer, error) {
	// setup helm
	actionConfig := new(action.Configuration)
	helmDriver := os.Getenv("HELM_DRIVER")

	k := kube.GetConfig("", "", namespace)

	if err := actionConfig.Init(k, namespace, helmDriver, func(format string, v ...interface{}) {
		log.Printf(format, v...)
	}); err != nil {
		return nil, errors.Wrap(err, "actionConfig.Init")
	}

	// load chart
	chart, err := loader.Load(chartName)
	if err != nil {
		return nil, errors.Wrap(err, "loader.Load")
	}

	return &helmer{namespace: namespace, chart: chart, actionConfig: actionConfig}, nil
}

func main() {
	h, err := newHelmer("./charts/mychart", "default")
	if err != nil {
		panic(err)
	}

	exists, err := h.exists("roel")
	if err != nil {
		panic(err)
	}
	if exists {
		if err := h.delete("roel"); err != nil {
			panic(err)
		}
	}
	if err := h.install("roel", map[string]interface{}{}); err != nil {
		panic(err)
	}
}
