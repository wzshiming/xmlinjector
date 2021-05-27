package xmlinjector

import (
	"reflect"
	"testing"
)

func TestInject(t *testing.T) {
	inject := func(args, origin []byte) []byte {
		return append(args[:len(args):len(args)], origin...)
	}
	type args struct {
		key    []byte
		data   []byte
		inject func(args, origin []byte) []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				key:    []byte(""),
				data:   []byte(""),
				inject: inject,
			},
			wantErr: true,
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("xxxx"),
				inject: inject,
			},
			want: []byte("xxxx"),
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("<!--key args... --><!--/key-->"),
				inject: inject,
			},
			want: []byte("<!--key args... -->args...<!--/key-->"),
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("<!--key args... -->XXXXXXXX<!--/key-->"),
				inject: inject,
			},
			want: []byte("<!--key args... -->args...XXXXXXXX<!--/key-->"),
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("<!--key args... -->XXXXXXXX<!--/key--><!--key args xxxx --><!--/key-->"),
				inject: inject,
			},
			want: []byte("<!--key args... -->args...XXXXXXXX<!--/key--><!--key args xxxx -->args xxxx<!--/key-->"),
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("<!--key args... / -->"),
				inject: inject,
			},
			want: []byte("<!--key args... / -->"),
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("<!--key args... /-->"),
				inject: inject,
			},
			want: []byte("<!--key args... -->args...<!-- /key -->"),
		},
		{
			args: args{
				key:    []byte("key"),
				data:   []byte("<!--key /--><!--key args... /-->"),
				inject: inject,
			},
			want: []byte("<!--key /--><!--key args... -->args...<!-- /key -->"),
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key args... /-->"),
				inject: func(args, origin []byte) []byte {
					return args
				},
			},
			want: []byte("<!--key args... -->args...<!-- /key -->"),
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key\nargs...\n/-->"),
				inject: func(args, origin []byte) []byte {
					return args
				},
			},
			want: []byte("<!--key\nargs...\n-->args...<!-- /key -->"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Inject(tt.args.key, tt.args.data, tt.args.inject)
			if (err != nil) != tt.wantErr {
				t.Errorf("Inject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inject() got = %q, want %q", got, tt.want)
			}
			if len(got) != cap(got) {
				t.Errorf("Inject() len != cap, len = %d, cap = %d", len(got), cap(got))
			}
		})
	}
}
