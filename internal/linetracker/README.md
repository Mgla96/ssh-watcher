<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# linetracker

```go
import "github.com/mgla96/ssh-watcher/internal/linetracker"
```

## Index

- [type FileProcessedLineTracker](<#FileProcessedLineTracker>)
  - [func NewFileProcessedLineTracker\(stateFilePath string\) FileProcessedLineTracker](<#NewFileProcessedLineTracker>)
  - [func \(f FileProcessedLineTracker\) GetLastProcessedLine\(\) \(int, error\)](<#FileProcessedLineTracker.GetLastProcessedLine>)
  - [func \(f FileProcessedLineTracker\) UpdateLastProcessedLine\(lineNumber int\) error](<#FileProcessedLineTracker.UpdateLastProcessedLine>)


<a name="FileProcessedLineTracker"></a>
## type [FileProcessedLineTracker](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/linetracker/line_tracker.go#L24-L26>)



```go
type FileProcessedLineTracker struct {
    StateFilePath string
}
```

<a name="NewFileProcessedLineTracker"></a>
### func [NewFileProcessedLineTracker](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/linetracker/line_tracker.go#L13>)

```go
func NewFileProcessedLineTracker(stateFilePath string) FileProcessedLineTracker
```



<a name="FileProcessedLineTracker.GetLastProcessedLine"></a>
### func \(FileProcessedLineTracker\) [GetLastProcessedLine](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/linetracker/line_tracker.go#L30>)

```go
func (f FileProcessedLineTracker) GetLastProcessedLine() (int, error)
```

GetLastProcessedLine reads the statefile and extracts the last processed line number in the ssh log file.

<a name="FileProcessedLineTracker.UpdateLastProcessedLine"></a>
### func \(FileProcessedLineTracker\) [UpdateLastProcessedLine](<https://github.com/Mgla96/ssh-watcher/blob/main/internal/linetracker/line_tracker.go#L58>)

```go
func (f FileProcessedLineTracker) UpdateLastProcessedLine(lineNumber int) error
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
