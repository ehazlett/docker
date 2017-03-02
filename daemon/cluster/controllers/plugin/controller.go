package plugin

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/swarmkit/api"
	"golang.org/x/net/context"
)

type pluginController struct{}

func NewPluginController() (*pluginController, error) {
	return &pluginController{}, nil
}

func (p *pluginController) Update(ctx context.Context, t *api.Task) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Update")
	return nil
}

func (p *pluginController) Prepare(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Prepare")
	return nil
}

func (p *pluginController) Start(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Start")
	return nil
}

func (p *pluginController) Wait(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Wait")
	return nil
}

func (p *pluginController) Shutdown(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Shutdown")
	return nil
}

func (p *pluginController) Terminate(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Terminate")
	return nil
}

func (p *pluginController) Remove(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Remove")
	return nil
}

func (p *pluginController) Close() error {
	logrus.WithFields(logrus.Fields{
		"controller": "plugin",
	}).Debug("Close")
	return nil
}
