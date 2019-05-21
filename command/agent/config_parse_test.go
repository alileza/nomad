package agent

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/nomad/structs/config"
	"github.com/stretchr/testify/require"
)

var basicConfig = &Config{
	Region:      "foobar",
	Datacenter:  "dc2",
	NodeName:    "my-web",
	DataDir:     "/tmp/nomad",
	PluginDir:   "/tmp/nomad-plugins",
	LogLevel:    "ERR",
	LogJson:     true,
	BindAddr:    "192.168.0.1",
	EnableDebug: true,
	Ports: &Ports{
		HTTP: 1234,
		RPC:  2345,
		Serf: 3456,
	},
	Addresses: &Addresses{
		HTTP: "127.0.0.1",
		RPC:  "127.0.0.2",
		Serf: "127.0.0.3",
	},
	AdvertiseAddrs: &AdvertiseAddrs{
		RPC:  "127.0.0.3",
		Serf: "127.0.0.4",
	},
	Client: &ClientConfig{
		Enabled:   true,
		StateDir:  "/tmp/client-state",
		AllocDir:  "/tmp/alloc",
		Servers:   []string{"a.b.c:80", "127.0.0.1:1234"},
		NodeClass: "linux-medium-64bit",
		ServerJoin: &ServerJoin{
			RetryJoin:        []string{"1.1.1.1", "2.2.2.2"},
			RetryInterval:    time.Duration(15) * time.Second,
			RetryIntervalHCL: "15s",
			RetryMaxAttempts: 3,
		},
		Meta: map[string]string{
			"foo": "bar",
			"baz": "zip",
		},
		Options: map[string]string{
			"foo": "bar",
			"baz": "zip",
		},
		ChrootEnv: map[string]string{
			"/opt/myapp/etc": "/etc",
			"/opt/myapp/bin": "/bin",
		},
		NetworkInterface: "eth0",
		NetworkSpeed:     100,
		CpuCompute:       4444,
		MemoryMB:         0,
		MaxKillTimeout:   "10s",
		ClientMinPort:    1000,
		ClientMaxPort:    2000,
		Reserved: &Resources{
			CPU:           10,
			MemoryMB:      10,
			DiskMB:        10,
			ReservedPorts: "1,100,10-12",
		},
		GCInterval:            6 * time.Second,
		GCIntervalHCL:         "6s",
		GCParallelDestroys:    6,
		GCDiskUsageThreshold:  82,
		GCInodeUsageThreshold: 91,
		GCMaxAllocs:           50,
		NoHostUUID:            helper.BoolToPtr(false),
		HostVolumes: []*structs.ClientHostVolumeConfig{
			{Name: "tmp", Source: "/tmp"},
		},
	},
	Server: &ServerConfig{
		Enabled:                true,
		AuthoritativeRegion:    "foobar",
		BootstrapExpect:        5,
		DataDir:                "/tmp/data",
		ProtocolVersion:        3,
		RaftProtocol:           3,
		NumSchedulers:          helper.IntToPtr(2),
		EnabledSchedulers:      []string{"test"},
		NodeGCThreshold:        "12h",
		EvalGCThreshold:        "12h",
		JobGCThreshold:         "12h",
		DeploymentGCThreshold:  "12h",
		HeartbeatGrace:         30 * time.Second,
		HeartbeatGraceHCL:      "30s",
		MinHeartbeatTTL:        33 * time.Second,
		MinHeartbeatTTLHCL:     "33s",
		MaxHeartbeatsPerSecond: 11.0,
		RetryJoin:              []string{"1.1.1.1", "2.2.2.2"},
		StartJoin:              []string{"1.1.1.1", "2.2.2.2"},
		RetryInterval:          15 * time.Second,
		RetryIntervalHCL:       "15s",
		RejoinAfterLeave:       true,
		RetryMaxAttempts:       3,
		NonVotingServer:        true,
		RedundancyZone:         "foo",
		UpgradeVersion:         "0.8.0",
		EncryptKey:             "abc",
		ServerJoin: &ServerJoin{
			RetryJoin:        []string{"1.1.1.1", "2.2.2.2"},
			RetryInterval:    time.Duration(15) * time.Second,
			RetryIntervalHCL: "15s",
			RetryMaxAttempts: 3,
		},
	},
	ACL: &ACLConfig{
		Enabled:          true,
		TokenTTL:         60 * time.Second,
		TokenTTLHCL:      "60s",
		PolicyTTL:        60 * time.Second,
		PolicyTTLHCL:     "60s",
		ReplicationToken: "foobar",
	},
	Telemetry: &Telemetry{
		StatsiteAddr:               "127.0.0.1:1234",
		StatsdAddr:                 "127.0.0.1:2345",
		PrometheusMetrics:          true,
		DisableHostname:            true,
		UseNodeName:                false,
		CollectionInterval:         "3s",
		collectionInterval:         3 * time.Second,
		PublishAllocationMetrics:   true,
		PublishNodeMetrics:         true,
		DisableTaggedMetrics:       true,
		BackwardsCompatibleMetrics: true,
	},
	LeaveOnInt:                true,
	LeaveOnTerm:               true,
	EnableSyslog:              true,
	SyslogFacility:            "LOCAL1",
	DisableUpdateCheck:        helper.BoolToPtr(true),
	DisableAnonymousSignature: true,
	Consul: &config.ConsulConfig{
		ServerServiceName:   "nomad",
		ServerHTTPCheckName: "nomad-server-http-health-check",
		ServerSerfCheckName: "nomad-server-serf-health-check",
		ServerRPCCheckName:  "nomad-server-rpc-health-check",
		ClientServiceName:   "nomad-client",
		ClientHTTPCheckName: "nomad-client-http-health-check",
		Addr:                "127.0.0.1:9500",
		Token:               "token1",
		Auth:                "username:pass",
		EnableSSL:           &trueValue,
		VerifySSL:           &trueValue,
		CAFile:              "/path/to/ca/file",
		CertFile:            "/path/to/cert/file",
		KeyFile:             "/path/to/key/file",
		ServerAutoJoin:      &trueValue,
		ClientAutoJoin:      &trueValue,
		AutoAdvertise:       &trueValue,
		ChecksUseAdvertise:  &trueValue,
		Timeout:             5 * time.Second,
	},
	Vault: &config.VaultConfig{
		Addr:                 "127.0.0.1:9500",
		AllowUnauthenticated: &trueValue,
		ConnectionRetryIntv:  config.DefaultVaultConnectRetryIntv,
		Enabled:              &falseValue,
		Role:                 "test_role",
		TLSCaFile:            "/path/to/ca/file",
		TLSCaPath:            "/path/to/ca",
		TLSCertFile:          "/path/to/cert/file",
		TLSKeyFile:           "/path/to/key/file",
		TLSServerName:        "foobar",
		TLSSkipVerify:        &trueValue,
		TaskTokenTTL:         "1s",
		Token:                "12345",
	},
	TLSConfig: &config.TLSConfig{
		EnableHTTP:                  true,
		EnableRPC:                   true,
		VerifyServerHostname:        true,
		CAFile:                      "foo",
		CertFile:                    "bar",
		KeyFile:                     "pipe",
		RPCUpgradeMode:              true,
		VerifyHTTPSClient:           true,
		TLSPreferServerCipherSuites: true,
		TLSCipherSuites:             "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		TLSMinVersion:               "tls12",
	},
	HTTPAPIResponseHeaders: map[string]string{
		"Access-Control-Allow-Origin": "*",
	},
	Sentinel: &config.SentinelConfig{
		Imports: []*config.SentinelImport{
			{
				Name: "foo",
				Path: "foo",
				Args: []string{"a", "b", "c"},
			},
			{
				Name: "bar",
				Path: "bar",
				Args: []string{"x", "y", "z"},
			},
		},
	},
	Autopilot: &config.AutopilotConfig{
		CleanupDeadServers:         &trueValue,
		ServerStabilizationTime:    23057 * time.Second,
		ServerStabilizationTimeHCL: "23057s",
		LastContactThreshold:       12705 * time.Second,
		LastContactThresholdHCL:    "12705s",
		MaxTrailingLogs:            17849,
		EnableRedundancyZones:      &trueValue,
		DisableUpgradeMigration:    &trueValue,
		EnableCustomUpgrades:       &trueValue,
	},
	Plugins: []*config.PluginConfig{
		{
			Name: "docker",
			Args: []string{"foo", "bar"},
			Config: map[string]interface{}{
				"foo": "bar",
				"nested": []map[string]interface{}{
					{
						"bam": 2,
					},
				},
			},
		},
		{
			Name: "exec",
			Config: map[string]interface{}{
				"foo": true,
			},
		},
	},
}

var pluginConfig = &Config{
	Region:         "",
	Datacenter:     "",
	NodeName:       "",
	DataDir:        "",
	PluginDir:      "",
	LogLevel:       "",
	BindAddr:       "",
	EnableDebug:    false,
	Ports:          nil,
	Addresses:      nil,
	AdvertiseAddrs: nil,
	Client: &ClientConfig{
		Enabled:               false,
		StateDir:              "",
		AllocDir:              "",
		Servers:               nil,
		NodeClass:             "",
		Meta:                  nil,
		Options:               nil,
		ChrootEnv:             nil,
		NetworkInterface:      "",
		NetworkSpeed:          0,
		CpuCompute:            0,
		MemoryMB:              5555,
		MaxKillTimeout:        "",
		ClientMinPort:         0,
		ClientMaxPort:         0,
		Reserved:              nil,
		GCInterval:            0,
		GCParallelDestroys:    0,
		GCDiskUsageThreshold:  0,
		GCInodeUsageThreshold: 0,
		GCMaxAllocs:           0,
		NoHostUUID:            nil,
	},
	Server:                    nil,
	ACL:                       nil,
	Telemetry:                 nil,
	LeaveOnInt:                false,
	LeaveOnTerm:               false,
	EnableSyslog:              false,
	SyslogFacility:            "",
	DisableUpdateCheck:        nil,
	DisableAnonymousSignature: false,
	Consul:                    nil,
	Vault:                     nil,
	TLSConfig:                 nil,
	HTTPAPIResponseHeaders:    nil,
	Sentinel:                  nil,
	Plugins: []*config.PluginConfig{
		{
			Name: "docker",
			Config: map[string]interface{}{
				"allow_privileged": true,
			},
		},
		{
			Name: "raw_exec",
			Config: map[string]interface{}{
				"enabled": true,
			},
		},
	},
}

var nonoptConfig = &Config{
	Region:         "",
	Datacenter:     "",
	NodeName:       "",
	DataDir:        "",
	PluginDir:      "",
	LogLevel:       "",
	BindAddr:       "",
	EnableDebug:    false,
	Ports:          nil,
	Addresses:      nil,
	AdvertiseAddrs: nil,
	Client: &ClientConfig{
		Enabled:               false,
		StateDir:              "",
		AllocDir:              "",
		Servers:               nil,
		NodeClass:             "",
		Meta:                  nil,
		Options:               nil,
		ChrootEnv:             nil,
		NetworkInterface:      "",
		NetworkSpeed:          0,
		CpuCompute:            0,
		MemoryMB:              5555,
		MaxKillTimeout:        "",
		ClientMinPort:         0,
		ClientMaxPort:         0,
		Reserved:              nil,
		GCInterval:            0,
		GCParallelDestroys:    0,
		GCDiskUsageThreshold:  0,
		GCInodeUsageThreshold: 0,
		GCMaxAllocs:           0,
		NoHostUUID:            nil,
	},
	Server:                    nil,
	ACL:                       nil,
	Telemetry:                 nil,
	LeaveOnInt:                false,
	LeaveOnTerm:               false,
	EnableSyslog:              false,
	SyslogFacility:            "",
	DisableUpdateCheck:        nil,
	DisableAnonymousSignature: false,
	Consul:                    nil,
	Vault:                     nil,
	TLSConfig:                 nil,
	HTTPAPIResponseHeaders:    nil,
	Sentinel:                  nil,
}

func TestConfig_Parse(t *testing.T) {
	t.Parallel()

	basicConfig.addDefaults()
	pluginConfig.addDefaults()
	nonoptConfig.addDefaults()

	cases := []struct {
		File   string
		Result *Config
		Err    bool
	}{
		{
			"basic.hcl",
			basicConfig,
			false,
		},
		{
			"basic.json",
			basicConfig,
			false,
		},
		{
			"plugin.hcl",
			pluginConfig,
			false,
		},
		{
			"plugin.json",
			pluginConfig,
			false,
		},
		{
			"non-optional.hcl",
			nonoptConfig,
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.File, func(t *testing.T) {
			require := require.New(t)
			path, err := filepath.Abs(filepath.Join("./testdata", tc.File))
			if err != nil {
				t.Fatalf("file: %s\n\n%s", tc.File, err)
			}

			actual, err := ParseConfigFile(path)
			if (err != nil) != tc.Err {
				t.Fatalf("file: %s\n\n%s", tc.File, err)
			}

			//panic(fmt.Sprintf("first: %+v \n second: %+v", actual.TLSConfig, tc.Result.TLSConfig))
			require.EqualValues(tc.Result, removeHelperAttributes(actual))
		})
	}
}

// In order to compare the Config struct after parsing, and from generating what
// is expected in the test, we need to remove helper attributes that are
// instantiated in the process of parsing the configuration
func removeHelperAttributes(c *Config) *Config {
	if c.TLSConfig != nil {
		c.TLSConfig.KeyLoader = nil
	}
	return c
}

func (c *Config) addDefaults() {
	if c.Client == nil {
		c.Client = &ClientConfig{}
	}
	if c.Client.ServerJoin == nil {
		c.Client.ServerJoin = &ServerJoin{}
	}
	if c.ACL == nil {
		c.ACL = &ACLConfig{}
	}
	if c.Consul == nil {
		c.Consul = config.DefaultConsulConfig()
	}
	if c.Autopilot == nil {
		c.Autopilot = config.DefaultAutopilotConfig()
	}
	if c.Vault == nil {
		c.Vault = config.DefaultVaultConfig()
	}
	if c.Telemetry == nil {
		c.Telemetry = &Telemetry{}
	}
	if c.Server == nil {
		c.Server = &ServerConfig{}
	}
	if c.Server.ServerJoin == nil {
		c.Server.ServerJoin = &ServerJoin{}
	}
}

// Tests for a panic parsing json with an object of exactly
// length 1 described in
// https://github.com/hashicorp/nomad/issues/1290
func TestConfig_ParsePanic(t *testing.T) {
	c, err := ParseConfigFile("./testdata/obj-len-one.hcl")
	if err != nil {
		t.Fatalf("parse error: %s\n", err)
	}

	d, err := ParseConfigFile("./testdata/obj-len-one.json")
	if err != nil {
		t.Fatalf("parse error: %s\n", err)
	}

	require.EqualValues(t, c, d)
}

// Top level keys left by hcl when parsing slices in the config
// structure should not be unexpected
func TestConfig_ParseSliceExtra(t *testing.T) {
	c, err := ParseConfigFile("./testdata/config-slices.json")
	if err != nil {
		t.Fatalf("parse error: %s\n", err)
	}

	opt := map[string]string{"o0": "foo", "o1": "bar"}
	meta := map[string]string{"m0": "foo", "m1": "bar"}
	env := map[string]string{"e0": "baz"}
	srv := []string{"foo", "bar"}

	require.EqualValues(t, opt, c.Client.Options)
	require.EqualValues(t, meta, c.Client.Meta)
	require.EqualValues(t, env, c.Client.ChrootEnv)
	require.EqualValues(t, srv, c.Client.Servers)
	require.EqualValues(t, srv, c.Server.EnabledSchedulers)
	require.EqualValues(t, srv, c.Server.StartJoin)
	require.EqualValues(t, srv, c.Server.RetryJoin)

	// the alt format is also accepted by hcl as valid config data
	c, err = ParseConfigFile("./testdata/config-slices-alt.json")
	if err != nil {
		t.Fatalf("parse error: %s\n", err)
	}

	require.EqualValues(t, opt, c.Client.Options)
	require.EqualValues(t, meta, c.Client.Meta)
	require.EqualValues(t, env, c.Client.ChrootEnv)
	require.EqualValues(t, srv, c.Client.Servers)
	require.EqualValues(t, srv, c.Server.EnabledSchedulers)
	require.EqualValues(t, srv, c.Server.StartJoin)
	require.EqualValues(t, srv, c.Server.RetryJoin)

	// small files keep more extra keys than large ones
	_, err = ParseConfigFile("./testdata/obj-len-one-server.json")
	if err != nil {
		t.Fatalf("parse error: %s\n", err)
	}
}

var sample0 = &Config{
	Region:     "global",
	Datacenter: "dc1",
	DataDir:    "/opt/data/nomad/data",
	LogLevel:   "INFO",
	BindAddr:   "0.0.0.0",
	AdvertiseAddrs: &AdvertiseAddrs{
		HTTP: "host.example.com",
		RPC:  "host.example.com",
		Serf: "host.example.com",
	},
	Client: &ClientConfig{ServerJoin: &ServerJoin{}},
	Server: &ServerConfig{
		Enabled:         true,
		BootstrapExpect: 3,
		RetryJoin:       []string{"10.0.0.101", "10.0.0.102", "10.0.0.103"},
		EncryptKey:      "sHck3WL6cxuhuY7Mso9BHA==",
		ServerJoin:      &ServerJoin{},
	},
	ACL: &ACLConfig{
		Enabled: true,
	},
	Telemetry: &Telemetry{
		PrometheusMetrics:        true,
		DisableHostname:          true,
		CollectionInterval:       "60s",
		collectionInterval:       60 * time.Second,
		PublishAllocationMetrics: true,
		PublishNodeMetrics:       true,
	},
	LeaveOnInt:     true,
	LeaveOnTerm:    true,
	EnableSyslog:   true,
	SyslogFacility: "LOCAL0",
	Consul: &config.ConsulConfig{
		Token:          "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		ServerAutoJoin: helper.BoolToPtr(false),
		ClientAutoJoin: helper.BoolToPtr(false),
		// Defaults
		ServerServiceName:   "nomad",
		ServerHTTPCheckName: "Nomad Server HTTP Check",
		ServerSerfCheckName: "Nomad Server Serf Check",
		ServerRPCCheckName:  "Nomad Server RPC Check",
		ClientServiceName:   "nomad-client",
		ClientHTTPCheckName: "Nomad Client HTTP Check",
		AutoAdvertise:       helper.BoolToPtr(true),
		ChecksUseAdvertise:  helper.BoolToPtr(false),
		Timeout:             5 * time.Second,
		EnableSSL:           helper.BoolToPtr(false),
		VerifySSL:           helper.BoolToPtr(true),
	},
	Vault: &config.VaultConfig{
		Enabled: helper.BoolToPtr(true),
		Role:    "nomad-cluster",
		Addr:    "http://host.example.com:8200",
		// Defaults
		AllowUnauthenticated: helper.BoolToPtr(true),
		ConnectionRetryIntv:  30 * time.Second,
	},
	TLSConfig: &config.TLSConfig{
		EnableHTTP:           true,
		EnableRPC:            true,
		VerifyServerHostname: true,
		CAFile:               "/opt/data/nomad/certs/nomad-ca.pem",
		CertFile:             "/opt/data/nomad/certs/server.pem",
		KeyFile:              "/opt/data/nomad/certs/server-key.pem",
	},
	Autopilot: &config.AutopilotConfig{
		CleanupDeadServers: helper.BoolToPtr(true),
		// Defaults
		ServerStabilizationTime: 10 * time.Second,
		LastContactThreshold:    200 * time.Millisecond,
		MaxTrailingLogs:         250,
	},
}

func TestConfig_ParseSample0(t *testing.T) {
	c, err := ParseConfigFile("./testdata/sample0.json")
	require.Nil(t, err)
	require.EqualValues(t, sample0, c)
}
