package blocks

import (
	"sync"
	"sync/atomic"
)

// AttestationVerifyStats tracks the success and failures of attestation verifications.
type AttestationVerifyStats struct {
	// Successes keeps track of correctly validated attestations count.
	Successes atomic.Uint64
	// Failures keeps track of failed validations mapping error strings to atomic counters.
	Failures sync.Map
}

// AttestationStats holds the tracking stats for processed attestations.
var AttestationStats = &AttestationVerifyStats{}

// RecordSuccess atomically increments the successes counter.
func (a *AttestationVerifyStats) RecordSuccess() {
	a.Successes.Add(1)
}

// RecordFailure atomically registers a failure with the given string error message as context.
func (a *AttestationVerifyStats) RecordFailure(err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	counterItem, _ := a.Failures.LoadOrStore(msg, &atomic.Uint64{})
	counter := counterItem.(*atomic.Uint64)
	counter.Add(1)
}

// SummaryAndReset extracts all tracked successes and failures into a simple summary map
// and resets the struct.
// The map will always contain the key "success" mapped to the number of successful attestations.
// Other keys correspond to failure error strings.
func (a *AttestationVerifyStats) SummaryAndReset() map[string]uint64 {
	summary := make(map[string]uint64)
	summary["success"] = a.Successes.Swap(0)

	a.Failures.Range(func(key, value any) bool {
		msg := key.(string)
		ptr := value.(*atomic.Uint64)
		count := ptr.Load()
		if count > 0 {
			summary[msg] = count
		}
		// Delete from the original.
		a.Failures.Delete(msg)
		return true
	})

	return summary
}
