package handlers

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/project"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/rancher-compose/lookup"
	"github.com/rancher/rancher-compose/rancher"
)

func constructProjectUpgrade(logger *logrus.Entry, env *client.Environment, upgradeOpts client.EnvironmentUpgrade, url, accessKey, secretKey string) (*project.Project, map[string]interface{}, error) {
	variables := map[string]interface{}{}
	for k, v := range env.Environment {
		variables[k] = v
	}

	for k, v := range upgradeOpts.Environment {
		variables[k] = v
	}

	context := rancher.Context{
		Context: project.Context{
			ProjectName: env.Name,
			ComposeBytes: [][]byte{
				[]byte(upgradeOpts.DockerCompose),
			},
			ResourceLookup: &lookup.FileConfigLookup{},
			EnvironmentLookup: &lookup.MapEnvLookup{
				Env: variables,
			},
		},
		Url:                 fmt.Sprintf("%s/projects/%s/schemas", url, env.AccountId),
		AccessKey:           accessKey,
		SecretKey:           secretKey,
		RancherComposeBytes: []byte(upgradeOpts.RancherCompose),
		Upgrade:             true,
	}

	p, err := rancher.NewProject(&context)
	if err != nil {
		return nil, nil, err
	}

	p.AddListener(NewListenLogger(logger, p))
	return p, variables, p.Parse()
}

func constructProject(logger *logrus.Entry, env *client.Environment, url, accessKey, secretKey string) (*rancher.Context, *project.Project, error) {
	context := rancher.Context{
		Context: project.Context{
			ProjectName: env.Name,
			ComposeBytes: [][]byte{
				[]byte(env.DockerCompose),
			},
			ResourceLookup: &lookup.FileConfigLookup{},
			EnvironmentLookup: &lookup.MapEnvLookup{
				Env: env.Environment,
			},
		},
		Url:                 fmt.Sprintf("%s/projects/%s/schemas", url, env.AccountId),
		AccessKey:           accessKey,
		SecretKey:           secretKey,
		RancherComposeBytes: []byte(env.RancherCompose),
	}

	p, err := rancher.NewProject(&context)
	if err != nil {
		return nil, nil, err
	}

	p.AddListener(NewListenLogger(logger, p))
	return &context, p, nil
}
