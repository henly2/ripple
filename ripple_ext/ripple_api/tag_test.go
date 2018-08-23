package ripple_api

import (
	"math"
	"testing"
)

func TestTag(t *testing.T) {
	type TagInfo struct {
		Str   string
		Tag   uint32
		NoErr bool

		Str2   string
		Equ   bool
	}

	tags := []TagInfo{
		TagInfo{Str: "", Tag: 0, NoErr: true, Str2:"", Equ: true},
		TagInfo{Str: "0", Tag: 0, NoErr: true, Str2:"0", Equ: true},
		TagInfo{Str: "0", Tag: 0, NoErr: true, Str2:"1", Equ: false},
		TagInfo{Str: "abcd", Tag: 0, NoErr: false, Str2:"", Equ: false},
		TagInfo{Str: "1234", Tag: 1234, NoErr: true, Str2:"1234", Equ: true},
		TagInfo{Str: "1234", Tag: 1234, NoErr: true, Str2:"12345", Equ: false},
		TagInfo{Str: "4294967295", Tag: math.MaxUint32, NoErr: true, Str2:"4294967295", Equ: true},
		TagInfo{Str: "4294967295", Tag: math.MaxUint32, NoErr: true, Str2:"4294967294", Equ: false},
		TagInfo{Str: "4294967296", Tag: math.MaxUint32, NoErr: false, Str2:"", Equ: false},
		TagInfo{Str: "-1", Tag: 1, NoErr: false, Str2:"", Equ: false},
	}

	for _, tag := range tags {
		a, err := BuildTag(tag.Str)
		if tag.NoErr != (err == nil) {
			t.Fatalf("need err =%v, real=%v", tag.NoErr, err)
		}

		if err != nil {
			continue
		} else if err == nil {
			if a != nil && *a != tag.Tag {
				t.Fatalf("%d!=%d", *a, tag.Tag)
			}
		}

		b := ParseTag(a)
		if tag.Equ != (tag.Str2 == b) {
			t.Fatalf("need str=%s, real=%v", tag.Str2, b)
		} else {
			continue
		}
	}

}
