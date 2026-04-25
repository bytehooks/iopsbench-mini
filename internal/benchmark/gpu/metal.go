//go:build darwin && cgo

package gpu

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Metal -framework Foundation
#import <Metal/Metal.h>
#include <mach/mach_time.h>
#include <string.h>

typedef struct {
    double h2d_latency_us;
    double d2h_latency_us;
    double h2d_bw_gbps;
    double d2h_bw_gbps;
    double kernel_latency_us;
    int    ok;
} GPUBenchResult;

static id<MTLDevice> device = nil;
static id<MTLCommandQueue> queue = nil;

int metalInit(void) {
    device = MTLCreateSystemDefaultDevice();
    if (!device) return -1;
    queue = [device newCommandQueue];
    if (!queue) return -1;
    return 0;
}

void metalShutdown(void) {
    queue = nil;
    device = nil;
}

static uint64_t nsElapsed(uint64_t start, uint64_t end) {
    static mach_timebase_info_data_t info = {0};
    if (info.denom == 0) mach_timebase_info(&info);
    return (end - start) * info.numer / info.denom;
}

static void commitAndWait(id<MTLCommandBuffer> cb) {
    dispatch_semaphore_t sem = dispatch_semaphore_create(0);
    [cb addCompletedHandler:^(id<MTLCommandBuffer> _Nonnull buf) {
        (void)buf;
        dispatch_semaphore_signal(sem);
    }];
    [cb commit];
    dispatch_semaphore_wait(sem, DISPATCH_TIME_FOREVER);
}

GPUBenchResult metalBenchmark(int bufferSize, int iterations) {
    GPUBenchResult res = {0};
    res.ok = 0;

    if (!device || !queue) return res;

    id<MTLBuffer> gpuBuf = [device newBufferWithLength:bufferSize
        options:MTLResourceStorageModePrivate];
    id<MTLBuffer> sharedBuf = [device newBufferWithLength:bufferSize
        options:MTLResourceStorageModeShared];

    if (!gpuBuf || !sharedBuf) return res;

    uint8_t *hostPtr = (uint8_t *)[sharedBuf contents];
    for (int i = 0; i < bufferSize; i++) hostPtr[i] = (uint8_t)(i & 0xFF);

    for (int i = 0; i < 3; i++) {
        id<MTLCommandBuffer> cb = [queue commandBuffer];
        id<MTLBlitCommandEncoder> blit = [cb blitCommandEncoder];
        [blit copyFromBuffer:sharedBuf sourceOffset:0
                  toBuffer:gpuBuf destinationOffset:0 size:bufferSize];
        [blit endEncoding];
        commitAndWait(cb);

        cb = [queue commandBuffer];
        blit = [cb blitCommandEncoder];
        [blit copyFromBuffer:gpuBuf sourceOffset:0
                  toBuffer:sharedBuf destinationOffset:0 size:bufferSize];
        [blit endEncoding];
        commitAndWait(cb);
    }

    uint64_t h2dStart = mach_absolute_time();
    for (int i = 0; i < iterations; i++) {
        id<MTLCommandBuffer> cb = [queue commandBuffer];
        id<MTLBlitCommandEncoder> blit = [cb blitCommandEncoder];
        [blit copyFromBuffer:sharedBuf sourceOffset:0
                  toBuffer:gpuBuf destinationOffset:0 size:bufferSize];
        [blit endEncoding];
        commitAndWait(cb);
    }
    uint64_t h2dEnd = mach_absolute_time();

    uint64_t d2hStart = mach_absolute_time();
    for (int i = 0; i < iterations; i++) {
        id<MTLCommandBuffer> cb = [queue commandBuffer];
        id<MTLBlitCommandEncoder> blit = [cb blitCommandEncoder];
        [blit copyFromBuffer:gpuBuf sourceOffset:0
                  toBuffer:sharedBuf destinationOffset:0 size:bufferSize];
        [blit endEncoding];
        commitAndWait(cb);
    }
    uint64_t d2hEnd = mach_absolute_time();

    NSString *kernelSrc = @"kernel void k(device uint *buf) { buf[0] += 1; }";
    NSError *err = nil;
    id<MTLLibrary> lib = [device newLibraryWithSource:kernelSrc
        options:nil error:&err];
    if (lib) {
        id<MTLFunction> fn = [lib newFunctionWithName:@"k"];
        id<MTLComputePipelineState> pso = [device newComputePipelineStateWithFunction:fn error:nil];
        id<MTLBuffer> smallBuf = [device newBufferWithLength:4
            options:MTLResourceStorageModeShared];

        uint64_t kStart = mach_absolute_time();
        for (int i = 0; i < iterations; i++) {
            id<MTLCommandBuffer> cb = [queue commandBuffer];
            id<MTLComputeCommandEncoder> enc = [cb computeCommandEncoder];
            [enc setComputePipelineState:pso];
            [enc setBuffer:smallBuf offset:0 atIndex:0];
            [enc dispatchThreadgroups:MTLSizeMake(1,1,1)
                threadsPerThreadgroup:MTLSizeMake(1,1,1)];
            [enc endEncoding];
            commitAndWait(cb);
        }
        uint64_t kEnd = mach_absolute_time();
        res.kernel_latency_us = (double)nsElapsed(kStart, kEnd) / (double)iterations / 1000.0;
    }

    uint64_t h2dNs = nsElapsed(h2dStart, h2dEnd);
    uint64_t d2hNs = nsElapsed(d2hStart, d2hEnd);

    res.h2d_latency_us = (double)h2dNs / (double)iterations / 1000.0;
    res.d2h_latency_us = (double)d2hNs / (double)iterations / 1000.0;

    double bytes = (double)bufferSize * (double)iterations;
    res.h2d_bw_gbps = (bytes / (h2dNs / 1e9)) / 1e9;
    res.d2h_bw_gbps = (bytes / (d2hNs / 1e9)) / 1e9;
    res.ok = 1;

    return res;
}
*/
import "C"

// New returns the Metal-backed GPU benchmarker for macOS.
func New() Benchmarker {
	return &metalBenchmarker{}
}

func newMetalBenchmarker() Benchmarker {
	return &metalBenchmarker{}
}

type metalBenchmarker struct{}

func (m *metalBenchmarker) Init() error {
	if C.metalInit() != 0 {
		return ErrNoDevice
	}
	return nil
}

func (m *metalBenchmarker) Shutdown() {
	C.metalShutdown()
}

func (m *metalBenchmarker) BackendName() string {
	return "Metal"
}

func (m *metalBenchmarker) Benchmark(bufferSize, iterations int) (Result, error) {
	res := C.metalBenchmark(C.int(bufferSize), C.int(iterations))
	if res.ok == 0 {
		return Result{}, ErrBenchmarkFailed
	}
	return Result{
		BufferSize:      int64(bufferSize),
		H2DLatencyUS:    float64(res.h2d_latency_us),
		D2HLatencyUS:    float64(res.d2h_latency_us),
		H2DBandwidthGBs: float64(res.h2d_bw_gbps),
		D2HBandwidthGBs: float64(res.d2h_bw_gbps),
		KernelLatencyUS: float64(res.kernel_latency_us),
	}, nil
}
