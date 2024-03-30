<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# app

```go
import "github.com/mgla96/ssh-watcher/internal/app"
```

## Index

- [type App](<#App>)
  - [func NewApp\(logFile string, notifier notifierClient, hostMachine string, watchSettings config.WatchSettings, processedLineTracker processedLineTracker\) App](<#NewApp>)
  - [func \(a App\) Watch\(\) error](<#App.Watch>)


<a name="App"></a>
## type [App](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/app/app.go#L42-L48>)



```go
type App struct {
    // contains filtered or unexported fields
}
```

<a name="NewApp"></a>
### func [NewApp](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/app/app.go#L32>)

```go
func NewApp(logFile string, notifier notifierClient, hostMachine string, watchSettings config.WatchSettings, processedLineTracker processedLineTracker) App
```



<a name="App.Watch"></a>
### func \(App\) [Watch](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/app/app.go#L124>)

```go
func (a App) Watch() error
```

TODO\(mgottlieb\) refactor this into more unit\-testable funcs.

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)