//go:build !integration

package app

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	cfg := (&Config{}).SetDefault()

	_, err := New(cfg)

	require.NoError(t, err)
}

// mockShutdownAwareService is a test service that implements ShutdownAware
type mockShutdownAwareService struct {
	shutdownRequester ShutdownRequester
	initCalled        bool
	startCalled       bool
	stopCalled        bool
	closeCalled       bool
	triggerShutdown   bool
	shutdownReason    string
	shutdownErr       error
	shutdownSource    string
}

func (m *mockShutdownAwareService) SetShutdownRequester(sr ShutdownRequester) {
	m.shutdownRequester = sr
}

func (m *mockShutdownAwareService) Init(_ context.Context) error {
	m.initCalled = true
	return nil
}

func (m *mockShutdownAwareService) Start(_ context.Context) error {
	m.startCalled = true

	// If configured to trigger shutdown, do it in a goroutine after start
	if m.triggerShutdown {
		go func() {
			time.Sleep(50 * time.Millisecond)
			m.shutdownRequester.RequestShutdown(ShutdownRequest{
				Reason: m.shutdownReason,
				Err:    m.shutdownErr,
				Source: m.shutdownSource,
			})
		}()
	}

	return nil
}

func (m *mockShutdownAwareService) Stop(_ context.Context) error {
	m.stopCalled = true
	return nil
}

func (m *mockShutdownAwareService) Close() error {
	m.closeCalled = true
	return nil
}

func TestShutdownRequester(t *testing.T) {
	t.Run("ShutdownRequester is injected into ShutdownAware services", func(t *testing.T) {
		cfg := (&Config{}).SetDefault()
		cfg.StartTimeout.Duration = 5 * time.Second
		cfg.StopTimeout.Duration = 5 * time.Second
		// Use random available ports to avoid conflicts
		cfg.Server.Entrypoint.Address = ":0"
		cfg.Healthz.Entrypoint.Address = ":0"

		app, err := New(cfg)
		require.NoError(t, err)

		svc := &mockShutdownAwareService{}
		app.RegisterService(svc)

		// Start app in a goroutine
		done := make(chan error, 1)
		go func() {
			done <- app.Run()
		}()

		// Wait a bit for initialization
		time.Sleep(200 * time.Millisecond)

		// Verify the shutdown requester was set
		require.NotNil(t, svc.shutdownRequester, "ShutdownRequester should be injected")

		// Send shutdown signal
		app.done <- syscall.SIGTERM

		// Wait for app to finish
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(10 * time.Second):
			t.Fatal("timeout waiting for app to stop")
		}

		// Verify lifecycle methods were called
		assert.True(t, svc.initCalled, "Init should be called")
		assert.True(t, svc.startCalled, "Start should be called")
		assert.True(t, svc.stopCalled, "Stop should be called")
		assert.True(t, svc.closeCalled, "Close should be called")
	})

	t.Run("Service can request shutdown via ShutdownRequester", func(t *testing.T) {
		cfg := (&Config{}).SetDefault()
		cfg.StartTimeout.Duration = 5 * time.Second
		cfg.StopTimeout.Duration = 5 * time.Second
		// Use random available ports to avoid conflicts
		cfg.Server.Entrypoint.Address = ":0"
		cfg.Healthz.Entrypoint.Address = ":0"

		app, err := New(cfg)
		require.NoError(t, err)

		testErr := errors.New("test error")
		svc := &mockShutdownAwareService{
			triggerShutdown: true,
			shutdownReason:  "service requested shutdown",
			shutdownErr:     testErr,
			shutdownSource:  "mockService",
		}
		app.RegisterService(svc)

		// Start app and wait for it to finish
		start := time.Now()
		err = app.Run()
		duration := time.Since(start)

		// Should complete without error
		require.NoError(t, err)

		// Should complete relatively quickly (not timeout)
		assert.Less(t, duration, 2*time.Second, "app should shutdown quickly when requested")

		// Verify all lifecycle methods were called
		assert.True(t, svc.initCalled, "Init should be called")
		assert.True(t, svc.startCalled, "Start should be called")
		assert.True(t, svc.stopCalled, "Stop should be called")
		assert.True(t, svc.closeCalled, "Close should be called")
	})

	t.Run("Multiple shutdown requests are handled gracefully", func(t *testing.T) {
		cfg := (&Config{}).SetDefault()
		cfg.StartTimeout.Duration = 5 * time.Second
		cfg.StopTimeout.Duration = 5 * time.Second
		// Use random available ports to avoid conflicts
		cfg.Server.Entrypoint.Address = ":0"
		cfg.Healthz.Entrypoint.Address = ":0"

		app, err := New(cfg)
		require.NoError(t, err)

		svc := &mockShutdownAwareService{
			triggerShutdown: false,
		}
		app.RegisterService(svc)

		// Start app in goroutine
		done := make(chan error, 1)
		go func() {
			done <- app.Run()
		}()

		// Wait for service to start
		time.Sleep(200 * time.Millisecond)

		// Request multiple shutdowns
		if svc.shutdownRequester != nil {
			go func() {
				for i := 0; i < 5; i++ {
					svc.shutdownRequester.RequestShutdown(ShutdownRequest{
						Reason: fmt.Sprintf("shutdown %d", i),
						Source: "test",
					})
					time.Sleep(10 * time.Millisecond)
				}
			}()
		}

		// Wait for app to finish
		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(10 * time.Second):
			t.Fatal("timeout waiting for app to stop")
		}

		// App should have stopped gracefully
		assert.True(t, svc.stopCalled, "Stop should be called")
		assert.True(t, svc.closeCalled, "Close should be called")
	})

	t.Run("ShutdownFunc adapter works correctly", func(t *testing.T) {
		var receivedRequest ShutdownRequest
		fn := ShutdownFunc(func(req ShutdownRequest) {
			receivedRequest = req
		})

		testErr := errors.New("test error")
		expectedReq := ShutdownRequest{
			Reason: "test reason",
			Err:    testErr,
			Source: "test source",
		}

		fn.RequestShutdown(expectedReq)

		assert.Equal(t, expectedReq.Reason, receivedRequest.Reason)
		assert.Equal(t, expectedReq.Err, receivedRequest.Err)
		assert.Equal(t, expectedReq.Source, receivedRequest.Source)
	})
}
