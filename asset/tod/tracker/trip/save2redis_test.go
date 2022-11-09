package trip

import (
	"reflect"
	"testing"

	"github.com/go-redis/redis/v8"
)

func TestS2Redis_Connect(t *testing.T) {
	type fields struct {
		Server   string
		Password string
		Client   *redis.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s2r := &S2Redis{
				Server:   tt.fields.Server,
				Password: tt.fields.Password,
				Client:   tt.fields.Client,
			}
			if err := s2r.Connect(); (err != nil) != tt.wantErr {
				t.Errorf("S2Redis.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestS2Redis_SaveTripDetail(t *testing.T) {
	type fields struct {
		Server   string
		Password string
		Client   *redis.Client
	}
	type args struct {
		id  string
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s2r := &S2Redis{
				Server:   tt.fields.Server,
				Password: tt.fields.Password,
				Client:   tt.fields.Client,
			}
			if err := s2r.SaveTripDetail(tt.args.id, tt.args.buf); (err != nil) != tt.wantErr {
				t.Errorf("S2Redis.SaveTripDetail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestS2Redis_TripDetail(t *testing.T) {
	type fields struct {
		Server   string
		Password string
		Client   *redis.Client
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s2r := &S2Redis{
				Server:   tt.fields.Server,
				Password: tt.fields.Password,
				Client:   tt.fields.Client,
			}
			got, err := s2r.TripDetail(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("S2Redis.TripDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S2Redis.TripDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS2Redis_AllTripDetail(t *testing.T) {
	type fields struct {
		Server   string
		Password string
		Client   *redis.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s2r := &S2Redis{
				Server:   tt.fields.Server,
				Password: tt.fields.Password,
				Client:   tt.fields.Client,
			}
			got, err := s2r.AllTripDetail()
			if (err != nil) != tt.wantErr {
				t.Errorf("S2Redis.AllTripDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("S2Redis.AllTripDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS2Redis_ClearTripDetail(t *testing.T) {
	type fields struct {
		Server   string
		Password string
		Client   *redis.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s2r := &S2Redis{
				Server:   tt.fields.Server,
				Password: tt.fields.Password,
				Client:   tt.fields.Client,
			}
			if err := s2r.ClearTripDetail(); (err != nil) != tt.wantErr {
				t.Errorf("S2Redis.ClearTripDetail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
