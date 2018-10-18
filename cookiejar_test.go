package cookiejar_test

import (
	"testing"
	"time"

	"github.com/cookiejars/cookiejar"
	"github.com/cookiejars/cookiejar/mock"
)

func TestCookieJar_Start(t *testing.T) {
	t.Parallel()

	var fetchResult []cookiejar.Cookie
	expectedContent := "test content"
	mockCookie := &mock.Cookie{
		ContentFn: func() ([]byte, error) {
			return []byte(expectedContent), nil
		},
		DeleteFn: func() error {
			fetchResult = []cookiejar.Cookie{}

			return nil
		},
	}
	fetchResult = []cookiejar.Cookie{mockCookie}
	mockJar := &mock.Jar{
		FetchFn: func() ([]cookiejar.Cookie, error) {
			return fetchResult, nil
		},
	}

	digestFn := func(cookie cookiejar.Cookie) error {
		got, err := cookie.Content()
		assertNil(t, err)

		assertEqual(t, expectedContent, string(got))

		return nil
	}

	d := cookiejar.NewDigester(1, mockJar, cookiejar.ConstantBackoff(time.Millisecond))
	err := d.Start(digestFn)
	assertNil(t, err)

	time.Sleep(2 * time.Millisecond)
	d.Stop()

	assertTrue(t, mockCookie.ContentInvoked)
	assertTrue(t, mockCookie.DeleteInvoked)
	assertTrue(t, mockJar.FetchInvoked)
}

func assertEqual(t *testing.T, expected, got interface{}) {
	if expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func assertTrue(t *testing.T, got interface{}) {
	assertEqual(t, true, got)
}

func assertNil(t *testing.T, got interface{}) {
	assertEqual(t, nil, got)
}
