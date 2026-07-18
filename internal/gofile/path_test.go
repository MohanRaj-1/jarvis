package gofile

import "testing"

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{name: "cleans Go path", path: "dir/../sample.GO", want: "sample.GO"},
		{name: "rejects empty path", path: " ", wantErr: true},
		{name: "rejects non-Go file", path: "sample.txt", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ValidatePath(test.path)
			if (err != nil) != test.wantErr {
				t.Fatalf("ValidatePath(%q) error = %v, wantErr %t", test.path, err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("ValidatePath(%q) = %q, want %q", test.path, got, test.want)
			}
		})
	}
}
