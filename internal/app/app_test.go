package app

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mgla96/ssh-watcher/internal/app/appfakes"
	"github.com/mgla96/ssh-watcher/internal/config"
	"github.com/mgla96/ssh-watcher/internal/notifier"
)

func TestLogWatcher_shouldSendMessage(t *testing.T) {
	type fields struct {
		watchSettings config.WatchSettings
	}
	type args struct {
		eventType notifier.EventType
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test should send accepted login",
			fields: fields{
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               false,
					FailedLoginInvalidUsername: false,
				},
			},
			args: args{
				eventType: notifier.LoggedIn,
			},
			want: true,
		},
		{
			name: "test should send failed login",
			fields: fields{
				watchSettings: config.WatchSettings{
					AcceptedLogins:             false,
					FailedLogins:               true,
					FailedLoginInvalidUsername: false,
				},
			},
			args: args{
				eventType: notifier.FailedLoginAttempt,
			},
			want: true,
		},
		{
			name: "test should send failed login invalid username",
			fields: fields{
				watchSettings: config.WatchSettings{
					AcceptedLogins:             false,
					FailedLogins:               false,
					FailedLoginInvalidUsername: true,
				},
			},
			args: args{
				eventType: notifier.FailedLoginAttemptInvalidUsername,
			},
			want: true,
		},
		{
			name: "test should not send accepted login",
			fields: fields{
				watchSettings: config.WatchSettings{
					AcceptedLogins:             false,
					FailedLogins:               true,
					FailedLoginInvalidUsername: true,
				},
			},
			args: args{
				eventType: notifier.LoggedIn,
			},
			want: false,
		},
		{
			name: "test should not send failed login",
			fields: fields{
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               false,
					FailedLoginInvalidUsername: true,
				},
			},
			args: args{
				eventType: notifier.FailedLoginAttempt,
			},
			want: false,
		},
		{
			name: "test should not send failed login invalid username",
			fields: fields{
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               true,
					FailedLoginInvalidUsername: false,
				},
			},
			args: args{
				eventType: notifier.FailedLoginAttemptInvalidUsername,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := App{
				watchSettings: tt.fields.watchSettings,
			}
			if got := w.shouldSendMessage(tt.args.eventType); got != tt.want {
				t.Errorf("App.shouldSendMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLogLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want notifier.LogLine
	}{
		{
			name: "accepted password",
			args: args{
				line: "Dec 1 10:0:0 fake sshd[0000]: Accepted password for foo from 1.2.3.4 port 57000 ssh2",
			},
			want: notifier.LogLine{
				Username:  "foo",
				IpAddress: "1.2.3.4",
				LoginTime: "Dec 1",
				EventType: notifier.LoggedIn,
			},
		},
		{
			name: "failed password",
			args: args{
				line: "Dec 1 10:0:0 fake sshd[0000]: Failed password for foo from 1.2.3.4 port 57000 ssh2",
			},
			want: notifier.LogLine{
				Username:  "foo",
				IpAddress: "1.2.3.4",
				LoginTime: "Dec 1",
				EventType: notifier.FailedLoginAttempt,
			},
		},
		{
			name: "invalid user",
			args: args{
				line: "Dec 1 10:0:0 fake sshd[0000]: Failed password for invalid user bar from 1.2.3.4 port 57000 ssh2",
			},
			want: notifier.LogLine{
				Username:  "bar",
				IpAddress: "1.2.3.4",
				LoginTime: "Dec 1",
				EventType: notifier.FailedLoginAttemptInvalidUsername,
			},
		},
		{
			name: "not sshd line",
			args: args{
				line: "Dec 1 10:0:0 fake foobar",
			},
			want: notifier.LogLine{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := App{}
			if got := w.parseLogLine(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLogLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_processLine(t *testing.T) {
	type fields struct {
		logFile              string
		notifier             notifierClient
		hostMachine          string
		watchSettings        config.WatchSettings
		processedLineTracker processedLineTracker
		file                 file
	}
	type args struct {
		line       string
		lineNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				notifier: &appfakes.FakeNotifierClient{
					NotifyStub: func(notifier.LogLine) error {
						return nil
					},
				},
				processedLineTracker: &appfakes.FakeProcessedLineTracker{
					UpdateLastProcessedLineStub: func(int) error {
						return nil
					},
				},
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               true,
					FailedLoginInvalidUsername: true,
				},
			},
			args: args{
				line:       "Mar 30 00:00:00 foo sshd[5052]: Invalid user foo from x.x.x.x port xxx",
				lineNumber: 1,
			},
			wantErr: false,
		},
		{
			name: "shouldn't send failed login message",
			fields: fields{
				notifier: &appfakes.FakeNotifierClient{
					NotifyStub: func(notifier.LogLine) error {
						return nil
					},
				},
				processedLineTracker: &appfakes.FakeProcessedLineTracker{
					UpdateLastProcessedLineStub: func(int) error {
						return nil
					},
				},
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               true,
					FailedLoginInvalidUsername: false,
				},
			},
			args: args{
				line:       "Mar 30 00:00:00 foo sshd[5052]: Invalid user foo from x.x.x.x port xxx",
				lineNumber: 1,
			},
			wantErr: false,
		},
		{
			name: "error notifying",
			fields: fields{
				notifier: &appfakes.FakeNotifierClient{
					NotifyStub: func(notifier.LogLine) error {
						return fmt.Errorf("error notifying")
					},
				},
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               true,
					FailedLoginInvalidUsername: true,
				},
			},
			args: args{
				line:       "Mar 30 00:00:00 foo sshd[5052]: Invalid user foo from x.x.x.x port xxx",
				lineNumber: 1,
			},
			wantErr: true,
		},
		{
			name: "error updating last processed line",
			fields: fields{
				notifier: &appfakes.FakeNotifierClient{
					NotifyStub: func(notifier.LogLine) error {
						return nil
					},
				},
				processedLineTracker: &appfakes.FakeProcessedLineTracker{
					UpdateLastProcessedLineStub: func(int) error {
						return fmt.Errorf("error updating last processed line")
					},
				},
				watchSettings: config.WatchSettings{
					AcceptedLogins:             true,
					FailedLogins:               true,
					FailedLoginInvalidUsername: true,
				},
			},
			args: args{
				line:       "Mar 30 00:00:00 foo sshd[5052]: Invalid user foo from x.x.x.x port xxx",
				lineNumber: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := App{
				logFile:              tt.fields.logFile,
				notifier:             tt.fields.notifier,
				hostMachine:          tt.fields.hostMachine,
				watchSettings:        tt.fields.watchSettings,
				processedLineTracker: tt.fields.processedLineTracker,
				file:                 tt.fields.file,
			}
			if err := a.processLine(tt.args.line, tt.args.lineNumber); (err != nil) != tt.wantErr {
				t.Errorf("App.processLine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
