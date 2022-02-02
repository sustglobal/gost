package component

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoadFromEnv(t *testing.T) {
	type Obj struct {
		A_Str  string        `env:"ENV_A"`
		B_Int  int           `env:"ENV_B"`
		C_Dur  time.Duration `env:"ENV_C"`
		D_Bool bool          `env:"ENV_D"`
	}

	var want Obj
	want.A_Str = "xyz"
	want.B_Int = 11
	want.C_Dur = 14 * time.Second
	want.D_Bool = true

	var got Obj
	os.Setenv("ENV_A", "xyz")
	os.Setenv("ENV_B", "11")
	os.Setenv("ENV_C", "14s")
	os.Setenv("ENV_D", "yes")

	err := LoadFromEnv(&got)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected ouptput: want=%+v got=%+v", want, got)
	}
}

func TestLoadFromEnvUnset(t *testing.T) {
	type Obj struct {
		A_Str  string        `env:"ENV_A"`
		B_Int  int           `env:"ENV_B"`
		C_Dur  time.Duration `env:"ENV_C"`
		D_Bool bool          `env:"ENV_D"`
	}

	var want Obj
	want.A_Str = ""
	want.B_Int = 0
	want.C_Dur = 0 * time.Second
	want.D_Bool = false

	var got Obj
	os.Unsetenv("ENV_A")
	os.Unsetenv("ENV_B")
	os.Unsetenv("ENV_C")
	os.Unsetenv("ENV_D")

	err := LoadFromEnv(&got)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected ouptput: want=%+v got=%+v", want, got)
	}
}

func TestLoadFromEnvInvalidInt(t *testing.T) {
	type Obj struct {
		Int int `env:"INT"`
	}

	var got Obj
	os.Setenv("INT", "xyz")

	err := LoadFromEnv(&got)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestLoadFromEnvInvalidDuration(t *testing.T) {
	type Obj struct {
		Dur time.Duration `env:"DUR"`
	}

	var got Obj
	os.Setenv("DUR", "xyz")

	err := LoadFromEnv(&got)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestLoadFromEnvBoolFalse(t *testing.T) {
	type Obj struct {
		Bool bool `env:"BOOL"`
	}

	var got Obj
	os.Setenv("BOOL", "xyz") // non-truthy value

	var want Obj
	want.Bool = false

	err := LoadFromEnv(&got)
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected ouptput: want=%+v got=%+v", want, got)
	}
}
