// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/drone/drone/cmd/drone-server/config"
	"github.com/drone/drone/operator/manager"

	"github.com/google/wire"

	"github.com/drone-runners/drone-runner-docker/engine"
	"github.com/drone-runners/drone-runner-docker/engine/compiler"
	"github.com/drone-runners/drone-runner-docker/engine/linter"
	"github.com/drone-runners/drone-runner-docker/engine/resource"

	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/environ/provider"
	"github.com/drone/runner-go/pipeline/reporter/history"
	"github.com/drone/runner-go/pipeline/reporter/remote"
	"github.com/drone/runner-go/pipeline/runtime"
	"github.com/drone/runner-go/poller"
	"github.com/drone/runner-go/registry"
	"github.com/drone/runner-go/secret"
)

// wire set for loading the server.
var runnerSet = wire.NewSet(
	provideRunner,
)

// provideRunner is a Wire provider function that returns a
// local build runner configured from the environment.
func provideRunner(
	manager manager.BuildManager,
	config config.Config,
) (*poller.Poller, error) {
	// the local runner is only created when the nomad scheduler,
	// kubernetes scheduler, and remote agents are disabled
	if config.Nomad.Enabled || config.Kube.Enabled || (config.Agent.Disabled == false) {
		return nil, nil
	}

	engine, err := engine.NewEnv(engine.Opts{})
	if err != nil {
		return nil, err
	}

	var cli client.Client
	remote := remote.New(cli)
	tracer := history.New(remote)

	runner := &runtime.Runner{
		Client:   cli,
		Machine:  "localhost",
		Reporter: tracer,
		Lookup:   resource.Lookup,
		Lint:     linter.New().Lint,
		Compiler: &compiler.Compiler{
			Clone:      config.Runner.Clone,
			Privileged: append(config.Runner.Privileged, compiler.Privileged...),
			Networks:   config.Runner.Networks,
			Volumes:    config.Runner.Volumes,
			Resources: compiler.Resources{
				Memory:     int64(config.Runner.Limits.MemLimit),
				MemorySwap: int64(config.Runner.Limits.MemSwapLimit),
				CPUQuota:   config.Runner.Limits.CPUQuota,
				CPUPeriod:  config.Runner.Limits.CPUPeriod,
				CPUShares:  config.Runner.Limits.CPUShares,
				CPUSet:     config.Runner.Limits.CPUSet,
				ShmSize:    int64(config.Runner.Limits.ShmSize),
			},
			Environ:  provider.Combine(),
			Registry: registry.Combine(),
			Secret:   secret.Combine(),
		},
		Exec: runtime.NewExecer(
			tracer,
			remote,
			engine,
			0,
		).Exec,
	}

	return &poller.Poller{
		Client:   cli,
		Dispatch: runner.Run,
		Filter: &client.Filter{
			Kind:    resource.Kind,
			Type:    resource.Type,
			OS:      config.Runner.OS,
			Arch:    config.Runner.Arch,
			Variant: config.Runner.Variant,
			Kernel:  config.Runner.Kernel,
			Labels:  config.Runner.Labels,
		},
	}, nil
}
