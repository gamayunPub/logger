package logger

import "testing"

func TestNew(t *testing.T) {
	testCases := []struct {
		cfg        *Config
		shouldFail bool
	}{
		{
			cfg: &Config{},
		},
		{
			cfg: &Config{
				Debug:       false,
				Level:       "",
				Output:      []string{},
				TimeEncoder: "",
			},
		},
		{
			cfg: &Config{
				Debug:       false,
				Level:       "",
				Output:      []string{},
				TimeEncoder: "epoch",
			},
		},
		{
			cfg: &Config{
				Debug:       false,
				Level:       InfoLevel,
				Output:      []string{},
				TimeEncoder: "epoch",
			},
		},
		{
			cfg: &Config{
				Debug:       false,
				Level:       "",
				Output:      []string{},
				TimeEncoder: "test",
			},
		},
		{
			cfg: &Config{
				Debug:       false,
				Level:       "test",
				Output:      []string{},
				TimeEncoder: "epoch",
			},
			shouldFail: true,
		},
	}

	var (
		i  int
		tc struct {
			cfg        *Config
			shouldFail bool
		}
	)

	defer func() {
		if recoveryMessage := recover(); recoveryMessage != nil && !tc.shouldFail {
			t.Fatalf("test case %d paniced: %s", i+1, recoveryMessage)
		}
	}()

	for i, tc = range testCases {
		if _, err := New(tc.cfg, "", ""); err != nil && !tc.shouldFail {
			t.Fatalf("test case %d: %s", i+1, err)
		}
	}
}
