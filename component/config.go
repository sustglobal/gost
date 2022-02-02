package component

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BindHTTPServer          string        `env:"GUST_BIND_HTTP_SERVER"`
	ExposeMetrics           bool          `env:"GUST_EXPOSE_METRICS"`
	ExposeHealth            bool          `env:"GUST_EXPOSE_HEALTH"`
	GracefulShutdownTimeout time.Duration `env:"GUST_GRACEFUL_SHUTDOWN_TIMEOUT"`
}

func DefaultConfig() Config {
	return Config{
		BindHTTPServer:          "0.0.0.0:8080",
		GracefulShutdownTimeout: 60 * time.Second,
	}
}

var (
	truthyValues = []string{"y", "yes", "1", "true", "t"}
)

func LoadFromEnv(obj interface{}) error {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()

	for i := 0; i < typ.NumField(); i++ {
		tagValue := typ.Field(i).Tag.Get("env")
		field := val.Elem().Field(i)

		if tagValue == "" {
			continue
		}

		envValStr := os.Getenv(tagValue)

		// Skip if env unset or set to empty string
		if envValStr == "" {
			continue
		}

		switch field.Type() {
		case reflect.TypeOf(""):
			field.SetString(envValStr)
		case reflect.TypeOf(1):
			envValInt, err := strconv.Atoi(envValStr)
			if err != nil {
				return fmt.Errorf("unable to parse %s as int: %v", tagValue, err)
			}
			field.SetInt(int64(envValInt))
		case reflect.TypeOf(false):
			envValStrLower := strings.ToLower(envValStr)
			var envValBool bool
			for _, val := range truthyValues {
				if envValStrLower == val {
					envValBool = true
				}
			}
			field.SetBool(envValBool)
		case reflect.TypeOf(time.Second):
			envValDur, err := time.ParseDuration(envValStr)
			if err != nil {
				return fmt.Errorf("unable to parse %s as time.Duration: %v", tagValue, err)
			}
			field.Set(reflect.ValueOf(envValDur))
		default:
			return fmt.Errorf("unsupported field type %+v of key %v", field.Type(), tagValue)
		}
	}

	return nil
}
