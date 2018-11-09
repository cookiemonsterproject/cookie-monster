package cookiejar_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/cookiejars/cookiejar"
	"github.com/cookiejars/cookiejar/mock"
)

func TestCookieJar(t *testing.T) {
	t.Parallel()

	expectedContent := "test content"
	mockCookie := &mock.Cookie{
		IDFn: func() string {
			return "test-cookie"
		},
		ContentFn: func() (interface{}, error) {
			return expectedContent, nil
		},
	}

	var cookies atomic.Value
	cookies.Store([]cookiejar.Cookie{mockCookie})
	mockJar := &mock.Jar{
		RetrieveFn: func() ([]cookiejar.Cookie, error) {
			return cookies.Load().([]cookiejar.Cookie), nil
		},
		RetireFn: func(cookie cookiejar.Cookie) error {
			cookies.Store([]cookiejar.Cookie{})
			return nil
		},
	}
	mockBackoff := &mock.Backoff{
		NextFn: func() {},
		CurrentFn: func() time.Duration {
			return time.Millisecond
		},
		ResetFn: func() {},
	}

	digestFn := func(cookie cookiejar.Cookie) error {
		got, err := cookie.Content()
		assertNil(t, err)

		assertEqual(t, expectedContent, got.(string))

		return nil
	}

	d := cookiejar.NewDigester(mockJar, cookiejar.SetWorkers(1), cookiejar.SetBackoff(mockBackoff))
	err := d.Start(digestFn)
	assertNil(t, err)

	time.Sleep(2 * time.Millisecond)
	d.Stop()

	assertTrue(t, mockCookie.IDInvoked)
	assertTrue(t, mockCookie.ContentInvoked)
	assertTrue(t, mockJar.RetrieveInvoked)
	assertTrue(t, mockJar.RetireInvoked)
	assertTrue(t, mockBackoff.NextInvoked)
	assertTrue(t, mockBackoff.CurrentInvoked)
	assertTrue(t, mockBackoff.ResetInvoked)
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
