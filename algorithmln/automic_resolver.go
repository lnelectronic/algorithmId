// ---------------------------------------------------------------------------
// LN-ELECTRONIC PROJECT LN-16C10R
// wwww.ln-electronic.com  ProjectManager : @Kimera
// FileData: 8/10/2564 7:33 2564  FileName : automic_resolver.go
// ---------------------------------------------------------------------------

package algorithmln

import "sync/atomic"

var lastTime int64
var lastSeq uint32

// AtomicResolver define as atomic sequence resolver, base on standard sync/atomic.
func AtomicResolver(ms int64) (uint16, error) {
	var last int64
	var seq, localSeq uint32

	for {
		last = atomic.LoadInt64(&lastTime)
		localSeq = atomic.LoadUint32(&lastSeq)
		if last > ms {
			return MaxSequence, nil
		}

		if last == ms {
			seq = MaxSequence & (localSeq + 1)
			if seq == 0 {
				return MaxSequence, nil
			}
		}

		if atomic.CompareAndSwapInt64(&lastTime, last, ms) && atomic.CompareAndSwapUint32(&lastSeq, localSeq, seq) {
			return uint16(seq), nil
		}
	}
}
