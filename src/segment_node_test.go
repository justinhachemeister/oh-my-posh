package main

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestNodeMatchesVersionFile(t *testing.T) {
	cases := []struct {
		Case      string
		Expected  bool
		RCVersion string
		Version   string
	}{
		{Case: "no file context", Expected: true, RCVersion: "", Version: "durp"},
		{Case: "version match", Expected: true, RCVersion: "durp", Version: "durp"},
		{Case: "version mismatch", Expected: false, RCVersion: "werp", Version: "durp"},
	}

	for _, tc := range cases {
		env := new(MockedEnvironment)
		env.On("getFileContent", ".nvmrc").Return(tc.RCVersion)
		node := &node{
			language: &language{
				env: env,
				activeCommand: &cmd{
					version: &version{
						full: tc.Version,
					},
				},
			},
		}
		assert.Equal(t, tc.Expected, node.matchesVersionFile(), tc.Case)
	}
}

func TestNodeInContext(t *testing.T) {
	cases := []struct {
		Case           string
		HasYarn        bool
		hasNPM         bool
		hasDefault     bool
		PkgMgrEnabled  bool
		ExpectedString string
	}{
		{Case: "no package manager file", ExpectedString: "", PkgMgrEnabled: true},
		{Case: "yarn", HasYarn: true, ExpectedString: "yarn", PkgMgrEnabled: true},
		{Case: "npm", hasNPM: true, ExpectedString: "npm", PkgMgrEnabled: true},
		{Case: "default", hasDefault: true, ExpectedString: "npm", PkgMgrEnabled: true},
		{Case: "disabled", HasYarn: true, ExpectedString: "", PkgMgrEnabled: false},
	}

	for _, tc := range cases {
		env := new(MockedEnvironment)
		env.On("hasFiles", "yarn.lock").Return(tc.HasYarn)
		env.On("hasFiles", "package-lock.json").Return(tc.hasNPM)
		env.On("hasFiles", "package.json").Return(tc.hasDefault)
		node := &node{
			language: &language{
				env: env,
				props: &properties{
					values: map[Property]interface{}{
						YarnIcon:              "yarn",
						NPMIcon:               "npm",
						DisplayPackageManager: tc.PkgMgrEnabled,
					},
				},
			},
		}
		node.loadContext()
		assert.Equal(t, tc.ExpectedString, node.packageManagerIcon, tc.Case)
	}
}