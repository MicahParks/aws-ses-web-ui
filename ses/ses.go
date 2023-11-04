package ses

import (
	"context"
	"fmt"
	netMail "net/mail"

	jt "github.com/MicahParks/jsontype"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/MicahParks/aws-ses-web-ui/model"
)

const (
	charSet = "UTF-8"
)

// Config is the configuration for the AWS SES provider. It includes assets to create an AWS session.
type Config struct {
	AWSRegion   string `json:"awsRegion"`
	AccessKeyID string `json:"accessKeyID"`
	SecretKey   string `json:"secretKey"`
}

// DefaultsAndValidate implements the jsontype.Config interface.
func (c Config) DefaultsAndValidate() (Config, error) {
	if c.AWSRegion == "" {
		return Config{}, fmt.Errorf("AWS region not provided in configuration: %w", jt.ErrDefaultsAndValidate)
	}
	if c.AccessKeyID == "" {
		return Config{}, fmt.Errorf("AWS access key ID not provided in configuration: %w", jt.ErrDefaultsAndValidate)
	}
	if c.SecretKey == "" {
		return Config{}, fmt.Errorf("AWS secret key not provided in configuration: %w", jt.ErrDefaultsAndValidate)
	}
	return c, nil
}

type SES struct {
	conf Config
	ses  *ses.SES
}

func New(c Config) (SES, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKeyID, c.SecretKey, ""),
		Region:      aws.String(c.AWSRegion),
	})
	if err != nil {
		return SES{}, fmt.Errorf("failed to create AWS session: %w", err)
	}

	s := SES{
		conf: c,
		ses:  ses.New(sess),
	}

	return s, nil
}

func (s SES) Send(ctx context.Context, c model.Compose) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  recipientsPtr(c.RecipientTo),
			CcAddresses:  recipientsPtr(c.RecipientCC),
			BccAddresses: recipientsPtr(c.RecipientBCC),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: nil,
				Text: &ses.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(c.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(c.Subject),
			},
		},
		Source: aws.String(c.From.Get().String()),
	}
	_, err := s.ses.SendEmailWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func recipientsPtr(s []*jt.JSONType[*netMail.Address]) []*string {
	ptrs := make([]*string, 0, len(s))
	for _, v := range s {
		addr := v.Get().String()
		ptrs = append(ptrs, &addr)
	}
	return ptrs
}
