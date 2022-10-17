package tcp

import (
	"reflect"
	"testing"

	"github.com/panjf2000/gnet/v2"
)

func TestEchoServer_OnBoot(t *testing.T) {
	type fields struct {
		BuiltinEventEngine gnet.BuiltinEventEngine
		eng                gnet.Engine
		addr               string
		multicore          bool
		sessions           map[string]string
	}
	type args struct {
		eng gnet.Engine
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   gnet.Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &EchoServer{
				BuiltinEventEngine: tt.fields.BuiltinEventEngine,
				eng:                tt.fields.eng,
				addr:               tt.fields.addr,
				multicore:          tt.fields.multicore,
				sessions:           tt.fields.sessions,
			}
			if got := es.OnBoot(tt.args.eng); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EchoServer.OnBoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEchoServer_OnOpen(t *testing.T) {
	type fields struct {
		BuiltinEventEngine gnet.BuiltinEventEngine
		eng                gnet.Engine
		addr               string
		multicore          bool
		sessions           map[string]string
	}
	type args struct {
		c gnet.Conn
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOut    []byte
		wantAction gnet.Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &EchoServer{
				BuiltinEventEngine: tt.fields.BuiltinEventEngine,
				eng:                tt.fields.eng,
				addr:               tt.fields.addr,
				multicore:          tt.fields.multicore,
				sessions:           tt.fields.sessions,
			}
			gotOut, gotAction := es.OnOpen(tt.args.c)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("EchoServer.OnOpen() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
			if !reflect.DeepEqual(gotAction, tt.wantAction) {
				t.Errorf("EchoServer.OnOpen() gotAction = %v, want %v", gotAction, tt.wantAction)
			}
		})
	}
}

func TestEchoServer_OnTraffic(t *testing.T) {
	type fields struct {
		BuiltinEventEngine gnet.BuiltinEventEngine
		eng                gnet.Engine
		addr               string
		multicore          bool
		sessions           map[string]string
	}
	type args struct {
		c gnet.Conn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   gnet.Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &EchoServer{
				BuiltinEventEngine: tt.fields.BuiltinEventEngine,
				eng:                tt.fields.eng,
				addr:               tt.fields.addr,
				multicore:          tt.fields.multicore,
				sessions:           tt.fields.sessions,
			}
			if got := es.OnTraffic(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EchoServer.OnTraffic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		port string
	}
	tests := []struct {
		name string
		args args
		want *EchoServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
