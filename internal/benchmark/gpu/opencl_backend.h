#ifndef OPENCL_BACKEND_H
#define OPENCL_BACKEND_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    double h2d_latency_us;
    double d2h_latency_us;
    double h2d_bw_gbps;
    double d2h_bw_gbps;
    double kernel_latency_us;
    int    ok;
} GPUBenchResult;

int openclInit(void);
void openclShutdown(void);
GPUBenchResult openclBenchmark(int bufferSize, int iterations);

#ifdef __cplusplus
}
#endif

#endif
