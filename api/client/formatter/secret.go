package formatter

import (
	"bytes"

	"github.com/docker/engine-api/types/secret"
)

const (
	defaultSecretQuietFormat = "{{.Name}}"
	defaultSecretTableFormat = "table {{.ID}}\t{{.Name}}"
	secretIDHeader           = "SECRET ID"
)

// SecretContext contains secret specific information required by the formatter,
// encapsulate a Context struct.
type SecretContext struct {
	Context
	// Secrets
	Secrets []*secret.Secret
}

func (ctx SecretContext) Write() {
	switch ctx.Format {
	case tableFormatKey:
		if ctx.Quiet {
			ctx.Format = defaultSecretQuietFormat
		} else {
			ctx.Format = defaultSecretTableFormat
		}
	case rawFormatKey:
		if ctx.Quiet {
			ctx.Format = `name: {{.Name}}`
		} else {
			ctx.Format = `name: {{.Name}}\nid: {{.ID}}\n`
		}
	}

	ctx.buffer = bytes.NewBufferString("")
	ctx.preformat()

	tmpl, err := ctx.parseFormat()
	if err != nil {
		return
	}

	for _, secret := range ctx.Secrets {
		secretCtx := &secretContext{
			s: secret,
		}
		err = ctx.contextFormat(tmpl, secretCtx)
		if err != nil {
			return
		}
	}

	ctx.postformat(tmpl, &secretContext{})
}

type secretContext struct {
	baseSubContext
	s *secret.Secret
}

func (c *secretContext) Name() string {
	c.addHeader(nameHeader)
	return c.s.Name
}

func (c *secretContext) ID() string {
	c.addHeader(secretIDHeader)
	return c.s.ID
}
