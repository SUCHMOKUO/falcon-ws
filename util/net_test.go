package util

import "testing"

func TestIsDomain(t *testing.T) {
	cases := []struct {
		domain string
		res bool
	} {
		{ "www.baidu.com", true },
		{ "google.com", true },
		{ "abc123.tk", true },
		{ "1.2.3.4", false },
		{ "8.8.8.8", false },
	}

	for _, test := range cases {
		if res := IsDomain(test.domain); res != test.res {
			t.Errorf("%v, %v, %v\n", test.domain, test.res, res)
		}
	}
}

func TestIsValidHost(t *testing.T) {
	cases := []struct {
		host string
		res bool
	} {
		{ "www.baidu.com", true },
		{ "google.com", true },
		{ "abc123.tk", true },
		{ "1.2.3.4", true },
		{ "8.8.8.8", true },
		{ "qwe", false },
		{ "百度", false },
	}

	for _, test := range cases {
		if res := IsValidHost(test.host); res != test.res {
			t.Errorf("host: %v, expect: %v, but: %v\n", test.host, test.res, res)
		}
	}
}
