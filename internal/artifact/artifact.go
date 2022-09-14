package artifact

import (
	"fmt"
	"path/filepath"

	"github.com/030/yaam/internal/pkg/project"
	"github.com/spf13/viper"
)

type artefact struct {
	path, url string
}

func allowedRepos(name string) ([]string, error) {
	h, err := project.Home()
	if err != nil {
		return []string{}, err
	}

	viper.SetConfigName("groups")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(h, "conf"))
	if err := viper.ReadInConfig(); err != nil {
		return []string{}, err
	}

	groups := viper.GetStringMapStringSlice("groups")
	var repos []string
	if values, ok := groups[name]; ok {
		repos = values
	} else {
		return []string{}, fmt.Errorf("group: '%s' not found in config file", name)
	}

	return repos, nil
}
