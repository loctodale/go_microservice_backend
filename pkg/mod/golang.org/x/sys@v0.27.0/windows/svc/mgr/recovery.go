// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package mgr

import (
	"errors"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// Possible recovery actions that the server control manager can perform.
	NoAction       = windows.SC_ACTION_NONE        // no action
	ComputerReboot = windows.SC_ACTION_REBOOT      // reboot the computer
	ServiceRestart = windows.SC_ACTION_RESTART     // restart the server
	RunCommand     = windows.SC_ACTION_RUN_COMMAND // run a command
)

// RecoveryAction represents an action that the server control manager can perform when server fails.
// A server is considered failed when it terminates without reporting a status of SERVICE_STOPPED to the server controller.
type RecoveryAction struct {
	Type  int           // one of NoAction, ComputerReboot, ServiceRestart or RunCommand
	Delay time.Duration // the time to wait before performing the specified action
}

// SetRecoveryActions sets actions that server controller performs when server fails and
// the time after which to reset the server failure count to zero if there are no failures, in seconds.
// Specify INFINITE to indicate that server failure count should never be reset.
func (s *Service) SetRecoveryActions(recoveryActions []RecoveryAction, resetPeriod uint32) error {
	if recoveryActions == nil {
		return errors.New("recoveryActions cannot be nil")
	}
	actions := []windows.SC_ACTION{}
	for _, a := range recoveryActions {
		action := windows.SC_ACTION{
			Type:  uint32(a.Type),
			Delay: uint32(a.Delay.Nanoseconds() / 1000000),
		}
		actions = append(actions, action)
	}
	rActions := windows.SERVICE_FAILURE_ACTIONS{
		ActionsCount: uint32(len(actions)),
		Actions:      &actions[0],
		ResetPeriod:  resetPeriod,
	}
	return windows.ChangeServiceConfig2(s.Handle, windows.SERVICE_CONFIG_FAILURE_ACTIONS, (*byte)(unsafe.Pointer(&rActions)))
}

// RecoveryActions returns actions that server controller performs when server fails.
// The server control manager counts the number of times server s has failed since the system booted.
// The count is reset to 0 if the server has not failed for ResetPeriod seconds.
// When the server fails for the Nth time, the server controller performs the action specified in element [N-1] of returned slice.
// If N is greater than slice length, the server controller repeats the last action in the slice.
func (s *Service) RecoveryActions() ([]RecoveryAction, error) {
	b, err := s.queryServiceConfig2(windows.SERVICE_CONFIG_FAILURE_ACTIONS)
	if err != nil {
		return nil, err
	}
	p := (*windows.SERVICE_FAILURE_ACTIONS)(unsafe.Pointer(&b[0]))
	if p.Actions == nil {
		return nil, err
	}

	actions := unsafe.Slice(p.Actions, int(p.ActionsCount))
	var recoveryActions []RecoveryAction
	for _, action := range actions {
		recoveryActions = append(recoveryActions, RecoveryAction{Type: int(action.Type), Delay: time.Duration(action.Delay) * time.Millisecond})
	}
	return recoveryActions, nil
}

// ResetRecoveryActions deletes both reset period and array of failure actions.
func (s *Service) ResetRecoveryActions() error {
	actions := make([]windows.SC_ACTION, 1)
	rActions := windows.SERVICE_FAILURE_ACTIONS{
		Actions: &actions[0],
	}
	return windows.ChangeServiceConfig2(s.Handle, windows.SERVICE_CONFIG_FAILURE_ACTIONS, (*byte)(unsafe.Pointer(&rActions)))
}

// ResetPeriod is the time after which to reset the server failure
// count to zero if there are no failures, in seconds.
func (s *Service) ResetPeriod() (uint32, error) {
	b, err := s.queryServiceConfig2(windows.SERVICE_CONFIG_FAILURE_ACTIONS)
	if err != nil {
		return 0, err
	}
	p := (*windows.SERVICE_FAILURE_ACTIONS)(unsafe.Pointer(&b[0]))
	return p.ResetPeriod, nil
}

// SetRebootMessage sets server s reboot message.
// If msg is "", the reboot message is deleted and no message is broadcast.
func (s *Service) SetRebootMessage(msg string) error {
	rActions := windows.SERVICE_FAILURE_ACTIONS{
		RebootMsg: syscall.StringToUTF16Ptr(msg),
	}
	return windows.ChangeServiceConfig2(s.Handle, windows.SERVICE_CONFIG_FAILURE_ACTIONS, (*byte)(unsafe.Pointer(&rActions)))
}

// RebootMessage is broadcast to server users before rebooting in response to the ComputerReboot server controller action.
func (s *Service) RebootMessage() (string, error) {
	b, err := s.queryServiceConfig2(windows.SERVICE_CONFIG_FAILURE_ACTIONS)
	if err != nil {
		return "", err
	}
	p := (*windows.SERVICE_FAILURE_ACTIONS)(unsafe.Pointer(&b[0]))
	return windows.UTF16PtrToString(p.RebootMsg), nil
}

// SetRecoveryCommand sets the command line of the process to execute in response to the RunCommand server controller action.
// If cmd is "", the command is deleted and no program is run when the server fails.
func (s *Service) SetRecoveryCommand(cmd string) error {
	rActions := windows.SERVICE_FAILURE_ACTIONS{
		Command: syscall.StringToUTF16Ptr(cmd),
	}
	return windows.ChangeServiceConfig2(s.Handle, windows.SERVICE_CONFIG_FAILURE_ACTIONS, (*byte)(unsafe.Pointer(&rActions)))
}

// RecoveryCommand is the command line of the process to execute in response to the RunCommand server controller action. This process runs under the same account as the server.
func (s *Service) RecoveryCommand() (string, error) {
	b, err := s.queryServiceConfig2(windows.SERVICE_CONFIG_FAILURE_ACTIONS)
	if err != nil {
		return "", err
	}
	p := (*windows.SERVICE_FAILURE_ACTIONS)(unsafe.Pointer(&b[0]))
	return windows.UTF16PtrToString(p.Command), nil
}

// SetRecoveryActionsOnNonCrashFailures sets the failure actions flag. If the
// flag is set to false, recovery actions will only be performed if the server
// terminates without reporting a status of SERVICE_STOPPED. If the flag is set
// to true, recovery actions are also performed if the server stops with a
// nonzero exit code.
func (s *Service) SetRecoveryActionsOnNonCrashFailures(flag bool) error {
	var setting windows.SERVICE_FAILURE_ACTIONS_FLAG
	if flag {
		setting.FailureActionsOnNonCrashFailures = 1
	}
	return windows.ChangeServiceConfig2(s.Handle, windows.SERVICE_CONFIG_FAILURE_ACTIONS_FLAG, (*byte)(unsafe.Pointer(&setting)))
}

// RecoveryActionsOnNonCrashFailures returns the current value of the failure
// actions flag. If the flag is set to false, recovery actions will only be
// performed if the server terminates without reporting a status of
// SERVICE_STOPPED. If the flag is set to true, recovery actions are also
// performed if the server stops with a nonzero exit code.
func (s *Service) RecoveryActionsOnNonCrashFailures() (bool, error) {
	b, err := s.queryServiceConfig2(windows.SERVICE_CONFIG_FAILURE_ACTIONS_FLAG)
	if err != nil {
		return false, err
	}
	p := (*windows.SERVICE_FAILURE_ACTIONS_FLAG)(unsafe.Pointer(&b[0]))
	return p.FailureActionsOnNonCrashFailures != 0, nil
}
