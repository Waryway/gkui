package bootstrap

import (
	"context"
	"gkui/pkg/logstream"
	"testing"
)

type exampleSettings struct {
	A string `yaml:"a"`
}

//func TestConfig_Init(t *testing.T) {
//	type testCase[boot BootStrap] struct {
//		name string
//		c    Config[boot]
//		want *Config[boot]
//	}
//	tests := []testCase[exampleSettings]{
//		{
//			name: "some test",
//			c: ,
//			want:,
//		},
//	},
//		for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.Init(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Init() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//
//func TestConfig_Load(t *testing.T) {
//	type testCase[boot BootStrap] struct {
//		name string
//		c    Config[boot]
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.c.Load()
//		})
//	}
//}
//
//func TestConfig_LoadEnv(t *testing.T) {
//	type testCase[boot BootStrap] struct {
//		name string
//		c    Config[boot]
//		want *Config[boot]
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.LoadEnv(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("LoadEnv() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestConfig_LoadYaml(t *testing.T) {
	bootstrap := Config[exampleSettings]{BootStrap: exampleSettings{}}
	expected := Config[exampleSettings]{BootStrap: exampleSettings{A: "thing"}}

	bCtx := context.Background()
	ctx, cancel := context.WithCancel(bCtx)
	ls := logstream.InitLogStream(ctx, cancel)
	bootstrap.Init(&ls)
	expected.Init(&ls)

	data := []byte(`
a: thing
b: a string from struct B`)
	bootstrap.file = &data
	type testCase[boot BootStrap] struct {
		name string
		c    Config[boot]
		want *Config[boot]
	}

	tests := []testCase[exampleSettings]{
		{
			name: "some test",
			c:    bootstrap,
			want: &expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				tt.c.BootCh <- *tt.c.LoadYaml()
			}()

			select {
			case err := <-tt.c.ErrCh:
				t.Error(err)
			case got := <-tt.c.BootCh:
				if got.BootStrap.A != "thing" {
					t.Errorf("LoadYaml() = %v, want %v", got.BootStrap.A, tt.want.BootStrap.A)
				}
			}
		})
	}
}

//
//func TestConfig_ReadSettingFile(t *testing.T) {
//	type testCase[boot BootStrap] struct {
//		name string
//		c    Config[boot]
//		want *Config[boot]
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.c.ReadSettingFile(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ReadSettingFile() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestConfig_Save(t *testing.T) {
//	type testCase[boot BootStrap] struct {
//		name string
//		c    Config[boot]
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.c.Save()
//		})
//	}
//}

func Test_loadFromEnv(t *testing.T) {
	type args struct {
		k       string
		def     string
		errorCh chan error
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := loadFromEnv(tt.args.k, tt.args.def, tt.args.errorCh)
			if got != tt.want {
				t.Errorf("loadFromEnv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("loadFromEnv() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
