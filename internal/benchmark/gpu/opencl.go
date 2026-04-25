//go:build linux && cgo

package gpu

/*
#cgo LDFLAGS: -lOpenCL
#include <CL/cl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

static cl_platform_id platform = NULL;
static cl_device_id device = NULL;
static cl_context context = NULL;
static cl_command_queue queue = NULL;

static uint64_t nsNow(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (uint64_t)ts.tv_sec * 1000000000ULL + (uint64_t)ts.tv_nsec;
}

int openclInit(void) {
    cl_uint numPlatforms = 0;
    if (clGetPlatformIDs(0, NULL, &numPlatforms) != CL_SUCCESS || numPlatforms == 0)
        return -1;

    cl_platform_id *platforms = (cl_platform_id *)malloc(sizeof(cl_platform_id) * numPlatforms);
    if (!platforms) return -1;
    if (clGetPlatformIDs(numPlatforms, platforms, NULL) != CL_SUCCESS) {
        free(platforms);
        return -1;
    }

    for (cl_uint i = 0; i < numPlatforms; i++) {
        cl_uint numDevices = 0;
        if (clGetDeviceIDs(platforms[i], CL_DEVICE_TYPE_GPU, 0, NULL, &numDevices) == CL_SUCCESS && numDevices > 0) {
            platform = platforms[i];
            break;
        }
    }
    if (!platform) platform = platforms[0];
    free(platforms);

    cl_int err;
    err = clGetDeviceIDs(platform, CL_DEVICE_TYPE_GPU, 1, &device, NULL);
    if (err != CL_SUCCESS) {
        err = clGetDeviceIDs(platform, CL_DEVICE_TYPE_ALL, 1, &device, NULL);
    }
    if (err != CL_SUCCESS) return -1;

    context = clCreateContext(NULL, 1, &device, NULL, NULL, &err);
    if (err != CL_SUCCESS) return -1;

#if CL_VERSION_2_0
    queue = clCreateCommandQueueWithProperties(context, device, 0, &err);
#else
    queue = clCreateCommandQueue(context, device, 0, &err);
#endif
    if (err != CL_SUCCESS) return -1;

    return 0;
}

void openclShutdown(void) {
    if (queue) { clReleaseCommandQueue(queue); queue = NULL; }
    if (context) { clReleaseContext(context); context = NULL; }
    device = NULL;
    platform = NULL;
}

static cl_int copyH2D(cl_mem dst, const void *src, size_t size) {
    return clEnqueueWriteBuffer(queue, dst, CL_TRUE, 0, size, src, 0, NULL, NULL);
}

static cl_int copyD2H(void *dst, cl_mem src, size_t size) {
    return clEnqueueReadBuffer(queue, src, CL_TRUE, 0, size, dst, 0, NULL, NULL);
}

typedef struct {
    double h2d_latency_us;
    double d2h_latency_us;
    double h2d_bw_gbps;
    double d2h_bw_gbps;
    double kernel_latency_us;
    int    ok;
} GPUBenchResult;

GPUBenchResult openclBenchmark(int bufferSize, int iterations) {
    GPUBenchResult res = {0};
    res.ok = 0;
    if (!context || !queue || !device) return res;

    cl_int err;
    cl_mem gpuBuf = clCreateBuffer(context, CL_MEM_READ_WRITE, bufferSize, NULL, &err);
    if (err != CL_SUCCESS) return res;

    void *hostBuf = malloc(bufferSize);
    if (!hostBuf) { clReleaseMemObject(gpuBuf); return res; }
    for (int i = 0; i < bufferSize; i++) ((uint8_t *)hostBuf)[i] = (uint8_t)(i & 0xFF);

    for (int i = 0; i < 3; i++) {
        copyH2D(gpuBuf, hostBuf, bufferSize);
        copyD2H(hostBuf, gpuBuf, bufferSize);
    }

    uint64_t h2dStart = nsNow();
    for (int i = 0; i < iterations; i++) {
        if (copyH2D(gpuBuf, hostBuf, bufferSize) != CL_SUCCESS) break;
    }
    uint64_t h2dEnd = nsNow();

    uint64_t d2hStart = nsNow();
    for (int i = 0; i < iterations; i++) {
        if (copyD2H(hostBuf, gpuBuf, bufferSize) != CL_SUCCESS) break;
    }
    uint64_t d2hEnd = nsNow();

    const char *src = "__kernel void k(__global uint *buf) { buf[0] += 1; }";
    cl_program prog = clCreateProgramWithSource(context, 1, &src, NULL, &err);
    if (prog && clBuildProgram(prog, 1, &device, "", NULL, NULL) == CL_SUCCESS) {
        cl_kernel kern = clCreateKernel(prog, "k", &err);
        cl_mem smallBuf = clCreateBuffer(context, CL_MEM_READ_WRITE, 4, NULL, &err);
        uint32_t zero = 0;
        clEnqueueWriteBuffer(queue, smallBuf, CL_TRUE, 0, 4, &zero, 0, NULL, NULL);
        clSetKernelArg(kern, 0, sizeof(cl_mem), &smallBuf);

        size_t global = 1;
        uint64_t kStart = nsNow();
        for (int i = 0; i < iterations; i++) {
            clEnqueueNDRangeKernel(queue, kern, 1, NULL, &global, &global, 0, NULL, NULL);
            clFinish(queue);
        }
        uint64_t kEnd = nsNow();
        res.kernel_latency_us = (double)(kEnd - kStart) / (double)iterations / 1000.0;

        clReleaseMemObject(smallBuf);
        clReleaseKernel(kern);
    }
    if (prog) clReleaseProgram(prog);

    uint64_t h2dNs = h2dEnd - h2dStart;
    uint64_t d2hNs = d2hEnd - d2hStart;

    res.h2d_latency_us = (double)h2dNs / (double)iterations / 1000.0;
    res.d2h_latency_us = (double)d2hNs / (double)iterations / 1000.0;

    double bytes = (double)bufferSize * (double)iterations;
    res.h2d_bw_gbps = (bytes / (h2dNs / 1e9)) / 1e9;
    res.d2h_bw_gbps = (bytes / (d2hNs / 1e9)) / 1e9;
    res.ok = 1;

    free(hostBuf);
    clReleaseMemObject(gpuBuf);
    return res;
}
*/
import "C"

// New returns the OpenCL-backed GPU benchmarker for Linux.
func New() Benchmarker {
	return &openclBenchmarker{}
}

func newOpenCLBenchmarker() Benchmarker {
	return &openclBenchmarker{}
}

type openclBenchmarker struct{}

func (o *openclBenchmarker) Init() error {
	if C.openclInit() != 0 {
		return ErrNoDevice
	}
	return nil
}

func (o *openclBenchmarker) Shutdown() {
	C.openclShutdown()
}

func (o *openclBenchmarker) BackendName() string {
	return "OpenCL"
}

func (o *openclBenchmarker) Benchmark(bufferSize, iterations int) (Result, error) {
	res := C.openclBenchmark(C.int(bufferSize), C.int(iterations))
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
