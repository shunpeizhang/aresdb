package query

import "github.com/uber/aresdb/memutils"

// NonAggrBatchExecutorImpl is batch executor implementation for non-aggregation query
type NonAggrBatchExecutorImpl struct {
	*BatchExecutorImpl
}

// NonAggrBatchExecutorImpl has empty reduce operation
func (e *NonAggrBatchExecutorImpl) reduce() {
	// nothing to do for non-aggregation
}

// project for non-aggregation query will figure out values for all selected columns,
// the calculation will be in CPU instead of GPU to avoid many vp movement while the result size will be always limitted by LIMIT
func (e *NonAggrBatchExecutorImpl) project() {
	if e.qc.OOPK.currentBatch.size == 0 {
		return
	}
	//allocate index vector in memory
	size := e.qc.OOPK.currentBatch.size * 4
	hostIndexH := memutils.HostAlloc(size)
	// Copy index vector from Device
	memutils.AsyncCopyDeviceToHost(hostIndexH, e.qc.OOPK.currentBatch.indexVectorD.getPointer(), size, e.stream, e.qc.Device)
	memutils.WaitForCudaStream(e.stream, e.qc.Device)

	// fill measure values row by row


	//release index space
	memutils.HostFree(hostIndexH)
}
