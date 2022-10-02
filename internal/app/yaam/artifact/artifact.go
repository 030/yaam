package artifact

import (
	"fmt"

	"github.com/spf13/viper"
)

type artefact struct {
	path, url string
}

func allowedRepos(name string) ([]string, error) {
	groups := viper.GetStringMapStringSlice("groups.maven")
	var repos []string
	if values, ok := groups[name]; ok {
		repos = values
	} else {
		return nil, fmt.Errorf("group: '%s' not found in config file", name)
	}

	return repos, nil
}
