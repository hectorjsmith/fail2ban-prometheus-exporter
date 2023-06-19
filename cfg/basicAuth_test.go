package cfg

import "testing"

func Test_hashedBasicAuth_DoesBasicAuthMatch(t *testing.T) {
	type args struct {
		username string
		password string
	}
	type fields struct {
		username string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Happy test #1", fields{username: "1234", password: "1234"}, args{username: "1234", password: "1234"}, true},
		{"Happy test #2", fields{username: "test", password: "1234"}, args{username: "test", password: "1234"}, true},
		{"Happy test #3", fields{username: "TEST", password: "1234"}, args{username: "TEST", password: "1234"}, true},
		{"Non match #1", fields{username: "test", password: "1234"}, args{username: "1234", password: "1234"}, false},
		{"Non match #2", fields{username: "1234", password: "test"}, args{username: "1234", password: "1234"}, false},
		{"Non match #3", fields{username: "1234", password: "test"}, args{username: "1234", password: "TEST"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basicAuth := newHashedBasicAuth(tt.fields.username, tt.fields.password)
			if got := basicAuth.DoesBasicAuthMatch(tt.args.username, tt.args.password); got != tt.want {
				t.Errorf("DoesBasicAuthMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashedBasicAuth_Enabled(t *testing.T) {
	type fields struct {
		username string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"Both blank", fields{username: "", password: ""}, false},
		{"Single blank #1", fields{username: "test", password: ""}, false},
		{"Single blank #1", fields{username: "", password: "test"}, false},
		{"Both populated", fields{username: "test", password: "test"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basicAuth := newHashedBasicAuth(tt.fields.username, tt.fields.password)
			if got := basicAuth.Enabled(); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
