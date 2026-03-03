package blocks

import (
	"errors"
	"sync"
	"testing"

	"github.com/OffchainLabs/prysm/v7/testing/assert"
)

func TestAttestationStats_SuccessTally(t *testing.T) {
	stats := &AttestationVerifyStats{}

	var wg sync.WaitGroup
	workers := 100
	for range workers {
		wg.Go(func() {
			stats.RecordSuccess()
		})
	}

	wg.Wait()
	summary := stats.SummaryAndReset()

	assert.Equal(t, uint64(100), summary["success"], "expected 100 successful tallies")

	// Ensure reset happened
	summaryAfterReset := stats.SummaryAndReset()
	assert.Equal(t, uint64(0), summaryAfterReset["success"])
}

func TestAttestationStats_FailureTallyAndReasons(t *testing.T) {
	stats := &AttestationVerifyStats{}
	var wg sync.WaitGroup
	errEpoch := errors.New("invalid epoch")
	errSig := errors.New("bad signature")

	// 50 invalid epoch errors
	epochWorkers := 50
	for range epochWorkers {
		wg.Go(func() {
			stats.RecordFailure(errEpoch)
		})
	}

	// 30 bad signature errors
	sigWorkers := 30
	for range sigWorkers {
		wg.Go(func() {
			stats.RecordFailure(errSig)
		})
	}

	// Record a nil error (should be ignored)
	stats.RecordFailure(nil)

	wg.Wait()
	summary := stats.SummaryAndReset()

	assert.Equal(t, uint64(0), summary["success"], "expected 0 successful tallies")
	assert.Equal(t, uint64(50), summary["invalid epoch"], "expected 50 invalid epoch tallies")
	assert.Equal(t, uint64(30), summary["bad signature"], "expected 30 bad signature tallies")

	// Ensure reset happened and the map was cleared out
	summaryAfterReset := stats.SummaryAndReset()
	assert.Equal(t, uint64(0), summaryAfterReset["invalid epoch"])
	assert.Equal(t, uint64(0), summaryAfterReset["bad signature"])
}
