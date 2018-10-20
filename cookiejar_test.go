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

	var deleted atomic.Value
	deleted.Store(false)

	expectedContent := "test content"
	mockCookie := &mock.Cookie{
		ContentFn: func() ([]byte, error) {
			return []byte(expectedContent), nil
		},
		DeleteFn: func() error {
			deleted.Store(true)
			return nil
		},
	}
	mockJar := &mock.Jar{
		FetchFn: func() ([]cookiejar.Cookie, error) {
			if deleted.Load().(bool) {
				return []cookiejar.Cookie{}, nil
			}
			return []cookiejar.Cookie{mockCookie}, nil
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

		assertEqual(t, expectedContent, string(got))

		return nil
	}

	d := cookiejar.NewDigester(1, mockJar, mockBackoff)
	err := d.Start(digestFn)
	assertNil(t, err)

	time.Sleep(2 * time.Millisecond)
	d.Stop()

	assertTrue(t, mockCookie.ContentInvoked)
	assertTrue(t, mockCookie.DeleteInvoked)
	assertTrue(t, mockJar.FetchInvoked)
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
