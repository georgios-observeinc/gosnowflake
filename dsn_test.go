// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	cr "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

type tcParseDSN struct {
	dsn      string
	config   *Config
	ocspMode string
	err      error
}

func TestParseDSN(t *testing.T) {
	privKeyPKCS8 := generatePKCS8StringSupress(testPrivKey)
	privKeyPKCS1 := generatePKCS1String(testPrivKey)
	testcases := []tcParseDSN{
		{
			dsn: "user:pass@ac-1-laksdnflaf.global/db/schema",
			config: &Config{
				Account: "ac-1", User: "user", Password: "pass", Region: "global",
				Protocol: "https", Host: "ac-1-laksdnflaf.global.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@ac-laksdnflaf.global/db/schema",
			config: &Config{
				Account: "ac", User: "user", Password: "pass", Region: "global",
				Protocol: "https", Host: "ac-laksdnflaf.global.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:p@asnowflakecomputing.com/db/pa?account=a&protocol=https&role=r&timezone=UTC&aehouse=w",
			config: &Config{Account: "a", User: "u", Password: "p", Database: "db", Schema: "pa",
				Protocol: "https", Role: "r", Host: "asnowflakecomputing.com.snowflakecomputing.com", Port: 443, Region: "com",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:p@/db?account=ac",
			config: &Config{
				Account: "ac", User: "u", Password: "p", Database: "db",
				Protocol: "https", Host: "ac.snowflakecomputing.com", Port: 443,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account-hfdw89q748ew9gqf48w9qgf.global/db/s",
			config: &Config{
				Account: "account", User: "user", Password: "pass", Region: "global",
				Protocol: "https", Host: "account-hfdw89q748ew9gqf48w9qgf.global.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s",
				ValidateDefaultParameters: ConfigBoolTrue,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account-hfdw89q748ew9gqf48w9qgf/db/s",
			config: &Config{
				Account: "account-hfdw89q748ew9gqf48w9qgf", User: "user", Password: "pass", Region: "",
				Protocol: "https", Host: "account-hfdw89q748ew9gqf48w9qgf.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s",
				ValidateDefaultParameters: ConfigBoolTrue,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account",
			config: &Config{
				Account: "account", User: "user", Password: "pass", Region: "",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account.eu-faraway",
			config: &Config{
				Account: "account", User: "user", Password: "pass", Region: "eu-faraway",
				Protocol: "https", Host: "account.eu-faraway.snowflakecomputing.com", Port: 443,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account?region=eu-faraway",
			config: &Config{
				Account: "account", User: "user", Password: "pass", Region: "eu-faraway",
				Protocol: "https", Host: "account.eu-faraway.snowflakecomputing.com", Port: 443,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account/db",
			config: &Config{
				Account: "account", User: "user", Password: "pass",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				Database:                  "db",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@host:123/db/schema?account=ac&protocol=http",
			config: &Config{
				Account: "ac", User: "user", Password: "pass",
				Protocol: "http", Host: "host", Port: 123,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user@host:123/db/schema?account=ac&protocol=http",
			config: &Config{
				Account: "ac", User: "user", Password: "pass",
				Protocol: "http", Host: "host", Port: 123,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      errEmptyPassword(),
		},
		{
			dsn: "@host:123/db/schema?account=ac&protocol=http",
			config: &Config{
				Account: "ac", User: "user", Password: "pass",
				Protocol: "http", Host: "host", Port: 123,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      errEmptyUsername(),
		},
		{
			dsn: "user:p@host:123/db/schema?protocol=http",
			config: &Config{
				Account: "ac", User: "user", Password: "pass",
				Protocol: "http", Host: "host", Port: 123,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      errEmptyAccount(),
		},
		{
			dsn: "u:p@a.snowflakecomputing.com/db/pa?account=a&protocol=https&role=r&timezone=UTC&warehouse=w",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "pa", Role: "r", Warehouse: "w",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:p@snowflake.local:9876?account=a&protocol=http",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "http", Host: "snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "snowflake.local:9876?account=a&protocol=http&authenticator=OAUTH",
			config: &Config{
				Account: "a", Authenticator: AuthTypeOAuth,
				Protocol: "http", Host: "snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:@a.snowflake.local:9876?account=a&protocol=http&authenticator=SNOWFLAKE_JWT",
			config: &Config{
				Account: "a", User: "u", Authenticator: AuthTypeJwt,
				Protocol: "http", Host: "a.snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},

		{
			dsn: "u:p@a?database=d&jwtTimeout=20",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.snowflakecomputing.com", Port: 443,
				Database: "d", Schema: "",
				JWTExpireTimeout:          20 * time.Second,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
		},
		{
			dsn: "u:p@a?database=d&externalBrowserTimeout=20",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.snowflakecomputing.com", Port: 443,
				Database: "d", Schema: "",
				ExternalBrowserTimeout:    20 * time.Second,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
			},
			ocspMode: ocspModeFailOpen,
		},
		{
			dsn: "u:p@a?database=d",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.snowflakecomputing.com", Port: 443,
				Database: "d", Schema: "",
				JWTExpireTimeout:          defaultJWTTimeout,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
		},
		{
			dsn: "u:p@snowflake.local:NNNN?account=a&protocol=http",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "http", Host: "snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err: &SnowflakeError{
				Message:     errMsgFailedToParsePort,
				MessageArgs: []interface{}{"NNNN"},
				Number:      ErrCodeFailedToParsePort,
			},
		},
		{
			dsn: "u:p@a?database=d&schema=s&role=r&application=aa&authenticator=snowflake&insecureMode=true&passcode=pp&passcodeInPassword=true",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.snowflakecomputing.com", Port: 443,
				Database: "d", Schema: "s", Role: "r", Authenticator: AuthTypeSnowflake, Application: "aa",
				InsecureMode: true, Passcode: "pp", PasscodeInPassword: true,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeInsecure,
			err:      nil,
		},
		{
			// schema should be ignored as no value is specified.
			dsn: "u:p@a?database=d&schema",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.snowflakecomputing.com", Port: 443,
				Database: "d", Schema: "",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn:    "u:p@a?database= %Sd",
			config: &Config{},
			err:    url.EscapeError(`invalid URL escape`),
		},
		{
			dsn:    "u:p@a?schema= %Sd",
			config: &Config{},
			err:    url.EscapeError(`invalid URL escape`),
		},
		{
			dsn:    "u:p@a?warehouse= %Sd",
			config: &Config{},
			err:    url.EscapeError(`invalid URL escape`),
		},
		{
			dsn:    "u:p@a?role= %Sd",
			config: &Config{},
			err:    url.EscapeError(`invalid URL escape`),
		},
		{
			dsn:    ":/",
			config: &Config{},
			err: &SnowflakeError{
				Number: ErrCodeFailedToParsePort,
			},
		},
		{
			dsn:    "u:u@/+/+?account=+&=0",
			config: &Config{},
			err:    errEmptyAccount(),
		},
		{
			dsn:    "u:u@/+/+?account=+&=+&=+",
			config: &Config{},
			err:    errEmptyAccount(),
		},
		{
			dsn: "user%40%2F1:p%3A%40s@/db%2F?account=ac",
			config: &Config{
				Account: "ac", User: "user@/1", Password: "p:@s", Database: "db/",
				Protocol: "https", Host: "ac.snowflakecomputing.com", Port: 443,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: fmt.Sprintf("u:p@ac.snowflake.local:9876?account=ac&protocol=http&authenticator=SNOWFLAKE_JWT&privateKey=%v", privKeyPKCS8),
			config: &Config{
				Account: "ac", User: "u", Password: "p",
				Authenticator: AuthTypeJwt, PrivateKey: testPrivKey,
				Protocol: "http", Host: "ac.snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: fmt.Sprintf("u:p@ac.snowflake.local:9876?account=ac&protocol=http&authenticator=%v", url.QueryEscape("https://ac.okta.com")),
			config: &Config{
				Account: "ac", User: "u", Password: "p",
				Authenticator: AuthTypeOkta,
				OktaURL: &url.URL{
					Scheme: "https",
					Host:   "ac.okta.com",
				},
				PrivateKey: testPrivKey,
				Protocol:   "http", Host: "ac.snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: fmt.Sprintf("u:p@a.snowflake.local:9876?account=a&protocol=http&authenticator=SNOWFLAKE_JWT&privateKey=%v", privKeyPKCS1),
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Authenticator: AuthTypeJwt, PrivateKey: testPrivKey,
				Protocol: "http", Host: "a.snowflake.local", Port: 9876,
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      &SnowflakeError{Number: ErrCodePrivateKeyParseError},
		},
		{
			dsn: "user:pass@account/db/s?ocspFailOpen=true",
			config: &Config{
				Account: "account", User: "user", Password: "pass",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", OCSPFailOpen: OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account/db/s?ocspFailOpen=false",
			config: &Config{
				Account: "account", User: "user", Password: "pass",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", OCSPFailOpen: OCSPFailOpenFalse,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailClosed,
			err:      nil,
		},
		{
			dsn: "user:pass@account/db/s?insecureMode=true&ocspFailOpen=false",
			config: &Config{
				Account: "account", User: "user", Password: "pass",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", OCSPFailOpen: OCSPFailOpenFalse, InsecureMode: true,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeInsecure,
			err:      nil,
		},
		{
			dsn: "user:pass@account/db/s?validateDefaultParameters=true",
			config: &Config{
				Account: "account", User: "user", Password: "pass",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", ValidateDefaultParameters: ConfigBoolTrue, OCSPFailOpen: OCSPFailOpenTrue,
				ClientTimeout:          defaultClientTimeout,
				JWTClientTimeout:       defaultJWTClientTimeout,
				ExternalBrowserTimeout: defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "user:pass@account/db/s?validateDefaultParameters=false",
			config: &Config{
				Account: "account", User: "user", Password: "pass",
				Protocol: "https", Host: "account.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", ValidateDefaultParameters: ConfigBoolFalse, OCSPFailOpen: OCSPFailOpenTrue,
				ClientTimeout:          defaultClientTimeout,
				JWTClientTimeout:       defaultJWTClientTimeout,
				ExternalBrowserTimeout: defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:p@a.r.c.snowflakecomputing.com/db/s?account=a.r.c&validateDefaultParameters=false",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.r.c.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", ValidateDefaultParameters: ConfigBoolFalse, OCSPFailOpen: OCSPFailOpenTrue,
				ClientTimeout:          defaultClientTimeout,
				JWTClientTimeout:       defaultJWTClientTimeout,
				ExternalBrowserTimeout: defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:p@a.r.c.snowflakecomputing.com/db/s?account=a.r.c&clientTimeout=300&jwtClientTimeout=45",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.r.c.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", ValidateDefaultParameters: ConfigBoolTrue, OCSPFailOpen: OCSPFailOpenTrue,
				ClientTimeout:          300 * time.Second,
				JWTClientTimeout:       45 * time.Second,
				ExternalBrowserTimeout: defaultExternalBrowserTimeout,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
		{
			dsn: "u:p@a.r.c.snowflakecomputing.com/db/s?account=a.r.c&tmpDirPath=%2Ftmp",
			config: &Config{
				Account: "a", User: "u", Password: "p",
				Protocol: "https", Host: "a.r.c.snowflakecomputing.com", Port: 443,
				Database: "db", Schema: "s", ValidateDefaultParameters: ConfigBoolTrue, OCSPFailOpen: OCSPFailOpenTrue,
				ClientTimeout:          defaultClientTimeout,
				JWTClientTimeout:       defaultJWTClientTimeout,
				ExternalBrowserTimeout: defaultExternalBrowserTimeout,
				TmpDirPath:             "/tmp",
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		},
	}

	for _, at := range []AuthType{AuthTypeExternalBrowser, AuthTypeOAuth} {
		testcases = append(testcases, tcParseDSN{
			dsn: fmt.Sprintf("@host:777/db/schema?account=ac&protocol=http&authenticator=%v", strings.ToLower(at.String())),
			config: &Config{
				Account: "ac", User: "", Password: "",
				Protocol: "http", Host: "host", Port: 777,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
				Authenticator:             at,
			},
			ocspMode: ocspModeFailOpen,
			err:      nil,
		})
	}

	for _, at := range []AuthType{AuthTypeSnowflake, AuthTypeUsernamePasswordMFA, AuthTypeJwt} {
		testcases = append(testcases, tcParseDSN{
			dsn: fmt.Sprintf("@host:888/db/schema?account=ac&protocol=http&authenticator=%v", strings.ToLower(at.String())),
			config: &Config{
				Account: "ac", User: "", Password: "",
				Protocol: "http", Host: "host", Port: 888,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
				Authenticator:             at,
			},
			ocspMode: ocspModeFailOpen,
			err:      errEmptyUsername(),
		})
	}

	for _, at := range []AuthType{AuthTypeSnowflake, AuthTypeUsernamePasswordMFA} {
		testcases = append(testcases, tcParseDSN{
			dsn: fmt.Sprintf("user@host:888/db/schema?account=ac&protocol=http&authenticator=%v", strings.ToLower(at.String())),
			config: &Config{
				Account: "ac", User: "user", Password: "",
				Protocol: "http", Host: "host", Port: 888,
				Database: "db", Schema: "schema",
				OCSPFailOpen:              OCSPFailOpenTrue,
				ValidateDefaultParameters: ConfigBoolTrue,
				ClientTimeout:             defaultClientTimeout,
				JWTClientTimeout:          defaultJWTClientTimeout,
				ExternalBrowserTimeout:    defaultExternalBrowserTimeout,
				Authenticator:             at,
			},
			ocspMode: ocspModeFailOpen,
			err:      errEmptyPassword(),
		})
	}

	for i, test := range testcases {
		t.Run(test.dsn, func(t *testing.T) {
			cfg, err := ParseDSN(test.dsn)
			switch {
			case test.err == nil:
				if err != nil {
					t.Fatalf("%d: Failed to parse the DSN. dsn: %v, err: %v", i, test.dsn, err)
				}
				if test.config.Host != cfg.Host {
					t.Fatalf("%d: Failed to match host. expected: %v, got: %v",
						i, test.config.Host, cfg.Host)
				}
				if test.config.Account != cfg.Account {
					t.Fatalf("%d: Failed to match account. expected: %v, got: %v",
						i, test.config.Account, cfg.Account)
				}
				if test.config.User != cfg.User {
					t.Fatalf("%d: Failed to match user. expected: %v, got: %v",
						i, test.config.User, cfg.User)
				}
				if test.config.Password != cfg.Password {
					t.Fatalf("%d: Failed to match password. expected: %v, got: %v",
						i, test.config.Password, cfg.Password)
				}
				if test.config.Database != cfg.Database {
					t.Fatalf("%d: Failed to match database. expected: %v, got: %v",
						i, test.config.Database, cfg.Database)
				}
				if test.config.Schema != cfg.Schema {
					t.Fatalf("%d: Failed to match schema. expected: %v, got: %v",
						i, test.config.Schema, cfg.Schema)
				}
				if test.config.Warehouse != cfg.Warehouse {
					t.Fatalf("%d: Failed to match warehouse. expected: %v, got: %v",
						i, test.config.Warehouse, cfg.Warehouse)
				}
				if test.config.Role != cfg.Role {
					t.Fatalf("%d: Failed to match role. expected: %v, got: %v",
						i, test.config.Role, cfg.Role)
				}
				if test.config.Region != cfg.Region {
					t.Fatalf("%d: Failed to match region. expected: %v, got: %v",
						i, test.config.Region, cfg.Region)
				}
				if test.config.Protocol != cfg.Protocol {
					t.Fatalf("%d: Failed to match protocol. expected: %v, got: %v",
						i, test.config.Protocol, cfg.Protocol)
				}
				if test.config.Passcode != cfg.Passcode {
					t.Fatalf("%d: Failed to match passcode. expected: %v, got: %v",
						i, test.config.Passcode, cfg.Passcode)
				}
				if test.config.PasscodeInPassword != cfg.PasscodeInPassword {
					t.Fatalf("%d: Failed to match passcodeInPassword. expected: %v, got: %v",
						i, test.config.PasscodeInPassword, cfg.PasscodeInPassword)
				}
				if test.config.Authenticator != cfg.Authenticator {
					t.Fatalf("%d: Failed to match Authenticator. expected: %v, got: %v",
						i, test.config.Authenticator.String(), cfg.Authenticator.String())
				}
				if test.config.Authenticator == AuthTypeOkta && *test.config.OktaURL != *cfg.OktaURL {
					t.Fatalf("%d: Failed to match okta URL. expected: %v, got: %v",
						i, test.config.OktaURL, cfg.OktaURL)
				}
				if test.config.OCSPFailOpen != cfg.OCSPFailOpen {
					t.Fatalf("%d: Failed to match OCSPFailOpen. expected: %v, got: %v",
						i, test.config.OCSPFailOpen, cfg.OCSPFailOpen)
				}
				if test.ocspMode != cfg.ocspMode() {
					t.Fatalf("%d: Failed to match OCSPMode. expected: %v, got: %v",
						i, test.ocspMode, cfg.ocspMode())
				}
				if test.config.ValidateDefaultParameters != cfg.ValidateDefaultParameters {
					t.Fatalf("%d: Failed to match ValidateDefaultParameters. expected: %v, got: %v",
						i, test.config.ValidateDefaultParameters, cfg.ValidateDefaultParameters)
				}
				if test.config.ClientTimeout != cfg.ClientTimeout {
					t.Fatalf("%d: Failed to match ClientTimeout. expected: %v, got: %v",
						i, test.config.ClientTimeout, cfg.ClientTimeout)
				}
				if test.config.JWTClientTimeout != cfg.JWTClientTimeout {
					t.Fatalf("%d: Failed to match JWTClientTimeout. expected: %v, got: %v",
						i, test.config.JWTClientTimeout, cfg.JWTClientTimeout)
				}
				if test.config.ExternalBrowserTimeout != cfg.ExternalBrowserTimeout {
					t.Fatalf("%d: Failed to match ExternalBrowserTimeout. expected: %v, got: %v",
						i, test.config.ExternalBrowserTimeout, cfg.ExternalBrowserTimeout)
				}
				if test.config.TmpDirPath != cfg.TmpDirPath {
					t.Fatalf("%v: Failed to match TmpDirPatch. expected: %v, got: %v", i, test.config.TmpDirPath, cfg.TmpDirPath)
				}
			case test.err != nil:
				driverErrE, okE := test.err.(*SnowflakeError)
				driverErrG, okG := err.(*SnowflakeError)
				if okE && !okG || !okE && okG {
					t.Fatalf("%d: Wrong error. expected: %v, got: %v", i, test.err, err)
				}
				if okE && okG {
					if driverErrE.Number != driverErrG.Number {
						t.Fatalf("%d: Wrong error number. expected: %v, got: %v", i, driverErrE.Number, driverErrG.Number)
					}
				} else {
					t1 := reflect.TypeOf(err)
					t2 := reflect.TypeOf(test.err)
					if t1 != t2 {
						t.Fatalf("%d: Wrong error. expected: %T:%v, got: %T:%v", i, test.err, test.err, err, err)
					}
				}
			}

		})
	}
}

type tcDSN struct {
	cfg *Config
	dsn string
	err error
}

func TestDSN(t *testing.T) {
	tmfmt := "MM-DD-YYYY"
	testcases := []tcDSN{
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a-aofnadsf.somewhere.azure",
			},
			dsn: "u:p@a-aofnadsf.somewhere.azure.snowflakecomputing.com:443?ocspFailOpen=true&region=somewhere.azure&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a-aofnadsf.global",
			},
			dsn: "u:p@a-aofnadsf.global.snowflakecomputing.com:443?ocspFailOpen=true&region=global&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a-aofnadsf.global",
				Region:   "us-west-2",
			},
			dsn: "u:p@a-aofnadsf.global.snowflakecomputing.com:443?ocspFailOpen=true&region=global&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a-aofnadsf.global",
				Region:   "r",
			},
			err: errInvalidRegion(),
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a",
			},
			dsn: "u:p@a.snowflakecomputing.com:443?ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a",
				Region:   "us-west-2",
			},
			dsn: "u:p@a.snowflakecomputing.com:443?ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a",
				Region:   "r",
			},
			dsn: "u:p@a.r.snowflakecomputing.com:443?ocspFailOpen=true&region=r&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:                   "u",
				Password:               "p",
				Account:                "a",
				Region:                 "r",
				ExternalBrowserTimeout: 20 * time.Second,
			},
			dsn: "u:p@a.r.snowflakecomputing.com:443?externalBrowserTimeout=20&ocspFailOpen=true&region=r&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "",
				Password: "p",
				Account:  "a",
			},
			err: errEmptyUsername(),
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "",
				Account:  "a",
			},
			err: errEmptyPassword(),
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "",
			},
			err: errEmptyAccount(),
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.e",
			},
			dsn: "u:p@a.e.snowflakecomputing.com:443?ocspFailOpen=true&region=e&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.e",
				Region:   "us-west-2",
			},
			dsn: "u:p@a.e.snowflakecomputing.com:443?ocspFailOpen=true&region=e&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.e",
				Region:   "r",
			},
			err: errInvalidRegion(),
		},
		{
			cfg: &Config{
				User:               "u",
				Password:           "p",
				Account:            "a",
				Database:           "db",
				Schema:             "sc",
				Role:               "ro",
				Region:             "b",
				Authenticator:      AuthTypeSnowflake,
				Passcode:           "db",
				PasscodeInPassword: true,
				LoginTimeout:       10 * time.Second,
				RequestTimeout:     300 * time.Second,
				Application:        "special go",
			},
			dsn: "u:p@a.b.snowflakecomputing.com:443?application=special+go&database=db&loginTimeout=10&ocspFailOpen=true&passcode=db&passcodeInPassword=true&region=b&requestTimeout=300&role=ro&schema=sc&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:                           "u",
				Password:                       "p",
				Account:                        "a",
				Authenticator:                  AuthTypeExternalBrowser,
				ClientStoreTemporaryCredential: ConfigBoolTrue,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?authenticator=externalbrowser&clientStoreTemporaryCredential=true&ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:                           "u",
				Password:                       "p",
				Account:                        "a",
				Authenticator:                  AuthTypeExternalBrowser,
				ClientStoreTemporaryCredential: ConfigBoolFalse,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?authenticator=externalbrowser&clientStoreTemporaryCredential=false&ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:          "u",
				Password:      "p",
				Account:       "a",
				Authenticator: AuthTypeOkta,
				OktaURL: &url.URL{
					Scheme: "https",
					Host:   "sc.okta.com",
				},
			},
			dsn: "u:p@a.snowflakecomputing.com:443?authenticator=https%3A%2F%2Fsc.okta.com&ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.e",
				Params: map[string]*string{
					"TIMESTAMP_OUTPUT_FORMAT": &tmfmt,
				},
			},
			dsn: "u:p@a.e.snowflakecomputing.com:443?TIMESTAMP_OUTPUT_FORMAT=MM-DD-YYYY&ocspFailOpen=true&region=e&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: ":@abc",
				Account:  "a.e",
				Params: map[string]*string{
					"TIMESTAMP_OUTPUT_FORMAT": &tmfmt,
				},
			},
			dsn: "u:%3A%40abc@a.e.snowflakecomputing.com:443?TIMESTAMP_OUTPUT_FORMAT=MM-DD-YYYY&ocspFailOpen=true&region=e&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:         "u",
				Password:     "p",
				Account:      "a",
				OCSPFailOpen: OCSPFailOpenTrue,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:         "u",
				Password:     "p",
				Account:      "a",
				OCSPFailOpen: OCSPFailOpenFalse,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?ocspFailOpen=false&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:                      "u",
				Password:                  "p",
				Account:                   "a",
				ValidateDefaultParameters: ConfigBoolFalse,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?ocspFailOpen=true&validateDefaultParameters=false",
		},
		{
			cfg: &Config{
				User:                      "u",
				Password:                  "p",
				Account:                   "a",
				ValidateDefaultParameters: ConfigBoolTrue,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:         "u",
				Password:     "p",
				Account:      "a",
				InsecureMode: true,
			},
			dsn: "u:p@a.snowflakecomputing.com:443?insecureMode=true&ocspFailOpen=true&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.b.c",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.b.c",
				Region:   "us-west-2",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.b.c",
				Region:   "r",
			},
			err: errInvalidRegion(),
		},
		{
			cfg: &Config{
				User:             "u",
				Password:         "p",
				Account:          "a.b.c",
				ClientTimeout:    300 * time.Second,
				JWTClientTimeout: 60 * time.Second,
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?clientTimeout=300&jwtClientTimeout=60&ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:             "u",
				Password:         "p",
				Account:          "a.b.c",
				ClientTimeout:    300 * time.Second,
				JWTExpireTimeout: 30 * time.Second,
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?clientTimeout=300&jwtTimeout=30&ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.b.c",
				Protocol: "http",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&protocol=http&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.b.c",
				Tracing:  "debug",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&region=b.c&tracing=debug&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:                  "u",
				Password:              "p",
				Account:               "a.b.c",
				Authenticator:         AuthTypeUsernamePasswordMFA,
				ClientRequestMfaToken: ConfigBoolTrue,
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?authenticator=username_password_mfa&clientRequestMfaToken=true&ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:                  "u",
				Password:              "p",
				Account:               "a.b.c",
				Authenticator:         AuthTypeUsernamePasswordMFA,
				ClientRequestMfaToken: ConfigBoolFalse,
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?authenticator=username_password_mfa&clientRequestMfaToken=false&ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:      "u",
				Password:  "p",
				Account:   "a.b.c",
				Warehouse: "wh",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&region=b.c&validateDefaultParameters=true&warehouse=wh",
		},
		{
			cfg: &Config{
				User:     "u",
				Password: "p",
				Account:  "a.b.c",
				Token:    "t",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&region=b.c&token=t&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:          "u",
				Password:      "p",
				Account:       "a.b.c",
				Authenticator: AuthTypeTokenAccessor,
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?authenticator=tokenaccessor&ocspFailOpen=true&region=b.c&validateDefaultParameters=true",
		},
		{
			cfg: &Config{
				User:       "u",
				Password:   "p",
				Account:    "a.b.c",
				TmpDirPath: "/tmp",
			},
			dsn: "u:p@a.b.c.snowflakecomputing.com:443?ocspFailOpen=true&region=b.c&tmpDirPath=%2Ftmp&validateDefaultParameters=true",
		},
	}
	for _, test := range testcases {
		t.Run(test.dsn, func(t *testing.T) {
			dsn, err := DSN(test.cfg)
			if test.err == nil && err == nil {
				if dsn != test.dsn {
					t.Errorf("failed to get DSN. expected: %v, got:\n %v", test.dsn, dsn)
				}
				_, err := ParseDSN(dsn)
				if err != nil {
					t.Errorf("failed to parse DSN. dsn: %v, err: %v", dsn, err)
				}
			}
			if test.err != nil && err == nil {
				t.Errorf("expected error. dsn: %v, err: %v", test.dsn, test.err)
			}
			if err != nil && test.err == nil {
				t.Errorf("failed to match. err: %v", err)
			}
		})
	}
}

func TestParsePrivateKeyFromFileMissingFile(t *testing.T) {
	_, err := parsePrivateKeyFromFile("nonexistent")

	if err == nil {
		t.Error("should report error for nonexistent file")
	}
}

func TestParsePrivateKeyFromFileIncorrectData(t *testing.T) {
	pemFile := createTmpFile("exampleKey.pem", []byte("gibberish"))
	_, err := parsePrivateKeyFromFile(pemFile)

	if err == nil {
		t.Error("should report error for wrong data in file")
	}
}

func TestParsePrivateKeyFromFile(t *testing.T) {
	generatedKey, _ := rsa.GenerateKey(cr.Reader, 1024)
	pemKey, _ := x509.MarshalPKCS8PrivateKey(generatedKey)
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pemKey,
		},
	)
	keyFile := createTmpFile("exampleKey.pem", pemData)
	defer os.Remove(keyFile)

	parsedKey, err := parsePrivateKeyFromFile(keyFile)
	if err != nil {
		t.Errorf("unable to parse pam file from path: %v, err: %v", keyFile, err)
	} else if !parsedKey.Equal(generatedKey) {
		t.Errorf("generated key does not equal to parsed key from file\ngeneratedKey=%v\nparsedKey=%v",
			generatedKey, parsedKey)
	}
}

func createTmpFile(fileName string, content []byte) string {
	tempFile, _ := os.CreateTemp("", fileName)
	tempFile.Write(content)
	absolutePath := tempFile.Name()
	return absolutePath
}

type configParamToValue struct {
	configParam string
	value       string
}

func TestGetConfigFromEnv(t *testing.T) {
	envMap := map[string]configParamToValue{
		"SF_TEST_ACCOUNT":     {"Account", "account"},
		"SF_TEST_USER":        {"User", "user"},
		"SF_TEST_PASSWORD":    {"Password", "password"},
		"SF_TEST_ROLE":        {"Role", "role"},
		"SF_TEST_HOST":        {"Host", "host"},
		"SF_TEST_PORT":        {"Port", "8080"},
		"SF_TEST_PROTOCOL":    {"Protocol", "http"},
		"SF_TEST_WAREHOUSE":   {"Warehouse", "warehouse"},
		"SF_TEST_DATABASE":    {"Database", "database"},
		"SF_TEST_REGION":      {"Region", "region"},
		"SF_TEST_PASSCODE":    {"Passcode", "passcode"},
		"SF_TEST_SCHEMA":      {"Schema", "schema"},
		"SF_TEST_APPLICATION": {"Application", "application"},
	}
	var properties = make([]*ConfigParam, len(envMap))
	i := 0
	for key, ctv := range envMap {
		os.Setenv(key, ctv.value)
		cfgParam := ConfigParam{ctv.configParam, key, true}
		properties[i] = &cfgParam
		i++
	}
	defer func() {
		for key := range envMap {
			os.Unsetenv(key)
		}
	}()

	cfg, err := GetConfigFromEnv(properties)
	if err != nil {
		t.Errorf("unable to parse env variables to Config, err: %v", err)
	}

	err = checkConfig(*cfg, envMap)
	if err != nil {
		t.Error(err)
	}
}

func checkConfig(cfg Config, envMap map[string]configParamToValue) error {
	appendError := func(errArray []string, envName string, expected string, received string) []string {
		errArray = append(errArray, fmt.Sprintf("field %v expected value: %v, received value: %v", envName, expected, received))
		return errArray
	}

	value := reflect.ValueOf(cfg)
	typeOfCfg := value.Type()
	cfgValues := make(map[string]interface{}, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		cfgValues[typeOfCfg.Field(i).Name] = value.Field(i).Interface()
	}

	var errArray []string
	for key, ctv := range envMap {
		if ctv.configParam == "Port" {
			if portStr := strconv.Itoa(cfgValues[ctv.configParam].(int)); portStr != ctv.value {
				errArray = appendError(errArray, key, ctv.value, cfgValues[ctv.configParam].(string))
			}
		} else if cfgValues[ctv.configParam] != ctv.value {
			errArray = appendError(errArray, key, ctv.value, cfgValues[ctv.configParam].(string))
		}
	}

	if errArray != nil {
		return fmt.Errorf(strings.Join(errArray, "\n"))
	}

	return nil
}

func TestConfigValidateTmpDirPath(t *testing.T) {
	cfg := &Config{
		TmpDirPath: "/not/existing",
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("Should fail on not existing TmpDirPath")
	}
}
