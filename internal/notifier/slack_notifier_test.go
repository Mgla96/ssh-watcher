package notifier

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/mgla96/ssh-watcher/internal/notifier/notifierfakes"
	"github.com/rs/zerolog"
)

func TestNewSlackNotifier(t *testing.T) {
	type args struct {
		webhookURL    string
		slackChannel  string
		slackUsername string
		slackIcon     string
		log           zerolog.Logger
	}
	tests := []struct {
		name string
		args args
		want SlackNotifier
	}{
		{
			name: "create slack notifier",
			args: args{
				webhookURL:    "http://localhost",
				slackChannel:  "test",
				slackUsername: "foobar",
				slackIcon:     ":ghost:",
				log:           zerolog.Nop(),
			},
			want: SlackNotifier{
				WebhookURL:    "http://localhost",
				SlackChannel:  "test",
				SlackUsername: "foobar",
				SlackIcon:     ":ghost:",
				HttpClient:    &http.Client{},
				log:           zerolog.Nop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlackNotifier(tt.args.webhookURL, tt.args.slackChannel, tt.args.slackUsername, tt.args.slackIcon, tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSlackNotifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackNotifier_Notify(t *testing.T) {
	type fields struct {
		WebhookURL    string
		SlackChannel  string
		SlackUsername string
		SlackIcon     string
		HttpClient    hTTPClient
		log           zerolog.Logger
	}
	type args struct {
		logLine LogLine
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error creating slack request",
			fields: fields{
				WebhookURL:    "http://localhost",
				SlackChannel:  "test",
				SlackUsername: "foobar",
				SlackIcon:     ":ghost:",
				HttpClient: &notifierfakes.FakeHTTPClient{
					DoStub: func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("error")
					},
				},
				log: zerolog.Nop(),
			},
			args: args{
				logLine: LogLine{
					Username:    "test",
					IpAddress:   "1.2.3.4",
					LoginTime:   "Dec 1",
					EventType:   LoggedIn,
					HostMachine: "foobar",
				},
			},
			wantErr: true,
		},
		{
			name: "unsuccessful status code",
			fields: fields{
				WebhookURL:    "http://localhost",
				SlackChannel:  "test",
				SlackUsername: "foobar",
				SlackIcon:     ":ghost:",
				HttpClient: &notifierfakes.FakeHTTPClient{
					DoStub: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusBadRequest,
							Body: &notifierfakes.FakeReadCloser{
								CloseStub: func() error {
									return nil
								},
								ReadStub: func(p []byte) (n int, err error) {
									return 0, nil
								},
							},
						}, nil
					},
				},
				log: zerolog.Nop(),
			},
			args: args{
				logLine: LogLine{
					Username:    "test",
					IpAddress:   "1.2.3.4",
					LoginTime:   "Dec 1",
					EventType:   LoggedIn,
					HostMachine: "foobar",
				},
			},
			wantErr: true,
		},
		{
			name: "successful status code",
			fields: fields{
				WebhookURL:    "http://localhost",
				SlackChannel:  "test",
				SlackUsername: "foobar",
				SlackIcon:     ":ghost:",
				HttpClient: &notifierfakes.FakeHTTPClient{
					DoStub: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: &notifierfakes.FakeReadCloser{
								CloseStub: func() error {
									return nil
								},
								ReadStub: func(p []byte) (n int, err error) {
									return 0, nil
								},
							},
						}, nil
					},
				},
				log: zerolog.Nop(),
			},
			args: args{
				logLine: LogLine{
					Username:    "test",
					IpAddress:   "1.2.3.4",
					LoginTime:   "Dec 1",
					EventType:   LoggedIn,
					HostMachine: "foobar",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SlackNotifier{
				WebhookURL:    tt.fields.WebhookURL,
				SlackChannel:  tt.fields.SlackChannel,
				SlackUsername: tt.fields.SlackUsername,
				SlackIcon:     tt.fields.SlackIcon,
				HttpClient:    tt.fields.HttpClient,
				log:           tt.fields.log,
			}
			if err := s.Notify(tt.args.logLine); (err != nil) != tt.wantErr {
				t.Errorf("SlackNotifier.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
