package test

import (
	"aliyun-oss/utils"
	"testing"
)

func TestProjectRoot(t *testing.T) {
	tests := []struct {
		expect string
	}{
		{"/home/fredom/workspace/goland/aliyun-bucket"},
	}

	for _, test := range tests {
		if utils.ProjectRoot() != test.expect {
			t.Errorf(
				"output: %s not matched expected: %s",
				utils.ProjectRoot(),
				test.expect,
			)
		}
	}

	utils.ProjectRoot()
}
