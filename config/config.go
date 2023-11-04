package config

import (
	"fmt"
	"net/mail"
	"strings"

	hhpostgres "github.com/MicahParks/httphandle/postgres"
	jt "github.com/MicahParks/jsontype"

	"github.com/MicahParks/aws-ses-web-ui/ses"
)

type Config struct {
	ASWU     ASWU              `json:"aswu"`
	Postgres hhpostgres.Config `json:"postgres"`
	SES      ses.Config        `json:"ses"`
}

func (c Config) DefaultsAndValidate() (Config, error) {
	var err error
	c.ASWU, err = c.ASWU.DefaultsAndValidate()
	if err != nil {
		return Config{}, fmt.Errorf("failed to validate ASWU configuration: %w", err)
	}
	if c.ASWU.UsePostgres {
		c.Postgres, err = c.Postgres.DefaultsAndValidate()
		if err != nil {
			return Config{}, fmt.Errorf("failed to validate PostgreSQL configuration: %w", err)
		}
	}
	c.SES, err = c.SES.DefaultsAndValidate()
	if err != nil {
		return Config{}, fmt.Errorf("failed to validate SES configuration: %w", err)
	}
	return c, nil
}

func (c Config) DevMode() bool {
	return c.ASWU.DevMode
}

type ASWU struct {
	AllowedFrom []string                    `json:"allowedFrom"`
	DefaultFrom *jt.JSONType[*mail.Address] `json:"defaultFrom"`
	DevMode     bool                        `json:"devMode"`
	UsePostgres bool                        `json:"usePostgres"`
}

func (a ASWU) DefaultsAndValidate() (ASWU, error) {
	for _, v := range a.AllowedFrom {
		if strings.HasPrefix(v, "@") {
			if len(v) < 2 {
				return ASWU{}, fmt.Errorf("invalid allowedFrom email address domain %q", v)
			}
		} else {
			_, err := mail.ParseAddress(v)
			if err != nil {
				return ASWU{}, fmt.Errorf("invalid allowedFrom email address %q: %w", v, err)
			}
		}
	}
	return a, nil
}
