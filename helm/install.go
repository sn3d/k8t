package helm

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/sn3d/k8t"
	"helm.sh/helm/v3/pkg/action"
)

type InstallOpts struct {

	// here you can pass your context. It it's not set, the default
	// context.Background() will be used
	Context context.Context

	// if it's not set, the helm's name is used
	ReleaseName string

	// namespace where helm will be installed. If it's not set, the
	// cluster's default test namespace will be used
	Namespace string

	// it's oposite of replace. Tests are different than production helm
	// installments. For testing, we want replace content by default.
	DontReplace bool

	// the helm chart is laoded as directory from this given FS implementation.
	// If it's empty, the current working dir as FS is used
	filesystem fs.FS
}

// Install helm chart to cluster from given directory, with provided
// values. This is simple less-verbose version of InstallWithOpts
func Install(cluster *k8t.Cluster, dir string, values Value) error {
	return InstallWithOpts(cluster, dir, values, InstallOpts{})
}

// install helm chart from given directory and with given values to cluster
func InstallWithOpts(cluster *k8t.Cluster, dir string, values Value, opts InstallOpts) error {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	namespace := opts.Namespace
	if namespace == "" {
		namespace = cluster.TestNamespace()
	}

	filesystem := opts.filesystem
	if filesystem == nil {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		filesystem = os.DirFS(wd)
	}

	// load chart
	chart, err := loadChartFromFS(filesystem, dir)
	if err != nil {
		return err
	}

	releaseName := opts.ReleaseName
	if releaseName == "" {
		releaseName = chart.Metadata.Name
	}

	// prepare helm installation
	cfg := &action.Configuration{}
	err = cfg.Init(newCustomRESTClientGetter(cluster), namespace, "", func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	})
	if err != nil {
		return err
	}

	client := action.NewInstall(cfg)
	client.Namespace = namespace
	client.Replace = !opts.DontReplace
	client.ReleaseName = releaseName

	// run installation
	rel, err := client.RunWithContext(ctx, chart, values)
	if err != nil {
		return err
	}

	fmt.Printf("installed %s into %s\n", rel.Name, rel.Namespace)
	return nil
}
