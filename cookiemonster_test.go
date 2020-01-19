package cookiemonster_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/cookiejars/cookiemonster"
	"github.com/cookiejars/cookiemonster/mock"
)

func TestCookieJar(t *testing.T) {
	t.Parallel()

	expectedContent := "test content"
	mockCookie := &mock.Cookie{
		IDFn: func() string {
			return "test-cookie"
		},
		ContentFn: func() interface{} {
			return expectedContent
		},
	}

	var cookies atomic.Value
	cookies.Store([]cookiemonster.Cookie{mockCookie})
	mockJar := &mock.Jar{
		RetrieveFn: func() ([]cookiemonster.Cookie, error) {
			return cookies.Load().([]cookiemonster.Cookie), nil
		},
		RetireFn: func(cookie cookiemonster.Cookie) error {
			cookies.Store([]cookiemonster.Cookie{})
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

	digestFn := func(cookie cookiemonster.Cookie) error {
		got := cookie.Content()
		assertEqual(t, expectedContent, got.(string))

		return nil
	}

	d := cookiemonster.NewDigester(mockJar, cookiemonster.SetWorkers(1), cookiemonster.SetBackoff(mockBackoff))
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
