package app

import (
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"testing"
	"time"

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

func TestApp_processNewLogLines(t *testing.T) {
	type fields struct {
		logFile              string
		notifier             notifierClient
		hostMachine          string
		watchSettings        config.WatchSettings
		processedLineTracker processedLineTracker
	}
	type args struct {
		file              reader
		lastProcessedLine int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "process log file with no lines",
			fields: fields{},
			args: args{
				file: &appfakes.FakeReader{
					ReadStub: func([]byte) (int, error) {
						return 0, io.EOF
					},
				},
				lastProcessedLine: 0,
			},
			wantErr: false,
		},
		{
			name: "happy path",
			fields: fields{
				notifier: &appfakes.FakeNotifierClient{
					NotifyStub: func(notifier.LogLine) error {
						return nil
					},
				},
			},
			args: args{
				file: &appfakes.FakeReader{
					ReadStub: func([]byte) (int, error) {
						return 1, nil
					},
				},
				lastProcessedLine: 5,
			},
			wantErr: false,
		},
		// {
		// 	name: "process multiple lines, skipping initial ones",
		// 	fields: fields{
		// 		notifier: &appfakes.FakeNotifierClient{
		// 			NotifyStub: func(notifier.LogLine) error {
		// 				return nil
		// 			},
		// 		},
		// 		processedLineTracker: &appfakes.FakeProcessedLineTracker{
		// 			UpdateLastProcessedLineStub: func(int) error {
		// 				return nil
		// 			},
		// 		},
		// 	},
		// 	args: args{
		// 		file: &appfakes.FakeReader{
		// 			ReadStub: func(p []byte) (int, error) {
		// 				switch callCount {
		// 				case 0:
		// 					copy(p, "line1\n")
		// 					return len("line1\n"), nil
		// 				case 1:
		// 					copy(p, "line2\n")
		// 					return len("line2\n"), nil
		// 				case 2:
		// 					return 0, io.EOF
		// 				}
		// 				return 0, io.EOF
		// 			},
		// 		},
		// 		lastProcessedLine: 0,
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "error from processLine",
			fields: fields{
				notifier: &appfakes.FakeNotifierClient{
					NotifyStub: func(notifier.LogLine) error {
						return fmt.Errorf("notify error")
					},
				},
			},
			args: args{
				file: &appfakes.FakeReader{
					ReadStub: func([]byte) (int, error) {
						return 1, nil
					},
				},
				lastProcessedLine: 1,
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
			}
			if err := a.processNewLogLines(tt.args.file, tt.args.lastProcessedLine); (err != nil) != tt.wantErr {
				t.Errorf("App.processNewLogLines() error = %v, wantErr %v", err, tt.wantErr)
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

func TestNew(t *testing.T) {
	fakeNotifier := &appfakes.FakeNotifierClient{}
	fakeFile := &appfakes.FakeFile{}

	type args struct {
		logFile              string
		notifier             notifierClient
		hostMachine          string
		watchSettings        config.WatchSettings
		processedLineTracker processedLineTracker
		file                 file
	}
	tests := []struct {
		name string
		args args
		want App
	}{
		{
			name: "create app happy path",
			args: args{
				logFile:       "foo",
				notifier:      fakeNotifier,
				hostMachine:   "bar",
				watchSettings: config.WatchSettings{},
				file:          fakeFile,
			},
			want: App{
				logFile:       "foo",
				notifier:      fakeNotifier,
				hostMachine:   "bar",
				watchSettings: config.WatchSettings{},
				file:          fakeFile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.logFile, tt.args.notifier, tt.args.hostMachine, tt.args.watchSettings, tt.args.processedLineTracker, tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLogRotated(t *testing.T) {
	timeA := time.Now()

	fileInfoA := &appfakes.FakeFileInfo{
		IsDirStub: func() bool {
			return false
		},
		ModTimeStub: func() time.Time {
			return timeA
		},
		NameStub: func() string {
			return "foo"
		},
		ModeStub: func() fs.FileMode {
			return 0644
		},
		SizeStub: func() int64 {
			return 100
		},
		SysStub: func() any {
			return nil
		},
	}
	type args struct {
		currentFileInfo fs.FileInfo
		lastFileInfo    fs.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "log not rotated",
			args: args{
				currentFileInfo: fileInfoA,
				lastFileInfo:    fileInfoA,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLogRotated(tt.args.currentFileInfo, tt.args.lastFileInfo); got != tt.want {
				t.Errorf("isLogRotated() = %v, want %v", got, tt.want)
			}
		})
	}
}
