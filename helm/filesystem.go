package helm

import (
	"errors"
	"io/fs"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// load helm chart from given directory and also filesystem. It returns
// you buffered helm chart, where all files are pre-loaded from given FS.
//
// Helm library didn't support such a functionality how to load chart from
// custom FS.
func loadChartFromFS(fileSystem fs.FS, dir string) (*chart.Chart, error) {
	files := make([]*loader.BufferedFile, 0)

	if fileSystem == nil {
		return nil, errors.New("no filesystem")
	}

	err := fs.WalkDir(fileSystem, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := fs.Stat(fileSystem, path)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			data, err := fs.ReadFile(fileSystem, path)
			if err != nil {
				return err
			}

			file := &loader.BufferedFile{
				Name: strings.TrimPrefix(path, dir+"/"),
				Data: data,
			}
			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	chart, err := loader.LoadFiles(files)
	return chart, err
}
