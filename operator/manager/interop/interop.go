package interop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone/operator/manager"
	"github.com/drone/runner-go/client"
)

var noContext = context.Background()

// An Client allows an embedded runner to interop directly with
// the server process.
type Client struct {
	m manager.BuildManager
}

// Join notifies the server the runner is joining the cluster.
func (p *Client) Join(_ context.Context, machine string) error {
	return nil // no-op
}

// Leave notifies the server the runner is leaving the cluster.
func (p *Client) Leave(_ context.Context, machine string) error {
	return nil // no-op
}

// Ping sends a ping message to the server to test connectivity.
func (p *Client) Ping(_ context.Context, machine string) error {
	return nil // no-op
}

// Request requests the next available build stage for execution.
func (p *Client) Request(_ context.Context, args *client.Filter) (*drone.Stage, error) {
	req := &manager.Request{
		Kind:    args.Kind,
		Type:    args.Type,
		OS:      args.OS,
		Arch:    args.Arch,
		Variant: args.Variant,
		Kernel:  args.Kernel,
		Labels:  args.Labels,
	}
	stage, err := p.m.Request(noContext, req)
	return toStage(stage), err
}

// Accept accepts the build stage for execution.
func (p *Client) Accept(_ context.Context, src *drone.Stage) error {
	dst, err := p.m.Accept(noContext, src.ID, src.Machine)
	if dst != nil {
		src.Updated = dst.Updated
		src.Version = dst.Version
	}
	return err
}

// TODO
// Detail gets the build stage details for execution.
func (p *Client) Detail(_ context.Context, stage *drone.Stage) (*client.Context, error) {
	res, err := p.m.Details(noContext, stage.ID)
	if err != nil {
		return nil, err
	}
	netrc, err := p.m.Netrc(noContext, res.Repo.ID)
	if err != nil {
		return nil, err
	}
	return &client.Context{
		Build:   toBuild(res.Build),
		Stage:   stage,
		Config:  nil,
		Netrc:   toNetrc(netrc),
		Repo:    toRepo(res.Repo),
		Secrets: nil,
		System:  nil,
	}, nil
}

// Update updates the build stage.
func (p *Client) Update(_ context.Context, src *drone.Stage) (err error) {
	for i, step := range src.Steps {
		// a properly implemented runner should never encounter
		// input errors. these checks are included to help
		// developers creating new runners.
		if step.Number == 0 {
			return fmt.Errorf("step[%d] missing number", i)
		}
		if step.StageID == 0 {
			return fmt.Errorf("step[%d] missing stage id", i)
		}
		if step.Status == drone.StatusRunning &&
			step.Started == 0 {
			return fmt.Errorf("step[%d] missing start time", i)
		}
	}
	dst := fromStage(src)
	if dst.Status == drone.StatusPending ||
		dst.Status == drone.StatusRunning {
		err = p.m.BeforeAll(noContext, dst)
	} else {
		err = p.m.AfterAll(noContext, dst)
	}

	src.Updated = dst.Updated
	src.Version = dst.Version

	set := map[int]*drone.Step{}
	for _, step := range dst.Steps {
		set[step.Number] = toStep(step)
	}
	for _, step := range src.Steps {
		from, ok := set[step.Number]
		if ok {
			step.ID = from.ID
			step.StageID = from.StageID
			step.Started = from.Started
			step.Stopped = from.Stopped
			step.Version = from.Version
		}
	}
	return nil
}

// UpdateStep updates the build step.
func (p *Client) UpdateStep(_ context.Context, src *drone.Step) (err error) {
	dst := fromStep(src)
	if dst.Status == drone.StatusPending ||
		dst.Status == drone.StatusRunning {
		err = p.m.Before(noContext, dst)
	} else {
		err = p.m.After(noContext, dst)
	}
	src.Version = dst.Version
	return err
}

// Watch watches for build cancellation requests.
func (p *Client) Watch(ctx context.Context, build int64) (bool, error) {
	return p.m.Watch(ctx, build) // TODO do we send context here?
}

// Batch batch writes logs to the build logs.
func (p *Client) Batch(_ context.Context, step int64, lines []*drone.Line) error {
	for _, line := range lines {
		err := p.m.Write(noContext, step, fromLine(line))
		if err != nil {
			return err
		}
	}
	return nil
}

// Upload uploads the full logs to the server.
func (p *Client) Upload(_ context.Context, step int64, lines []*drone.Line) error {
	raw, err := json.Marshal(lines)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(raw)
	return p.m.Upload(noContext, step, buf)
}
