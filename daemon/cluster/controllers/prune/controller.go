package prune

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/filters"
	executorpkg "github.com/docker/docker/daemon/cluster/executor"
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

type pruneController struct {
	backend executorpkg.Backend
	task    *api.Task
}

func NewPruneController(b executorpkg.Backend, task *api.Task) (*pruneController, error) {
	return &pruneController{
		backend: b,
		task:    task,
	}, nil
}

func (p *pruneController) Update(ctx context.Context, t *api.Task) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Update")

	return nil
}

func (p *pruneController) Prepare(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Prepare")
	return nil
}

func (p *pruneController) Start(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Start")

	f := filters.NewArgs()
	r := p.task.Spec.GetRuntime()
	resp, err := p.backend.ContainersPrune(f)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"controller": "prune",
		"prune":      fmt.Sprintf("%+v", resp),
	}).Debug("prune output")
	return nil
}

func (p *pruneController) Wait(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Waiting")

	return nil
}

func (p *pruneController) Shutdown(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Shutdown")
	return nil
}

func (p *pruneController) Terminate(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Terminate")
	return nil
}

func (p *pruneController) Remove(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Remove")
	return nil
}

func (p *pruneController) Close() error {
	logrus.WithFields(logrus.Fields{
		"controller": "prune",
	}).Debug("Close")

	return nil
}
