package binary

import (
	"github.com/portainer/portainer/pkg/libhelm/options"
	"github.com/portainer/portainer/pkg/libhelm/release"

	"github.com/pkg/errors"
	"github.com/segmentio/encoding/json"
)

// Upgrade runs `helm install` with specified install options.
// The install options translate to CLI arguments which are passed in to the helm binary when executing install.
func (hbpm *helmBinaryPackageManager) Upgrade(upgradeOpts options.UpgradeOptions) (*release.Release, error) {
	if upgradeOpts.Name == "" {
		upgradeOpts.Name = "--generate-name"
	}
	args := []string{
		// installs a new chart if it doesn't exist, otherwise upgrades the existing release
		"--install",
		upgradeOpts.Name,
		upgradeOpts.Chart,
		"--repo", upgradeOpts.Repo,
		"--output", "json",
	}
	if upgradeOpts.Namespace != "" {
		args = append(args, "--namespace", upgradeOpts.Namespace)
	}
	if upgradeOpts.ValuesFile != "" {
		args = append(args, "--values", upgradeOpts.ValuesFile)
	}
	if upgradeOpts.Wait {
		args = append(args, "--wait")
	}
	if upgradeOpts.PostRenderer != "" {
		args = append(args, "--post-renderer", upgradeOpts.PostRenderer)
	}

	result, err := hbpm.runWithKubeConfig("upgrade", args, upgradeOpts.KubernetesClusterAccess, upgradeOpts.Env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run helm upgrade on specified args")
	}

	response := &release.Release{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal helm install response to Release struct")
	}

	return response, nil
}
