package mongodb

import (
	"os"
	"testing"
)

func Test_getAppName(t *testing.T) {
	tests := []struct {
		name string
		want func() string
	}{
		{
			name: "test case 1",
			want: func() string {
				name := "TEST_SERVICE"
				if err := os.Setenv("APP_NAME", name); err != nil {
					return "Go-Service"
				}
				return name
			},
		},
		{
			name: "test case 2",
			want: func() string {
				// delete the env variable
				if err := os.Setenv("APP_NAME", ""); err != nil {
					return os.Getenv("APP_NAME")
				}
				return "Go-Service"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want()
			if got := getAppName(); got != want {
				t.Errorf("getAppName() = %v, want %v", got, want)
			}
		})
	}
}

func Test_getDbName(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want string
	}{
		{
			name: "test case 1",
			uri:  "mongodb://localhost:27017/bankaoolbio_devDB?authSource=admin&replicaSet=rs_dev",
			want: "bankaoolbio_devDB",
		},
		{
			name: "test case 2",
			uri:  "mongodb://localhost:27017/?authSource=admin",
			want: "admin",
		},
		{
			name: "test case 3",
			uri:  "mongodb://localhost:27017/bankaoolbio_devDB",
			want: "bankaoolbio_devDB",
		},
		{
			name: "test case 4",
			uri:  "mongodb://localhost:27017/",
			want: "",
		},
		{
			name: "test case 5",
			uri:  "mongodb://localhost:27017?authSource=admin",
			want: "",
		},
		{
			name: "test case 6",
			uri:  "mongodb://username:password@host_01:27017,host_02:27017,host_03:27017/bankaoolbio_devDB?authSource=admin",
			want: "bankaoolbio_devDB",
		},
		{
			name: "test case 7",
			uri:  "mongodb://username:password@host_01:27017,host_02:27017,host_03:27017/?authSource=admin",
			want: "admin",
		},
		{
			name: "test case 8",
			uri:  "mongodb://username:password@host_01:27017,host_02:27017,host_03:27017",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDbName(tt.uri); got != tt.want {
				t.Errorf("getDbName() = %v, want %v", got, tt.want)
			}
		})
	}
}
