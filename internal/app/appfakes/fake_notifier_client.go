// Code generated by counterfeiter. DO NOT EDIT.
package appfakes

import (
	"sync"

	"github.com/mgla96/ssh-watcher/internal/notifier"
)

type FakeNotifierClient struct {
	NotifyStub        func(notifier.LogLine) error
	notifyMutex       sync.RWMutex
	notifyArgsForCall []struct {
		arg1 notifier.LogLine
	}
	notifyReturns struct {
		result1 error
	}
	notifyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNotifierClient) Notify(arg1 notifier.LogLine) error {
	fake.notifyMutex.Lock()
	ret, specificReturn := fake.notifyReturnsOnCall[len(fake.notifyArgsForCall)]
	fake.notifyArgsForCall = append(fake.notifyArgsForCall, struct {
		arg1 notifier.LogLine
	}{arg1})
	stub := fake.NotifyStub
	fakeReturns := fake.notifyReturns
	fake.recordInvocation("Notify", []interface{}{arg1})
	fake.notifyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeNotifierClient) NotifyCallCount() int {
	fake.notifyMutex.RLock()
	defer fake.notifyMutex.RUnlock()
	return len(fake.notifyArgsForCall)
}

func (fake *FakeNotifierClient) NotifyCalls(stub func(notifier.LogLine) error) {
	fake.notifyMutex.Lock()
	defer fake.notifyMutex.Unlock()
	fake.NotifyStub = stub
}

func (fake *FakeNotifierClient) NotifyArgsForCall(i int) notifier.LogLine {
	fake.notifyMutex.RLock()
	defer fake.notifyMutex.RUnlock()
	argsForCall := fake.notifyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeNotifierClient) NotifyReturns(result1 error) {
	fake.notifyMutex.Lock()
	defer fake.notifyMutex.Unlock()
	fake.NotifyStub = nil
	fake.notifyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNotifierClient) NotifyReturnsOnCall(i int, result1 error) {
	fake.notifyMutex.Lock()
	defer fake.notifyMutex.Unlock()
	fake.NotifyStub = nil
	if fake.notifyReturnsOnCall == nil {
		fake.notifyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.notifyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeNotifierClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.notifyMutex.RLock()
	defer fake.notifyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeNotifierClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
