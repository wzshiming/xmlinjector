package xmlinjector

import (
	"reflect"
	"testing"
)

func Test_scanAnnotationElement(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
		wantBegin   int
		wantEnd     int
		wantErr     bool
	}{
		{
			args: args{
				data: []byte(""),
			},
			wantErr: true,
		},
		{
			args: args{
				data: []byte("<!-- "),
			},
			wantErr: true,
		},
		{
			args: args{
				data: []byte("<!--X-->"),
			},
			wantContent: []byte("X"),
			wantBegin:   0,
			wantEnd:     8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, gotBegin, gotEnd, err := scanAnnotationElement(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanAnnotationElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("scanAnnotationElement() gotContent = %q, want %q", gotContent, tt.wantContent)
			}
			if gotBegin != tt.wantBegin {
				t.Errorf("scanAnnotationElement() gotBegin = %v, want %v", gotBegin, tt.wantBegin)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("scanAnnotationElement() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func Test_scanPairAnnotationElement(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := []struct {
		name       string
		args       args
		wantArgs   []byte
		wantBegin  int
		wantEnd    int
		wantSingle bool
		wantErr    bool
	}{
		{
			args: args{
				key:  []byte(""),
				data: []byte(""),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte(""),
			},
			wantErr: true,
		},

		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--x key--><!--/key-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key--><!--/key-->"),
			},
			wantArgs:  nil,
			wantBegin: 10,
			wantEnd:   10,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--k-->      <!--/k-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key-->      <!--/k-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key-->      <!--key-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key2--><!--/key2-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key--><!--/key2-->"),
			},
			wantErr: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!-- key --><!-- /key -->"),
			},
			wantArgs:  nil,
			wantBegin: 12,
			wantEnd:   12,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key--> <!--/key-->"),
			},
			wantArgs:  nil,
			wantBegin: 10,
			wantEnd:   11,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!-- key --> <!-- /key -->"),
			},
			wantArgs:  nil,
			wantBegin: 12,
			wantEnd:   13,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--key xxxx--> <!--/key-->"),
			},
			wantArgs:  []byte("xxxx"),
			wantBegin: 15,
			wantEnd:   16,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!-- key xxxx --> <!-- /key -->"),
			},
			wantArgs:  []byte("xxxx"),
			wantBegin: 17,
			wantEnd:   18,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!-- key xxxx /-->"),
			},
			wantArgs:   []byte("xxxx"),
			wantBegin:  14,
			wantEnd:    15,
			wantSingle: true,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--\nkey\nxxxx\n--> <!--\n/key\n-->"),
			},
			wantArgs:  []byte("xxxx"),
			wantBegin: 17,
			wantEnd:   18,
		},
		{
			args: args{
				key:  []byte("key"),
				data: []byte("<!--\nkey\nxxxx\n/-->"),
			},
			wantArgs:   []byte("xxxx"),
			wantBegin:  14,
			wantEnd:    15,
			wantSingle: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotBegin, gotEnd, gotSingle, err := scanPairAnnotationElement(tt.args.key, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("scanPairAnnotationElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("scanPairAnnotationElement() gotArgs = %q, want %q", gotArgs, tt.wantArgs)
			}
			if gotBegin != tt.wantBegin {
				t.Errorf("scanPairAnnotationElement() gotBegin = %v, want %v", gotBegin, tt.wantBegin)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("scanPairAnnotationElement() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
			if gotSingle != tt.wantSingle {
				t.Errorf("scanPairAnnotationElement() gotSingle = %v, want %v", gotSingle, tt.wantSingle)
			}
		})
	}
}
