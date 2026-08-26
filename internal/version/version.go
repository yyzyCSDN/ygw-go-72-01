package version

import "powergw/internal/model"

type Controller struct {
	manager *Manager
}

func NewController(manager *Manager) *Controller {
	return &Controller{manager: manager}
}

func (c *Controller) Activate(ver model.ProtocolVersion, tables TableApplier) error {
	ver.State = model.VersionDraft
	return c.manager.Apply(ver, tables)
}

func (c *Controller) SwitchTo(ver model.ProtocolVersion, tables TableApplier, resetter DedupResetter) error {
	return c.manager.Switch(ver, tables, resetter)
}

func (c *Controller) ProtocolFor(channelID string) (model.Protocol, error) {
	proto, _, err := c.manager.ActiveFor(channelID)
	return proto, err
}

func (c *Controller) Current() (model.ProtocolVersion, bool) {
	return c.manager.Active()
}

func (c *Controller) Snapshot() *model.VersionSnapshot {
	return c.manager.Snapshot()
}

func (c *Controller) RegisterChannel(channelID string, proto model.Protocol) {
	c.manager.RegisterChannel(channelID, proto)
}

func (c *Controller) Supersede() {
	c.manager.Supersede()
}

func (c *Controller) Reset() {
	c.manager.ResetState()
}
