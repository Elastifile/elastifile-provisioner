package retry

import (
	"fmt"
	"math"
	"time"

	"github.com/elastifile/errors"

	log "gopkg.in/inconshreveable/log15.v2"
)

var Log = log.New()

func init() {
	Log.SetHandler(log.DiscardHandler())
}

////////////////////////////////////////////////////////////////////////////////

// Based on the idea from this post: http://dave.cheney.net/2014/12/24/inspecting-errors
type TimeoutError struct {
	timeout time.Duration
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("action timed out after %v", e.timeout)
}

func (e *TimeoutError) Timeout() bool {
	return true
}

type TemporaryError struct {
	Err error
}

func (e *TemporaryError) Error() string {
	return fmt.Sprintf("temporary error: %v", e.Err)
}

func (e *TemporaryError) Temporary() bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////

type retrier interface {
	nextTimeout() time.Duration
	totalTimeout() time.Duration
	keepTrying() bool
}

func totalTimeout(r retrier) (result time.Duration) {
	for r.keepTrying() {
		// One millisecond is chosen to cover for discrepancy in timers
		result += r.nextTimeout() + time.Millisecond
	}

	return result
}

func retryDo(r retrier, action func() error) error {
	var (
		result error
		delta  time.Duration
	)

	done := make(chan error)
	total := r.totalTimeout()
	timeoutPassed := time.After(total)

	for r.keepTrying() {
		projected := time.Now().Add(r.nextTimeout())

		go func() { done <- action() }()

		select {
		case result = <-done:
		case <-timeoutPassed:
			result = &TimeoutError{total}
		}

		if result == nil || !IsTemporary(result) {
			break
		}

		actual := time.Now()

		if actual.Add(delta).Before(projected) {
			time.Sleep(projected.Sub(actual.Add(delta)))
			delta = 0
		} else {
			delta += actual.Sub(projected)
		}
	}

	if result != nil {
		if real, ok := result.(*TemporaryError); ok {
			result = real.Err
		}
		// In case someone wrapped nil in TemporaryError, we don't
		// want to fail because of it.
		if result != nil {
			result = errors.New(result)
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////

// Common code, to be embedded by specific retriers.

type retrierCommon struct {
	timeout time.Duration
	retries int
}

func (cr *retrierCommon) keepTrying() bool {
	return cr.retries > 0
}

////////////////////////////////////////////////////////////////////////////////

type Basic struct {
	Timeout time.Duration
	Retries int
}

func (b Basic) Do(action func() error) error {
	br := &basicRetrier{}
	br.timeout = b.Timeout
	br.retries = b.Retries

	return retryDo(br, action)
}

type basicRetrier struct {
	retrierCommon
}

func (br *basicRetrier) nextTimeout() time.Duration {
	br.retries--
	return br.timeout
}

func (br *basicRetrier) totalTimeout() time.Duration {
	// Pass a copy to the called function, since it mutates the state.
	rCopy := *br
	return totalTimeout(&rCopy)
}

////////////////////////////////////////////////////////////////////////////////

type Linear struct {
	Timeout time.Duration
	Retries int
}

func (l Linear) Do(action func() error) error {
	lr := &linearRetrier{}
	lr.timeout = l.Timeout
	lr.retries = l.Retries

	return retryDo(lr, action)
}

type linearRetrier struct {
	retrierCommon
	step int
}

func (lr *linearRetrier) nextTimeout() time.Duration {
	lr.retries--
	lr.step++
	return lr.timeout * time.Duration(lr.step)
}

func (lr *linearRetrier) totalTimeout() time.Duration {
	// Pass a copy to the called function, since it mutates the state.
	rCopy := *lr
	return totalTimeout(&rCopy)
}

////////////////////////////////////////////////////////////////////////////////

type Sigmoid struct {
	Limit   time.Duration
	Retries int
}

func (s Sigmoid) Do(action func() error) error {
	sr := &sigmoidRetrier{
		step:  s.Retries / -2,
		upper: s.Limit,
	}
	sr.retries = s.Retries

	return retryDo(sr, action)
}

type sigmoidRetrier struct {
	retrierCommon
	step  int
	upper time.Duration
}

func (sr *sigmoidRetrier) nextTimeout() time.Duration {
	sr.retries--
	sr.step++
	return time.Duration(2 / (1 + math.Pow(math.E, float64(-sr.step))) * float64(sr.upper))
}

func (sr *sigmoidRetrier) totalTimeout() time.Duration {
	// Pass a copy to the called function, since it mutates the state.
	rCopy := *sr
	result := totalTimeout(&rCopy)
	return result
}

////////////////////////////////////////////////////////////////////////////////

type Root struct {
	Increment time.Duration
	Retries   int
}

func (r Root) Do(action func() error) error {
	rr := &rootRetrier{
		increment: r.Increment,
	}
	rr.retries = r.Retries

	return retryDo(rr, action)
}

type rootRetrier struct {
	retrierCommon
	step      int
	increment time.Duration
}

func (rr *rootRetrier) nextTimeout() time.Duration {
	rr.retries--
	rr.step++
	return time.Duration(math.Sqrt(float64(rr.step)) * float64(rr.increment))
}

func (rr *rootRetrier) totalTimeout() time.Duration {
	// Pass a copy to the called function, since it mutates the state.
	rCopy := *rr
	return totalTimeout(&rCopy)
}

////////////////////////////////////////////////////////////////////////////////

type Opportunistic struct {
	Total   time.Duration
	Retries int
}

func (o Opportunistic) Do(action func() error) error {
	opr := &opportunisticRetrier{
		total: o.Total,
	}
	opr.retries = o.Retries

	return retryDo(opr, action)
}

type opportunisticRetrier struct {
	retrierCommon
	total time.Duration
}

func (opr *opportunisticRetrier) nextTimeout() time.Duration {
	opr.retries--
	return 0
}

func (opr *opportunisticRetrier) totalTimeout() time.Duration {
	return opr.total
}

////////////////////////////////////////////////////////////////////////////////

// Do retries an operation until it completes without error or times out
func Do(timeout time.Duration, action func() error) error {
	return Sigmoid{
		Limit:   timeout / 10,
		Retries: 10,
	}.Do(action)
}

func IsTemporary(err error) (result bool) {
	type temporary interface {
		Temporary() bool
	}
	switch e := err.(type) {
	case temporary:
		result = e.Temporary()
	case *errors.Error:
		result = IsTemporary(e.Err)
	}
	Log.Debug("IsTemporary", "result", result, "err", err)
	return result
}

func IsTimeout(err error) (result bool) {
	type timeout interface {
		Timeout() bool
	}
	switch e := err.(type) {
	case timeout:
		result = e.Timeout()
	case *errors.Error:
		result = IsTimeout(e.Err)
	}
	Log.Debug("IsTimeout", "result", result, "err", err)
	return result
}
