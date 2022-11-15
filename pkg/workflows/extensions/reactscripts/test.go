package reactscripts

import (
	_ "embed"

	"dagger.io/dagger"
	"github.com/trustacks/trustacks/pkg/workflows/utils"
)

func Test(project *utils.Project, client *dagger.Client) (*dagger.Container, error) {
	action := &utils.Action{
		Image: "node",
		Run: []string{
			"npm install",
			"CI=true npx react-scripts test",
		},
	}
	return action.Exec(client, project, &utils.RunOpts{
		Caches: []string{"/src/node_modules"},
	})
}
