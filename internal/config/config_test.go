package config

// *** THIS IS GENERATED CODE: DO NOT EDIT ***

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	falseVal = false
	trueVal  = true
)

func TestDuration_String(t *testing.T) {
	f := func(d time.Duration, exp string) {
		actual := Duration(d).String()
		if actual != exp {
			t.Errorf("String() for duration %d expected to return %s, but got %s", d, exp, actual)
		}
	}
	f(time.Second, "1s")
	f(time.Second*30, "30s")
	f(time.Minute, "1m0s")
	f(time.Second*90, "1m30s")
	f(0, "0s")
}

func TestDuration_JSON(t *testing.T) {
	f := func(d time.Duration, exp string) {
		v := Duration(d)
		bytes, err := json.Marshal(&v)
		if err != nil {
			t.Fatalf("Unable to json.Marshal our Duration of %+v: %v", v, err)
		}
		if string(bytes) != exp {
			t.Errorf("Marshaled duration expected to generate %v, but got %v", exp, string(bytes))
		}
		var decoded Duration
		if err := json.Unmarshal(bytes, &decoded); err != nil {
			t.Errorf("Got error trying to unmarshal %v to a Duration: %v", string(bytes), err)
		}
		if decoded != v {
			t.Errorf("Encoded/Decoded duration no longer equal!, original %v, round-tripped %v", v, decoded)
		}
	}
	f(0, `"0s"`)
	f(time.Second, `"1s"`)
	f(time.Minute*5, `"5m0s"`)
	f(time.Second*90, `"1m30s"`)
	f(time.Hour*2, `"2h0m0s"`)
	f(time.Millisecond*10, `"10ms"`)
}

func TestDuration_JSONDecode(t *testing.T) {
	f := func(j string, exp time.Duration) {
		var act Duration
		err := json.Unmarshal([]byte(j), &act)
		if err != nil {
			t.Fatalf("Unable to json.Unmarshal %s: %v", j, err)
		}
		if act.TimeDuration() != exp {
			t.Errorf("Expecting json of %s to production duration %s, but got %s", j, exp, act)
		}
	}
	f(`"5m"`, time.Minute*5)
	f(`120`, time.Second*120)
	f(`0`, 0)
	f(`"1m5s"`, time.Second*65)
}

func Test_overrideBool(t *testing.T) {
	d := &trueVal
	var zero *bool
	overrideBool(&d, &zero)
	require.NotEqual(t, d, zero, "overrideBool shouldn't have overriden the value as the override is the default/zero value. value now %v", d)
	o := &falseVal
	overrideBool(&d, &o)
	require.Equal(t, d, o, "overrideBool should of overriden the value but didn't. value %v, expecting %v", d, o)
}

func Test_overrideInt(t *testing.T) {
	d := -42
	var zero int
	overrideInt(&d, &zero)
	require.NotEqual(t, d, zero, "overrideInt shouldn't have overriden the value as the override is the default/zero value. value now %v", d)
	o := 42
	overrideInt(&d, &o)
	require.Equal(t, d, o, "overrideInt should of overriden the value but didn't. value %v, expecting %v", d, o)
}

func Test_overrideRepoLogLevelSlice(t *testing.T) {
	d := []RepoLogLevel{
		{
			Repo:    "one",
			Package: "one",
			Level:   "one"},
	}
	var zero []RepoLogLevel
	overrideRepoLogLevelSlice(&d, &zero)
	require.NotEqual(t, d, zero, "overrideRepoLogLevelSlice shouldn't have overriden the value as the override is the default/zero value. value now %v", d)
	o := []RepoLogLevel{
		{
			Repo:    "two",
			Package: "two",
			Level:   "two"},
	}
	overrideRepoLogLevelSlice(&d, &o)
	require.Equal(t, d, o, "overrideRepoLogLevelSlice should of overriden the value but didn't. value %v, expecting %v", d, o)
}

func Test_overrideString(t *testing.T) {
	d := "one"
	var zero string
	overrideString(&d, &zero)
	require.NotEqual(t, d, zero, "overrideString shouldn't have overriden the value as the override is the default/zero value. value now %v", d)
	o := "two"
	overrideString(&d, &o)
	require.Equal(t, d, o, "overrideString should of overriden the value but didn't. value %v, expecting %v", d, o)
}

func Test_overrideStrings(t *testing.T) {
	d := []string{"a"}
	var zero []string
	overrideStrings(&d, &zero)
	require.NotEqual(t, d, zero, "overrideStrings shouldn't have overriden the value as the override is the default/zero value. value now %v", d)
	o := []string{"b", "b"}
	overrideStrings(&d, &o)
	require.Equal(t, d, o, "overrideStrings should of overriden the value but didn't. value %v, expecting %v", d, o)
}

func TestAuthz_overrideFrom(t *testing.T) {
	orig := Authz{
		Allow:        []string{"a"},
		AllowAny:     []string{"a"},
		AllowAnyRole: []string{"a"},
		LogAllowed:   &trueVal,
		LogDenied:    &trueVal,
		CertMapper:   "one",
		APIKeyMapper: "one",
		JWTMapper:    "one"}
	dest := orig
	var zero Authz
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "Authz.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := Authz{
		Allow:        []string{"b", "b"},
		AllowAny:     []string{"b", "b"},
		AllowAnyRole: []string{"b", "b"},
		LogAllowed:   &falseVal,
		LogDenied:    &falseVal,
		CertMapper:   "two",
		APIKeyMapper: "two",
		JWTMapper:    "two"}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "Authz.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := Authz{
		Allow: []string{"a"}}
	dest.overrideFrom(&o2)
	exp := o

	exp.Allow = o2.Allow
	require.Equal(t, dest, exp, "Authz.overrideFrom should have overriden the field Allow. value now %#v, expecting %#v", dest, exp)
}

func TestAuthz_Getters(t *testing.T) {
	orig := Authz{
		Allow:        []string{"a"},
		AllowAny:     []string{"a"},
		AllowAnyRole: []string{"a"},
		LogAllowed:   &trueVal,
		LogDenied:    &trueVal,
		CertMapper:   "one",
		APIKeyMapper: "one",
		JWTMapper:    "one"}

	gv0 := orig.GetAllow()
	require.Equal(t, orig.Allow, gv0, "Authz.GetAllowCfg() does not match")

	gv1 := orig.GetAllowAny()
	require.Equal(t, orig.AllowAny, gv1, "Authz.GetAllowAnyCfg() does not match")

	gv2 := orig.GetAllowAnyRole()
	require.Equal(t, orig.AllowAnyRole, gv2, "Authz.GetAllowAnyRoleCfg() does not match")

	gv3 := orig.GetLogAllowed()
	require.Equal(t, orig.LogAllowed, &gv3, "Authz.GetLogAllowed() does not match")

	gv4 := orig.GetLogDenied()
	require.Equal(t, orig.LogDenied, &gv4, "Authz.GetLogDenied() does not match")

	gv5 := orig.GetCertMapper()
	require.Equal(t, orig.CertMapper, gv5, "Authz.GetCertMapperCfg() does not match")

	gv6 := orig.GetAPIKeyMapper()
	require.Equal(t, orig.APIKeyMapper, gv6, "Authz.GetAPIKeyMapperCfg() does not match")

	gv7 := orig.GetJWTMapper()
	require.Equal(t, orig.JWTMapper, gv7, "Authz.GetJWTMapperCfg() does not match")

}

func TestCORS_overrideFrom(t *testing.T) {
	orig := CORS{
		Enabled:            &trueVal,
		MaxAge:             -42,
		AllowedOrigins:     []string{"a"},
		AllowedMethods:     []string{"a"},
		AllowedHeaders:     []string{"a"},
		ExposedHeaders:     []string{"a"},
		AllowCredentials:   &trueVal,
		OptionsPassthrough: &trueVal,
		Debug:              &trueVal}
	dest := orig
	var zero CORS
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "CORS.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := CORS{
		Enabled:            &falseVal,
		MaxAge:             42,
		AllowedOrigins:     []string{"b", "b"},
		AllowedMethods:     []string{"b", "b"},
		AllowedHeaders:     []string{"b", "b"},
		ExposedHeaders:     []string{"b", "b"},
		AllowCredentials:   &falseVal,
		OptionsPassthrough: &falseVal,
		Debug:              &falseVal}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "CORS.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := CORS{
		Enabled: &trueVal}
	dest.overrideFrom(&o2)
	exp := o

	exp.Enabled = o2.Enabled
	require.Equal(t, dest, exp, "CORS.overrideFrom should have overriden the field Enabled. value now %#v, expecting %#v", dest, exp)
}

func TestCORS_Getters(t *testing.T) {
	orig := CORS{
		Enabled:            &trueVal,
		MaxAge:             -42,
		AllowedOrigins:     []string{"a"},
		AllowedMethods:     []string{"a"},
		AllowedHeaders:     []string{"a"},
		ExposedHeaders:     []string{"a"},
		AllowCredentials:   &trueVal,
		OptionsPassthrough: &trueVal,
		Debug:              &trueVal}

	gv0 := orig.GetEnabled()
	require.Equal(t, orig.Enabled, &gv0, "CORS.GetEnabled() does not match")

	gv1 := orig.GetMaxAge()
	require.Equal(t, orig.MaxAge, gv1, "CORS.GetMaxAgeCfg() does not match")

	gv2 := orig.GetAllowedOrigins()
	require.Equal(t, orig.AllowedOrigins, gv2, "CORS.GetAllowedOriginsCfg() does not match")

	gv3 := orig.GetAllowedMethods()
	require.Equal(t, orig.AllowedMethods, gv3, "CORS.GetAllowedMethodsCfg() does not match")

	gv4 := orig.GetAllowedHeaders()
	require.Equal(t, orig.AllowedHeaders, gv4, "CORS.GetAllowedHeadersCfg() does not match")

	gv5 := orig.GetExposedHeaders()
	require.Equal(t, orig.ExposedHeaders, gv5, "CORS.GetExposedHeadersCfg() does not match")

	gv6 := orig.GetAllowCredentials()
	require.Equal(t, orig.AllowCredentials, &gv6, "CORS.GetAllowCredentials() does not match")

	gv7 := orig.GetOptionsPassthrough()
	require.Equal(t, orig.OptionsPassthrough, &gv7, "CORS.GetOptionsPassthrough() does not match")

	gv8 := orig.GetDebug()
	require.Equal(t, orig.Debug, &gv8, "CORS.GetDebug() does not match")

}

func TestConfiguration_overrideFrom(t *testing.T) {
	orig := Configuration{
		Datacenter:  "one",
		Environment: "one",
		ServiceName: "one",
		HTTP: HTTPServer{
			ServiceName: "one",
			Disabled:    &trueVal,
			VIPName:     "one",
			BindAddr:    "one",
			ServerTLS: TLSInfo{
				CertFile:       "one",
				KeyFile:        "one",
				TrustedCAFile:  "one",
				ClientCertAuth: "one"},
			PackageLogger:  "one",
			AllowProfiling: &trueVal,
			ProfilerDir:    "one",
			Services:       []string{"a"},
			HeartbeatSecs:  -42,
			CORS: CORS{
				Enabled:            &trueVal,
				MaxAge:             -42,
				AllowedOrigins:     []string{"a"},
				AllowedMethods:     []string{"a"},
				AllowedHeaders:     []string{"a"},
				ExposedHeaders:     []string{"a"},
				AllowCredentials:   &trueVal,
				OptionsPassthrough: &trueVal,
				Debug:              &trueVal}},
		HTTPS: HTTPServer{
			ServiceName: "one",
			Disabled:    &trueVal,
			VIPName:     "one",
			BindAddr:    "one",
			ServerTLS: TLSInfo{
				CertFile:       "one",
				KeyFile:        "one",
				TrustedCAFile:  "one",
				ClientCertAuth: "one"},
			PackageLogger:  "one",
			AllowProfiling: &trueVal,
			ProfilerDir:    "one",
			Services:       []string{"a"},
			HeartbeatSecs:  -42,
			CORS: CORS{
				Enabled:            &trueVal,
				MaxAge:             -42,
				AllowedOrigins:     []string{"a"},
				AllowedMethods:     []string{"a"},
				AllowedHeaders:     []string{"a"},
				ExposedHeaders:     []string{"a"},
				AllowCredentials:   &trueVal,
				OptionsPassthrough: &trueVal,
				Debug:              &trueVal}},
		Authz: Authz{
			Allow:        []string{"a"},
			AllowAny:     []string{"a"},
			AllowAnyRole: []string{"a"},
			LogAllowed:   &trueVal,
			LogDenied:    &trueVal,
			CertMapper:   "one",
			APIKeyMapper: "one",
			JWTMapper:    "one"},
		Audit: Logger{
			Directory:  "one",
			MaxAgeDays: -42,
			MaxSizeMb:  -42},
		Metrics: Metrics{
			Provider: "one"},
		Logger: Logger{
			Directory:  "one",
			MaxAgeDays: -42,
			MaxSizeMb:  -42},
		LogLevels: []RepoLogLevel{
			{
				Repo:    "one",
				Package: "one",
				Level:   "one"},
		},
		RootCA: "one",
		Jokes: Jokes{
			NameSvcURL:  "one",
			Categories:  []string{"a"},
			JokesSvcURL: "one"}}
	dest := orig
	var zero Configuration
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "Configuration.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := Configuration{
		Datacenter:  "two",
		Environment: "two",
		ServiceName: "two",
		HTTP: HTTPServer{
			ServiceName: "two",
			Disabled:    &falseVal,
			VIPName:     "two",
			BindAddr:    "two",
			ServerTLS: TLSInfo{
				CertFile:       "two",
				KeyFile:        "two",
				TrustedCAFile:  "two",
				ClientCertAuth: "two"},
			PackageLogger:  "two",
			AllowProfiling: &falseVal,
			ProfilerDir:    "two",
			Services:       []string{"b", "b"},
			HeartbeatSecs:  42,
			CORS: CORS{
				Enabled:            &falseVal,
				MaxAge:             42,
				AllowedOrigins:     []string{"b", "b"},
				AllowedMethods:     []string{"b", "b"},
				AllowedHeaders:     []string{"b", "b"},
				ExposedHeaders:     []string{"b", "b"},
				AllowCredentials:   &falseVal,
				OptionsPassthrough: &falseVal,
				Debug:              &falseVal}},
		HTTPS: HTTPServer{
			ServiceName: "two",
			Disabled:    &falseVal,
			VIPName:     "two",
			BindAddr:    "two",
			ServerTLS: TLSInfo{
				CertFile:       "two",
				KeyFile:        "two",
				TrustedCAFile:  "two",
				ClientCertAuth: "two"},
			PackageLogger:  "two",
			AllowProfiling: &falseVal,
			ProfilerDir:    "two",
			Services:       []string{"b", "b"},
			HeartbeatSecs:  42,
			CORS: CORS{
				Enabled:            &falseVal,
				MaxAge:             42,
				AllowedOrigins:     []string{"b", "b"},
				AllowedMethods:     []string{"b", "b"},
				AllowedHeaders:     []string{"b", "b"},
				ExposedHeaders:     []string{"b", "b"},
				AllowCredentials:   &falseVal,
				OptionsPassthrough: &falseVal,
				Debug:              &falseVal}},
		Authz: Authz{
			Allow:        []string{"b", "b"},
			AllowAny:     []string{"b", "b"},
			AllowAnyRole: []string{"b", "b"},
			LogAllowed:   &falseVal,
			LogDenied:    &falseVal,
			CertMapper:   "two",
			APIKeyMapper: "two",
			JWTMapper:    "two"},
		Audit: Logger{
			Directory:  "two",
			MaxAgeDays: 42,
			MaxSizeMb:  42},
		Metrics: Metrics{
			Provider: "two"},
		Logger: Logger{
			Directory:  "two",
			MaxAgeDays: 42,
			MaxSizeMb:  42},
		LogLevels: []RepoLogLevel{
			{
				Repo:    "two",
				Package: "two",
				Level:   "two"},
		},
		RootCA: "two",
		Jokes: Jokes{
			NameSvcURL:  "two",
			Categories:  []string{"b", "b"},
			JokesSvcURL: "two"}}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "Configuration.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := Configuration{
		Datacenter: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.Datacenter = o2.Datacenter
	require.Equal(t, dest, exp, "Configuration.overrideFrom should have overriden the field Datacenter. value now %#v, expecting %#v", dest, exp)
}

func TestHTTPServer_overrideFrom(t *testing.T) {
	orig := HTTPServer{
		ServiceName: "one",
		Disabled:    &trueVal,
		VIPName:     "one",
		BindAddr:    "one",
		ServerTLS: TLSInfo{
			CertFile:       "one",
			KeyFile:        "one",
			TrustedCAFile:  "one",
			ClientCertAuth: "one"},
		PackageLogger:  "one",
		AllowProfiling: &trueVal,
		ProfilerDir:    "one",
		Services:       []string{"a"},
		HeartbeatSecs:  -42,
		CORS: CORS{
			Enabled:            &trueVal,
			MaxAge:             -42,
			AllowedOrigins:     []string{"a"},
			AllowedMethods:     []string{"a"},
			AllowedHeaders:     []string{"a"},
			ExposedHeaders:     []string{"a"},
			AllowCredentials:   &trueVal,
			OptionsPassthrough: &trueVal,
			Debug:              &trueVal}}
	dest := orig
	var zero HTTPServer
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "HTTPServer.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := HTTPServer{
		ServiceName: "two",
		Disabled:    &falseVal,
		VIPName:     "two",
		BindAddr:    "two",
		ServerTLS: TLSInfo{
			CertFile:       "two",
			KeyFile:        "two",
			TrustedCAFile:  "two",
			ClientCertAuth: "two"},
		PackageLogger:  "two",
		AllowProfiling: &falseVal,
		ProfilerDir:    "two",
		Services:       []string{"b", "b"},
		HeartbeatSecs:  42,
		CORS: CORS{
			Enabled:            &falseVal,
			MaxAge:             42,
			AllowedOrigins:     []string{"b", "b"},
			AllowedMethods:     []string{"b", "b"},
			AllowedHeaders:     []string{"b", "b"},
			ExposedHeaders:     []string{"b", "b"},
			AllowCredentials:   &falseVal,
			OptionsPassthrough: &falseVal,
			Debug:              &falseVal}}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "HTTPServer.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := HTTPServer{
		ServiceName: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.ServiceName = o2.ServiceName
	require.Equal(t, dest, exp, "HTTPServer.overrideFrom should have overriden the field ServiceName. value now %#v, expecting %#v", dest, exp)
}

func TestHTTPServer_Getters(t *testing.T) {
	orig := HTTPServer{
		ServiceName: "one",
		Disabled:    &trueVal,
		VIPName:     "one",
		BindAddr:    "one",
		ServerTLS: TLSInfo{
			CertFile:       "one",
			KeyFile:        "one",
			TrustedCAFile:  "one",
			ClientCertAuth: "one"},
		PackageLogger:  "one",
		AllowProfiling: &trueVal,
		ProfilerDir:    "one",
		Services:       []string{"a"},
		HeartbeatSecs:  -42,
		CORS: CORS{
			Enabled:            &trueVal,
			MaxAge:             -42,
			AllowedOrigins:     []string{"a"},
			AllowedMethods:     []string{"a"},
			AllowedHeaders:     []string{"a"},
			ExposedHeaders:     []string{"a"},
			AllowCredentials:   &trueVal,
			OptionsPassthrough: &trueVal,
			Debug:              &trueVal}}

	gv0 := orig.GetServiceName()
	require.Equal(t, orig.ServiceName, gv0, "HTTPServer.GetServiceNameCfg() does not match")

	gv1 := orig.GetDisabled()
	require.Equal(t, orig.Disabled, &gv1, "HTTPServer.GetDisabled() does not match")

	gv2 := orig.GetVIPName()
	require.Equal(t, orig.VIPName, gv2, "HTTPServer.GetVIPNameCfg() does not match")

	gv3 := orig.GetBindAddr()
	require.Equal(t, orig.BindAddr, gv3, "HTTPServer.GetBindAddrCfg() does not match")

	gv4 := orig.GetServerTLSCfg()
	require.Equal(t, orig.ServerTLS, *gv4, "HTTPServer.GetServerTLSCfg() does not match")

	gv5 := orig.GetPackageLogger()
	require.Equal(t, orig.PackageLogger, gv5, "HTTPServer.GetPackageLoggerCfg() does not match")

	gv6 := orig.GetAllowProfiling()
	require.Equal(t, orig.AllowProfiling, &gv6, "HTTPServer.GetAllowProfiling() does not match")

	gv7 := orig.GetProfilerDir()
	require.Equal(t, orig.ProfilerDir, gv7, "HTTPServer.GetProfilerDirCfg() does not match")

	gv8 := orig.GetServices()
	require.Equal(t, orig.Services, gv8, "HTTPServer.GetServicesCfg() does not match")

	gv9 := orig.GetHeartbeatSecs()
	require.Equal(t, orig.HeartbeatSecs, gv9, "HTTPServer.GetHeartbeatSecsCfg() does not match")

	gv10 := orig.GetCORSCfg()
	require.Equal(t, orig.CORS, *gv10, "HTTPServer.GetCORSCfg() does not match")

}

func TestJokes_overrideFrom(t *testing.T) {
	orig := Jokes{
		NameSvcURL:  "one",
		Categories:  []string{"a"},
		JokesSvcURL: "one"}
	dest := orig
	var zero Jokes
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "Jokes.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := Jokes{
		NameSvcURL:  "two",
		Categories:  []string{"b", "b"},
		JokesSvcURL: "two"}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "Jokes.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := Jokes{
		NameSvcURL: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.NameSvcURL = o2.NameSvcURL
	require.Equal(t, dest, exp, "Jokes.overrideFrom should have overriden the field NameSvcURL. value now %#v, expecting %#v", dest, exp)
}

func TestLogger_overrideFrom(t *testing.T) {
	orig := Logger{
		Directory:  "one",
		MaxAgeDays: -42,
		MaxSizeMb:  -42}
	dest := orig
	var zero Logger
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "Logger.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := Logger{
		Directory:  "two",
		MaxAgeDays: 42,
		MaxSizeMb:  42}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "Logger.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := Logger{
		Directory: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.Directory = o2.Directory
	require.Equal(t, dest, exp, "Logger.overrideFrom should have overriden the field Directory. value now %#v, expecting %#v", dest, exp)
}

func TestLogger_Getters(t *testing.T) {
	orig := Logger{
		Directory:  "one",
		MaxAgeDays: -42,
		MaxSizeMb:  -42}

	gv0 := orig.GetDirectory()
	require.Equal(t, orig.Directory, gv0, "Logger.GetDirectoryCfg() does not match")

	gv1 := orig.GetMaxAgeDays()
	require.Equal(t, orig.MaxAgeDays, gv1, "Logger.GetMaxAgeDaysCfg() does not match")

	gv2 := orig.GetMaxSizeMb()
	require.Equal(t, orig.MaxSizeMb, gv2, "Logger.GetMaxSizeMbCfg() does not match")

}

func TestMetrics_overrideFrom(t *testing.T) {
	orig := Metrics{
		Provider: "one"}
	dest := orig
	var zero Metrics
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "Metrics.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := Metrics{
		Provider: "two"}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "Metrics.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := Metrics{
		Provider: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.Provider = o2.Provider
	require.Equal(t, dest, exp, "Metrics.overrideFrom should have overriden the field Provider. value now %#v, expecting %#v", dest, exp)
}

func TestRepoLogLevel_overrideFrom(t *testing.T) {
	orig := RepoLogLevel{
		Repo:    "one",
		Package: "one",
		Level:   "one"}
	dest := orig
	var zero RepoLogLevel
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "RepoLogLevel.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := RepoLogLevel{
		Repo:    "two",
		Package: "two",
		Level:   "two"}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "RepoLogLevel.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := RepoLogLevel{
		Repo: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.Repo = o2.Repo
	require.Equal(t, dest, exp, "RepoLogLevel.overrideFrom should have overriden the field Repo. value now %#v, expecting %#v", dest, exp)
}

func TestTLSInfo_overrideFrom(t *testing.T) {
	orig := TLSInfo{
		CertFile:       "one",
		KeyFile:        "one",
		TrustedCAFile:  "one",
		ClientCertAuth: "one"}
	dest := orig
	var zero TLSInfo
	dest.overrideFrom(&zero)
	require.Equal(t, dest, orig, "TLSInfo.overrideFrom shouldn't have overriden the value as the override is the default/zero value. value now %#v", dest)
	o := TLSInfo{
		CertFile:       "two",
		KeyFile:        "two",
		TrustedCAFile:  "two",
		ClientCertAuth: "two"}
	dest.overrideFrom(&o)
	require.Equal(t, dest, o, "TLSInfo.overrideFrom should have overriden the value as the override. value now %#v, expecting %#v", dest, o)
	o2 := TLSInfo{
		CertFile: "one"}
	dest.overrideFrom(&o2)
	exp := o

	exp.CertFile = o2.CertFile
	require.Equal(t, dest, exp, "TLSInfo.overrideFrom should have overriden the field CertFile. value now %#v, expecting %#v", dest, exp)
}

func TestTLSInfo_Getters(t *testing.T) {
	orig := TLSInfo{
		CertFile:       "one",
		KeyFile:        "one",
		TrustedCAFile:  "one",
		ClientCertAuth: "one"}

	gv0 := orig.GetCertFile()
	require.Equal(t, orig.CertFile, gv0, "TLSInfo.GetCertFileCfg() does not match")

	gv1 := orig.GetKeyFile()
	require.Equal(t, orig.KeyFile, gv1, "TLSInfo.GetKeyFileCfg() does not match")

	gv2 := orig.GetTrustedCAFile()
	require.Equal(t, orig.TrustedCAFile, gv2, "TLSInfo.GetTrustedCAFileCfg() does not match")

	gv3 := orig.GetClientCertAuth()
	require.Equal(t, orig.ClientCertAuth, gv3, "TLSInfo.GetClientCertAuthCfg() does not match")

}

func Test_LoadOverrides(t *testing.T) {

	c := Configurations{
		Defaults: Configuration{
			Datacenter:  "two",
			Environment: "two",
			ServiceName: "two",
			HTTP: HTTPServer{
				ServiceName: "two",
				Disabled:    &falseVal,
				VIPName:     "two",
				BindAddr:    "two",
				ServerTLS: TLSInfo{
					CertFile:       "two",
					KeyFile:        "two",
					TrustedCAFile:  "two",
					ClientCertAuth: "two"},
				PackageLogger:  "two",
				AllowProfiling: &falseVal,
				ProfilerDir:    "two",
				Services:       []string{"b", "b"},
				HeartbeatSecs:  42,
				CORS: CORS{
					Enabled:            &falseVal,
					MaxAge:             42,
					AllowedOrigins:     []string{"b", "b"},
					AllowedMethods:     []string{"b", "b"},
					AllowedHeaders:     []string{"b", "b"},
					ExposedHeaders:     []string{"b", "b"},
					AllowCredentials:   &falseVal,
					OptionsPassthrough: &falseVal,
					Debug:              &falseVal}},
			HTTPS: HTTPServer{
				ServiceName: "two",
				Disabled:    &falseVal,
				VIPName:     "two",
				BindAddr:    "two",
				ServerTLS: TLSInfo{
					CertFile:       "two",
					KeyFile:        "two",
					TrustedCAFile:  "two",
					ClientCertAuth: "two"},
				PackageLogger:  "two",
				AllowProfiling: &falseVal,
				ProfilerDir:    "two",
				Services:       []string{"b", "b"},
				HeartbeatSecs:  42,
				CORS: CORS{
					Enabled:            &falseVal,
					MaxAge:             42,
					AllowedOrigins:     []string{"b", "b"},
					AllowedMethods:     []string{"b", "b"},
					AllowedHeaders:     []string{"b", "b"},
					ExposedHeaders:     []string{"b", "b"},
					AllowCredentials:   &falseVal,
					OptionsPassthrough: &falseVal,
					Debug:              &falseVal}},
			Authz: Authz{
				Allow:        []string{"b", "b"},
				AllowAny:     []string{"b", "b"},
				AllowAnyRole: []string{"b", "b"},
				LogAllowed:   &falseVal,
				LogDenied:    &falseVal,
				CertMapper:   "two",
				APIKeyMapper: "two",
				JWTMapper:    "two"},
			Audit: Logger{
				Directory:  "two",
				MaxAgeDays: 42,
				MaxSizeMb:  42},
			Metrics: Metrics{
				Provider: "two"},
			Logger: Logger{
				Directory:  "two",
				MaxAgeDays: 42,
				MaxSizeMb:  42},
			LogLevels: []RepoLogLevel{
				{
					Repo:    "two",
					Package: "two",
					Level:   "two"},
			},
			RootCA: "two",
			Jokes: Jokes{
				NameSvcURL:  "two",
				Categories:  []string{"b", "b"},
				JokesSvcURL: "two"}},
		Hosts: map[string]string{"bob": "example2", "bob2": "missing"},
		Overrides: map[string]Configuration{
			"example2": {
				Datacenter:  "three",
				Environment: "three",
				ServiceName: "three",
				HTTP: HTTPServer{
					ServiceName: "three",
					Disabled:    &trueVal,
					VIPName:     "three",
					BindAddr:    "three",
					ServerTLS: TLSInfo{
						CertFile:       "three",
						KeyFile:        "three",
						TrustedCAFile:  "three",
						ClientCertAuth: "three"},
					PackageLogger:  "three",
					AllowProfiling: &trueVal,
					ProfilerDir:    "three",
					Services:       []string{"c", "c", "c"},
					HeartbeatSecs:  1234,
					CORS: CORS{
						Enabled:            &trueVal,
						MaxAge:             1234,
						AllowedOrigins:     []string{"c", "c", "c"},
						AllowedMethods:     []string{"c", "c", "c"},
						AllowedHeaders:     []string{"c", "c", "c"},
						ExposedHeaders:     []string{"c", "c", "c"},
						AllowCredentials:   &trueVal,
						OptionsPassthrough: &trueVal,
						Debug:              &trueVal}},
				HTTPS: HTTPServer{
					ServiceName: "three",
					Disabled:    &trueVal,
					VIPName:     "three",
					BindAddr:    "three",
					ServerTLS: TLSInfo{
						CertFile:       "three",
						KeyFile:        "three",
						TrustedCAFile:  "three",
						ClientCertAuth: "three"},
					PackageLogger:  "three",
					AllowProfiling: &trueVal,
					ProfilerDir:    "three",
					Services:       []string{"c", "c", "c"},
					HeartbeatSecs:  1234,
					CORS: CORS{
						Enabled:            &trueVal,
						MaxAge:             1234,
						AllowedOrigins:     []string{"c", "c", "c"},
						AllowedMethods:     []string{"c", "c", "c"},
						AllowedHeaders:     []string{"c", "c", "c"},
						ExposedHeaders:     []string{"c", "c", "c"},
						AllowCredentials:   &trueVal,
						OptionsPassthrough: &trueVal,
						Debug:              &trueVal}},
				Authz: Authz{
					Allow:        []string{"c", "c", "c"},
					AllowAny:     []string{"c", "c", "c"},
					AllowAnyRole: []string{"c", "c", "c"},
					LogAllowed:   &trueVal,
					LogDenied:    &trueVal,
					CertMapper:   "three",
					APIKeyMapper: "three",
					JWTMapper:    "three"},
				Audit: Logger{
					Directory:  "three",
					MaxAgeDays: 1234,
					MaxSizeMb:  1234},
				Metrics: Metrics{
					Provider: "three"},
				Logger: Logger{
					Directory:  "three",
					MaxAgeDays: 1234,
					MaxSizeMb:  1234},
				LogLevels: []RepoLogLevel{
					{
						Repo:    "three",
						Package: "three",
						Level:   "three"},
				},
				RootCA: "three",
				Jokes: Jokes{
					NameSvcURL:  "three",
					Categories:  []string{"c", "c", "c"},
					JokesSvcURL: "three"}},
		},
	}
	f, err := ioutil.TempFile("", "config")
	if err != nil {
		t.Fatalf("Uanble to create temp file: %v", err)
	}
	json.NewEncoder(f).Encode(&c)
	f.Close()
	defer os.Remove(f.Name())
	config, err := Load(f.Name(), "", "")
	if err != nil {
		t.Fatalf("Unexpected error loading config: %v", err)
	}
	require.Equal(t, c.Defaults, *config, "Loaded configuration should match default, but doesn't, expecting %#v, got %#v", c.Defaults, *config)
	config, err = Load(f.Name(), "", "bob")
	if err != nil {
		t.Fatalf("Unexpected error loading config: %v", err)
	}
	require.Equal(t, c.Overrides["example2"], *config, "Loaded configuration should match default, but doesn't, expecting %#v, got %#v", c.Overrides["example2"], *config)
	_, err = Load(f.Name(), "", "bob2")
	if err == nil || err.Error() != "Configuration for host bob2 specified override set missing but that doesn't exist" {
		t.Errorf("Should of gotten error about missing override set, but got %v", err)
	}
}

func Test_LoadMissingFile(t *testing.T) {
	f, err := ioutil.TempFile("", "missing")
	f.Close()
	os.Remove(f.Name())
	_, err = Load(f.Name(), "", "")
	if !os.IsNotExist(err) {
		t.Errorf("Expecting a file doesn't exist error when trying to load from a non-existant file, but got %v", err)
	}
}

func Test_LoadInvalidJson(t *testing.T) {
	f, err := ioutil.TempFile("", "invalid")
	f.WriteString("{boom}")
	f.Close()
	defer os.Remove(f.Name())
	_, err = Load(f.Name(), "", "")
	if err == nil || err.Error() != "invalid character 'b' looking for beginning of object key string" {
		t.Errorf("Should get a json error with an invalid config file, but got %v", err)
	}
}

func loadJSONEWithENV(filename string, v interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	val := strings.ReplaceAll(string(bytes), "${ENV}", "ENV_VALUE")
	return json.NewDecoder(strings.NewReader(val)).Decode(v)
}

func Test_LoadCustomJSON(t *testing.T) {

	c := Configurations{
		Defaults: Configuration{
			Datacenter:  "two",
			Environment: "two",
			ServiceName: "two",
			HTTP: HTTPServer{
				ServiceName: "two",
				Disabled:    &falseVal,
				VIPName:     "two",
				BindAddr:    "two",
				ServerTLS: TLSInfo{
					CertFile:       "two",
					KeyFile:        "two",
					TrustedCAFile:  "two",
					ClientCertAuth: "two"},
				PackageLogger:  "two",
				AllowProfiling: &falseVal,
				ProfilerDir:    "two",
				Services:       []string{"b", "b"},
				HeartbeatSecs:  42,
				CORS: CORS{
					Enabled:            &falseVal,
					MaxAge:             42,
					AllowedOrigins:     []string{"b", "b"},
					AllowedMethods:     []string{"b", "b"},
					AllowedHeaders:     []string{"b", "b"},
					ExposedHeaders:     []string{"b", "b"},
					AllowCredentials:   &falseVal,
					OptionsPassthrough: &falseVal,
					Debug:              &falseVal}},
			HTTPS: HTTPServer{
				ServiceName: "two",
				Disabled:    &falseVal,
				VIPName:     "two",
				BindAddr:    "two",
				ServerTLS: TLSInfo{
					CertFile:       "two",
					KeyFile:        "two",
					TrustedCAFile:  "two",
					ClientCertAuth: "two"},
				PackageLogger:  "two",
				AllowProfiling: &falseVal,
				ProfilerDir:    "two",
				Services:       []string{"b", "b"},
				HeartbeatSecs:  42,
				CORS: CORS{
					Enabled:            &falseVal,
					MaxAge:             42,
					AllowedOrigins:     []string{"b", "b"},
					AllowedMethods:     []string{"b", "b"},
					AllowedHeaders:     []string{"b", "b"},
					ExposedHeaders:     []string{"b", "b"},
					AllowCredentials:   &falseVal,
					OptionsPassthrough: &falseVal,
					Debug:              &falseVal}},
			Authz: Authz{
				Allow:        []string{"b", "b"},
				AllowAny:     []string{"b", "b"},
				AllowAnyRole: []string{"b", "b"},
				LogAllowed:   &falseVal,
				LogDenied:    &falseVal,
				CertMapper:   "two",
				APIKeyMapper: "two",
				JWTMapper:    "two"},
			Audit: Logger{
				Directory:  "two",
				MaxAgeDays: 42,
				MaxSizeMb:  42},
			Metrics: Metrics{
				Provider: "two"},
			Logger: Logger{
				Directory:  "two",
				MaxAgeDays: 42,
				MaxSizeMb:  42},
			LogLevels: []RepoLogLevel{
				{
					Repo:    "two",
					Package: "two",
					Level:   "two"},
			},
			RootCA: "two",
			Jokes: Jokes{
				NameSvcURL:  "two",
				Categories:  []string{"b", "b"},
				JokesSvcURL: "two"}},
		Hosts: map[string]string{"bob": "${ENV}"},
		Overrides: map[string]Configuration{
			"${ENV}": {
				Datacenter:  "three",
				Environment: "three",
				ServiceName: "three",
				HTTP: HTTPServer{
					ServiceName: "three",
					Disabled:    &trueVal,
					VIPName:     "three",
					BindAddr:    "three",
					ServerTLS: TLSInfo{
						CertFile:       "three",
						KeyFile:        "three",
						TrustedCAFile:  "three",
						ClientCertAuth: "three"},
					PackageLogger:  "three",
					AllowProfiling: &trueVal,
					ProfilerDir:    "three",
					Services:       []string{"c", "c", "c"},
					HeartbeatSecs:  1234,
					CORS: CORS{
						Enabled:            &trueVal,
						MaxAge:             1234,
						AllowedOrigins:     []string{"c", "c", "c"},
						AllowedMethods:     []string{"c", "c", "c"},
						AllowedHeaders:     []string{"c", "c", "c"},
						ExposedHeaders:     []string{"c", "c", "c"},
						AllowCredentials:   &trueVal,
						OptionsPassthrough: &trueVal,
						Debug:              &trueVal}},
				HTTPS: HTTPServer{
					ServiceName: "three",
					Disabled:    &trueVal,
					VIPName:     "three",
					BindAddr:    "three",
					ServerTLS: TLSInfo{
						CertFile:       "three",
						KeyFile:        "three",
						TrustedCAFile:  "three",
						ClientCertAuth: "three"},
					PackageLogger:  "three",
					AllowProfiling: &trueVal,
					ProfilerDir:    "three",
					Services:       []string{"c", "c", "c"},
					HeartbeatSecs:  1234,
					CORS: CORS{
						Enabled:            &trueVal,
						MaxAge:             1234,
						AllowedOrigins:     []string{"c", "c", "c"},
						AllowedMethods:     []string{"c", "c", "c"},
						AllowedHeaders:     []string{"c", "c", "c"},
						ExposedHeaders:     []string{"c", "c", "c"},
						AllowCredentials:   &trueVal,
						OptionsPassthrough: &trueVal,
						Debug:              &trueVal}},
				Authz: Authz{
					Allow:        []string{"c", "c", "c"},
					AllowAny:     []string{"c", "c", "c"},
					AllowAnyRole: []string{"c", "c", "c"},
					LogAllowed:   &trueVal,
					LogDenied:    &trueVal,
					CertMapper:   "three",
					APIKeyMapper: "three",
					JWTMapper:    "three"},
				Audit: Logger{
					Directory:  "three",
					MaxAgeDays: 1234,
					MaxSizeMb:  1234},
				Metrics: Metrics{
					Provider: "three"},
				Logger: Logger{
					Directory:  "three",
					MaxAgeDays: 1234,
					MaxSizeMb:  1234},
				LogLevels: []RepoLogLevel{
					{
						Repo:    "three",
						Package: "three",
						Level:   "three"},
				},
				RootCA: "three",
				Jokes: Jokes{
					NameSvcURL:  "three",
					Categories:  []string{"c", "c", "c"},
					JokesSvcURL: "three"}},
		},
	}
	f, err := ioutil.TempFile("", "customjson")
	if err != nil {
		t.Fatalf("Uanble to create temp file: %v", err)
	}
	json.NewEncoder(f).Encode(&c)
	f.Close()
	defer os.Remove(f.Name())

	JSONLoader = loadJSONEWithENV
	config, err := Load(f.Name(), "", "")
	if err != nil {
		t.Fatalf("Unexpected error loading config: %v", err)
	}
	require.Equal(t, c.Defaults, *config, "Loaded configuration should match default, but doesn't, expecting %#v, got %#v", c.Defaults, *config)

	JSONLoader = loadJSONEWithENV
	config, err = Load(f.Name(), "", "bob")
	if err != nil {
		t.Fatalf("Unexpected error loading config: %v", err)
	}
	require.Equal(t, c.Overrides["${ENV}"], *config,
		"Loaded configuration should match default, but doesn't, expecting %#v, got %#v\nOverrides: %v", c.Overrides["${ENV}"], *config, c.Overrides)
}
