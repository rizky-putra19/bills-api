package converter

import (
	"reflect"
	"testing"
)

func TestToString(t *testing.T) {
	type args struct {
		v interface{}
	}
	type Person struct {
		Name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nil value",
			args: args{
				nil,
			},
			want: "",
		},
		{
			name: "string correct format",
			args: args{
				"123 asdf",
			},
			want: "123 asdf",
		},
		{
			name: "int correct format",
			args: args{
				123,
			},
			want: "123",
		},
		{
			name: "int32 correct format",
			args: args{
				int32(123),
			},
			want: "123",
		},
		{
			name: "int64 correct format",
			args: args{
				int64(123),
			},
			want: "123",
		},
		{
			name: "bool correct format",
			args: args{
				true,
			},
			want: "true",
		},
		{
			name: "float32 correct format",
			args: args{
				float32(123.456),
			},
			want: "123.456",
		},
		{
			name: "float64 correct format",
			args: args{
				123.456,
			},
			want: "123.456",
		},
		{
			name: "uint8 list correct format",
			args: args{
				[]uint8{1, 2, 3},
			},
			want: "",
		},
		{
			name: "default value with error",
			args: args{
				make(chan int),
			},
			want: "",
		},
		{
			name: "default value no error",
			args: args{
				Person{"toped"},
			},
			want: `{"Name":"toped"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.v); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "type string, accepted format, return true",
			args: args{
				v: "true",
			},
			want: true,
		},
		{
			name: "type string, accepted format, return false",
			args: args{
				v: "f",
			},
			want: false,
		},
		{
			name: "type string, unaccepted format, return false",
			args: args{
				v: "betul betul betul",
			},
			want: false,
		},
		{
			name: "type int = 0, return false",
			args: args{
				v: 0,
			},
			want: false,
		},
		{
			name: "type int != 0, return true",
			args: args{
				v: 1,
			},
			want: true,
		},
		{
			name: "default, return false",
			args: args{
				v: []byte("asd"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBool(tt.args.v); got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "nil value",
			args: args{
				nil,
			},
			want: 0,
		},
		{
			name: "string correct format",
			args: args{
				"123",
			},
			want: 123,
		},
		{
			name: "string bad format",
			args: args{
				"asdfd",
			},
			want: 0,
		},
		{
			name: "int correct format",
			args: args{
				123,
			},
			want: 123,
		},
		{
			name: "int32 correct format",
			args: args{
				int32(123),
			},
			want: 123,
		},
		{
			name: "int64 correct format",
			args: args{
				int64(123),
			},
			want: 123,
		},
		{
			name: "bool correct format",
			args: args{
				true,
			},
			want: 0,
		},
		{
			name: "float32 correct format",
			args: args{
				float32(123.00),
			},
			want: 123,
		},
		{
			name: "float64 correct format",
			args: args{
				123.456,
			},
			want: 123,
		},
		{
			name: "[]byte correct format",
			args: args{
				[]byte("123"),
			},
			want: 123,
		},
		{
			name: "string with decimal value",
			args: args{
				[]byte("999.00"),
			},
			want: 999,
		},
		{
			name: "default value with error",
			args: args{
				make(chan int),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.args.v); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToIntWithDecimal(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "nil value",
			args: args{
				nil,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt(tt.args.v); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestToInt64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "nil value",
			args: args{
				nil,
			},
			want: 0,
		},
		{
			name: "string correct format",
			args: args{
				"123",
			},
			want: 123,
		},
		{
			name: "string bad format",
			args: args{
				"asdfd",
			},
			want: 0,
		},
		{
			name: "int correct format",
			args: args{
				123,
			},
			want: 123,
		},
		{
			name: "int32 correct format",
			args: args{
				int32(123),
			},
			want: 123,
		},
		{
			name: "int64 correct format",
			args: args{
				int64(123),
			},
			want: 123,
		},
		{
			name: "bool correct format",
			args: args{
				true,
			},
			want: 0,
		},
		{
			name: "float32 correct format",
			args: args{
				float32(123.00),
			},
			want: 123,
		},
		{
			name: "float64 correct format",
			args: args{
				123.456,
			},
			want: 123,
		},
		{
			name: "[]byte correct format",
			args: args{
				[]byte("123"),
			},
			want: 123,
		},
		{
			name: "default value with error",
			args: args{
				make(chan int),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt64(tt.args.v); got != tt.want {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToArrayOfInt(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Args that passed is not string, return nil",
			args: args{
				v: 1,
			},
			want: nil,
		},
		{
			name: "Args that passed is string invalid, return nil",
			args: args{
				v: "[1,2",
			},
			want: nil,
		},
		{
			name: "Args that passed is string valid, return nil",
			args: args{
				v: "[1,2]",
			},
			want: []int{1, 2},
		},
		{
			name: "Args that passed is slice string valid, return nil",
			args: args{
				v: []string{"1", "2"},
			},
			want: []int{1, 2},
		},
		{
			name: "Args that passed is byte array, return nil",
			args: args{
				v: [][]byte{
					[]byte("1234"),
					[]byte("5678"),
				},
			},
			want: []int{1234, 5678},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToArrayOfInt(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToArrayOfInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToArrayOfString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Args that passed is not string, return nil",
			args: args{
				v: 1,
			},
			want: nil,
		},
		{
			name: "Args that passed is string invalid, return nil",
			args: args{
				v: "[1,2",
			},
			want: nil,
		},
		{
			name: "Args that passed is int array, return nil",
			args: args{
				v: "[\"1\",\"2\"]",
			},
			want: []string{"1", "2"},
		},
		{
			name: "Args that passed is string array, return nil",
			args: args{
				v: "[\"1\",\"2\"]",
			},
			want: []string{"1", "2"},
		},
		{
			name: "Args that passed is byte array, return nil",
			args: args{
				v: [][]byte{
					[]byte("foo"),
					[]byte("bar"),
				},
			},
			want: []string{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToArrayOfString(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToArrayOfString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromStringToIntAmount(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "10K with . delimiter",
			args: args{
				val: "1000.000",
			},
			want: 1000,
		},

		{
			name: "10K with no delimiter",
			args: args{
				val: "1000",
			},
			want: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromStringToIntAmount(tt.args.val); got != tt.want {
				t.Errorf("FromStringToIntAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}
