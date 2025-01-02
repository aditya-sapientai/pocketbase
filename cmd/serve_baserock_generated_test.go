package cmd

import (
	"errors"
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockApp struct {
	mock.Mock
}

func (m *mockApp) Bootstrap() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockApp) ResetBootstrapState() {
	m.Called()
}

func (m *mockApp) IsBootstrapped() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockApp) Cleanup() {
	m.Called()
}

func TestNewServeCommand(t *testing.T) {
	mockApp := &mockApp{}
	cmd := NewServeCommand(mockApp, true)

	assert.NotNil(t, cmd)
	assert.Equal(t, "serve", cmd.Use)
	assert.Equal(t, "Starts the web server (default to 127.0.0.1:8090 if no domain is specified)", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
}

func TestServeCommandRunE(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		httpAddr       string
		httpsAddr      string
		allowedOrigins []string
		expectedConfig apis.ServeConfig
		serveError     error
		expectedError  error
	}{
		{
			name: "Default configuration",
			args: []string{},
			expectedConfig: apis.ServeConfig{
				HttpAddr:        "127.0.0.1:8090",
				ShowStartBanner: true,
				AllowedOrigins:  []string{"*"},
			},
		},
		{
			name: "With domain",
			args: []string{"example.com"},
			expectedConfig: apis.ServeConfig{
				HttpAddr:           "0.0.0.0:80",
				HttpsAddr:          "0.0.0.0:443",
				ShowStartBanner:    true,
				AllowedOrigins:     []string{"*"},
				CertificateDomains: []string{"example.com"},
			},
		},
		{
			name:     "With custom HTTP address",
			args:     []string{},
			httpAddr: "0.0.0.0:8080",
			expectedConfig: apis.ServeConfig{
				HttpAddr:        "0.0.0.0:8080",
				ShowStartBanner: true,
				AllowedOrigins:  []string{"*"},
			},
		},
		{
			name:      "With custom HTTPS address",
			args:      []string{"example.com"},
			httpsAddr: "0.0.0.0:8443",
			expectedConfig: apis.ServeConfig{
				HttpAddr:           "0.0.0.0:80",
				HttpsAddr:          "0.0.0.0:8443",
				ShowStartBanner:    true,
				AllowedOrigins:     []string{"*"},
				CertificateDomains: []string{"example.com"},
			},
		},
		{
			name:           "With custom allowed origins",
			args:           []string{},
			allowedOrigins: []string{"http://localhost:3000", "https://example.com"},
			expectedConfig: apis.ServeConfig{
				HttpAddr:        "127.0.0.1:8090",
				ShowStartBanner: true,
				AllowedOrigins:  []string{"http://localhost:3000", "https://example.com"},
			},
		},
		{
			name:          "Server closed error",
			args:          []string{},
			serveError:    http.ErrServerClosed,
			expectedError: nil,
		},
		{
			name:          "Other error",
			args:          []string{},
			serveError:    errors.New("test error"),
			expectedError: errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockApp := &mockApp{}
			cmd := NewServeCommand(mockApp, true)

			cmd.SetArgs(tt.args)
			if tt.httpAddr != "" {
				cmd.PersistentFlags().Set("http", tt.httpAddr)
			}
			if tt.httpsAddr != "" {
				cmd.PersistentFlags().Set("https", tt.httpsAddr)
			}
			if len(tt.allowedOrigins) > 0 {
				cmd.PersistentFlags().Set("origins", tt.allowedOrigins[0])
				for _, origin := range tt.allowedOrigins[1:] {
					cmd.PersistentFlags().Set("origins", origin)
				}
			}

			originalServe := apis.Serve
			defer func() { apis.Serve = originalServe }()

			apis.Serve = func(app core.App, config apis.ServeConfig) error {
				assert.Equal(t, tt.expectedConfig, config)
				return tt.serveError
			}

			err := cmd.Execute()

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServeCommandFlags(t *testing.T) {
	mockApp := &mockApp{}
	cmd := NewServeCommand(mockApp, true)

	flags := cmd.PersistentFlags()

	originFlag := flags.Lookup("origins")
	assert.NotNil(t, originFlag)
	assert.Equal(t, "string", originFlag.Value.Type())
	assert.Equal(t, "[*]", originFlag.DefValue)

	httpFlag := flags.Lookup("http")
	assert.NotNil(t, httpFlag)
	assert.Equal(t, "string", httpFlag.Value.Type())
	assert.Equal(t, "", httpFlag.DefValue)

	httpsFlag := flags.Lookup("https")
	assert.NotNil(t, httpsFlag)
	assert.Equal(t, "string", httpsFlag.Value.Type())
	assert.Equal(t, "", httpsFlag.DefValue)
}
