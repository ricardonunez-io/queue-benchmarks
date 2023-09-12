# Async/Synchronous Queue Benchmarks
## Javascript (Bun) vs. Go

With the recent release of Bun 1.0, there's been a lot of talk on Twitter about its incredibly impressive performance, as well as how it fares against faster languages like Go.

I'm currently building an application that's incredibly data intensive where each service deals with hundreds of thousands of datapoints per request, with arrays reaching up to the millions after serializing to flatbuffers, so I wanted to make a benchmark to see if Bun would fit my use case as many of the services are already built in Go, but I'm building the frontend soon, so if I can leverage Javascript for more of the backend, it'll make the rest of development less fragmented for my one person team (of me).

Many of the benchmarks I've seen online have been on IO/syscall intensive tasks where the program isn't CPU or memory bound, so I wanted to create a benchmark more representative of real-world use cases with data transformations among large data objects rather than syscall driven benchmarks like DB retrieval or JSON parsing with SIMD where Bun is heavily optimized but Go can leverage CGo to use SIMDJSON or another library to match it (in which case it's a battle of libraries, not languages).

The benchmarks themselves are two-fold for each language. 

---

### Async

Firstly, I wanted to test the async capabilities of both language when it comes to data processing in a queue.

> It's worth noting that for the sake of the benchmark, "async" doesn't mean parallel execution. For CPU-bound tasks, Javascript's event loop doesn't allow for parallel execution without the usage of worker threads, so instead, I simply used `Promise.all()` in JS (which still executes each promise serially) and goroutines & channels in Go (which also execute each concurrent task serially due to channels' locking nature). If I wanted to test true parallel execution, I would've used a SharedArrayBuffer in Javascript across many worker-threads and a similar channel structure or mutexes in Go, but then dealing with race conditions would've been less fun.

Both implementations were as close to identical as I could make them, and I tried to limit object creation in Javascript as much as possible (to avoid unnecessary/expensive garbage collection) so as to not give Go an unfair advantage. I used pre-allocation for the array to avoid expensive copies of objects, similar to how Go can instantiate a channel/array with a predefined limit.

As for the benchmark itself, we're going to be testing each language/runtime's ability to concurrently "process" a queue (array in JS/channel in Go) by having each async task/goroutine go through an array, assign a variable to the value of that task number's index in the array, increment the variable, and then replace that same index's value with the incremented variable.

The channel/array will be of length `n` where `n` is also equal to the number of tasks/goroutines we want to test.

---

### Synchronous

For the synchronous benchmark, the goal was simple. Write a for loop that accomplishes the same as above where each iteration over the loop retrieves the value of the array/channel at the index of the iteration's number, assigns it to a variable, icrements the variable, and reassigns the item back to the array/channel.

The channel/array will also be of length `n` where `n` is equal to the number of tasks/goroutines we want to test.

---

## Testing

To perform the testing, I used [Hyperfine](https://github.com/sharkdp/hyperfine) and [GNU Time](https://www.gnu.org/software/time/) to run multiple executions of each incrementing number of tasks/goroutines where each benchmark's execution is increased by 10x, starting at 100 tasks/goroutines and ending at 10,000,000.

### Execution

To execute the test, first, you have to clone the repo:

```zsh
git clone https://github.com/ricardonunez-io/queue-benchmarks/
```

Then, you have to build the Go binaries:

```zsh
cd queue-benchmarks/go/async && go build main.go
cd queue-benchmarks/go/sync && go build main.go
```

Finally, you can run the tests by `cd`ing back to the `queue-benchmarks` folder and executing `tests.sh` and `memory-tests.sh`, passing in only the argument of how many tasks/goroutines you'd like to run per test, i.e.:

```zsh
./tests.sh 1000
./memory-tests.sh 1000
```

## Results

### Resource Usage

I ran the tests on a 1 core, 2GB RAM server hosted on Linode's US-East to both limit any multithreading advantage from Go as well as represent a real world scenario for a realistic container's resources.

Firstly, resource usage. Speed might seem like an important test, but resource usage on a server tells you a lot about the constraints of a lanaguage, especially if there's bottlenecks that could translate to requests-per-second limits in production, heap memory errors, forced server restarts, etc.

For CPU usage, Bun competed incredibly well with Go. We'll get to why in a minute, but in terms of memory usage, Go's memory management and compiled nature makes it very difficult for any VM-bound runtime to compete, even one as performant as Bun's.

Here are the results from running GNU Time on each language's async/synchronous program after warming each file's execution up 100 times.

| Language   | Async   | Tasks/Goroutines   | CPU usage (%)   | Maximum Resident Set Size (Mb)   |
|:-----------|:--------|:-------------------|:----------------|:---------------------------------|
| Go         | Yes     | 100                | 91%             | 3.336                            |
| Javascript | Yes     | 100                | 85%             | 50.164                           |
| Go         | No      | 100                | 81%             | 3.352                            |
| Javascript | No      | 100                | 100%            | 47.656                           |
| Go         | Yes     | 1,000              | 100%            | 5.398                            |
| Javascript | Yes     | 1,000              | 74%             | 61.512                           |
| Go         | No      | 1,000              | 100%            | 3.352                            |
| Javascript | No      | 1,000              | 90%             | 50.211                           |
| Go         | Yes     | 10,000             | 88%             | 21.855                           |
| Javascript | Yes     | 10,000             | 52%             | 65.277                           |
| Go         | No      | 10,000             | 83%             | 3.355                            |
| Javascript | No      | 10,000             | 100%            | 54.715                           |
| Go         | Yes     | 100,000            | 100%            | 87.863                           |
| Javascript | Yes     | 100,000            | 41%             | 112.918                          |
| Go         | No      | 100,000            | 86%             | 3.367                            |
| Javascript | No      | 100,000            | 90%             | 58.684                           |
| Go         | Yes     | 1,000,000          | 99%             | 180.340                          |
| Javascript | Yes     | 1,000,000          | 38%             | 429.168                          |
| Go         | No      | 1,000,000          | 100%            | 11.488                           |
| Javascript | No      | 1,000,000          | 92%             | 77.320                           |
| Go         | Yes     | 10,000,000         | 99%             | 359.570                          |
| Javascript | Yes     | 10,000,000         | 37%             | 1,807.570                        |
| Go         | No      | 10,000,000         | 93%             | 80.656                           |
| Javascript | No      | 10,000,000         | 90%             | 139.293                          |

As you can see, for each iteration of the async programs, memory usage is dramatically higher than the synchronous equivalent. However, the most important thing to note here is that with the exception of the program ran with 100,000 async tasks where Bun's memory usage is only about 28% higher, Go's memory usage is usually about 2-3x lower than Bun's at 10,000 tasks and above, with the lower iterations being between 12x-16x less memory intensive.

This is a dramatic difference, and running the tests multiple times gave very similar results each time, though I can't say what the statistical significance was on each iteration, especially that of the iteration with 100,000 goroutines ran on `go/async/main`, as that iteration specifically seemed to be Bun's sweet spot and Go's weak spot.

### Speed

Now, for the speed tests. Here's where the CPU usage really begins to explain itself. Given that these programs are, for the most part, CPU-bound due to the fact they are performing computations (incrementing an item that's assigned from an array and then placing it back in the array) more so than waiting on syscalls or performing IO bound tasks, Bun's lower CPU usage can be explained by one thing: speed.

Although it does fair well against Go in terms of speed for IO/syscall based operations like running a web server where JSON serialization of database objects might be the most expensive operation ([Elysia is a great example of how performant this can be in Javascript](https://elysiajs.com)), Go's speed is still much farther ahead of Bun's when it comes to processing instructions efficiently and quickly, which means Go's CPU usage % will also be higher due to the fact that it's simply processing more instructions quicker.

However, Bun's speed is impressive when it hits the sweet spot of 100k to 1M async tasks, where it's actually able to compete with Go, having only a 20% slowdown in speed at 100k tasks and only a ~2% slowdown at 1M tasks, *within a standard deviation* of Go's benchmark.

At 10 million, the difference widens again in Go's favor, but having 1 million promises run concurrently on a 1 million item array in Javascript and being at essentially the same speed as Go using goroutines is incredibly impressive, and for those working on workloads that aren't as memory constrained (i.e. longer running servers or larger containers), this speed is going to make a huge difference, along with Bun's very fast startup times which are going to make an even bigger difference for those running on serverless like Lambda or Vercel.

However, because the memory usage is still quite high compared to Go, you still might run into `heap out of memory` errors if running programs that process this much data.

| Language   | Async | Tasks/Goroutines | Mean Time | Std. Dev. | User Time | System Time | Range (min) | Range (max) | Runs |
|:-----------|:------|:-----------------|:----------|:----------|:----------|:------------|:------------|:------------|:-----|
| Go         | Yes   |              100 | 0.6 ms    | 0.1 ms    | 0.5 ms    | 0.1 ms      | 0.5 ms      | 3.1 ms      | 4161 |
| Javascript | Yes   |              100 | 13.1 ms   | 0.8 ms    | 6.8 ms    | 6.0 ms      | 11.3 ms     | 16.5 ms     | 213  |
| Go         | No    |              100 | 0.6 ms    | 0.1 ms    | 0.5 ms    | 0.1 ms      | 0.5 ms      | 2.6 ms      | 4262 |
| Javascript | No    |              100 | 11.3 ms   | 0.7 ms    | 5.9 ms    | 5.2 ms      | 8.8 ms      | 13.1 ms     | 291  |
| Go         | Yes   |            1,000 | 1.1 ms    | 0.1 ms    | 0.9 ms    | 0.1 ms      | 0.9 ms      | 3.3 ms      | 2421 |
| Javascript | Yes   |            1,000 | 22.3 ms   | 1.5 ms    | 15.9 ms   | 6.1 ms      | 17.4 ms     | 27.3 ms     | 152  |
| Go         | No    |            1,000 | 0.6 ms    | 0.1 ms    | 0.5 ms    | 0.1 ms      | 0.5 ms      | 2.0 ms      | 4493 |
| Javascript | No    |            1,000 | 11.5 ms   | 0.6 ms    | 5.7 ms    | 5.6 ms      | 10.1 ms     | 13.9 ms     | 269  |
| Go         | Yes   |           10,000 | 7.9 ms    | 0.4 ms    | 6.4 ms    | 1.4 ms      | 7.5 ms      | 10.6 ms     | 369  |
| Javascript | Yes   |           10,000 | 40.1 ms   | 1.5 ms    | 30.0 ms   | 9.8 ms      | 36.7 ms     | 42.9 ms     | 72   |
| Go         | No    |           10,000 | 0.6 ms    | 0.0 ms    | 0.5 ms    | 0.1 ms      | 0.5 ms      | 1.0 ms      | 3985 |
| Javascript | No    |           10,000 | 15.1 ms   | 0.7 ms    | 9.4 ms    | 5.5 ms      | 13.5 ms     | 16.7 ms     | 193  |
| Go         | Yes   |          100,000 | 98.2 ms   | 6.7 ms    | 87.7 ms   | 10.1 ms     | 90.1 ms     | 122.1 ms    | 31   |
| Javascript | Yes   |          100,000 | 117.7 ms  | 2.5 ms    | 94.2 ms   | 23.0 ms     | 112.4 ms    | 125.4 ms    | 25   |
| Go         | No    |          100,000 | 0.6 ms    | 0.1 ms    | 0.5 ms    | 0.1 ms      | 0.5 ms      | 1.0 ms      | 4055 |
| Javascript | No    |          100,000 | 18.5 ms   | 1.1 ms    | 12.4 ms   | 5.9 ms      | 15.9 ms     | 22.2 ms     | 179  |
| Go         | Yes   |        1,000,000 | 756.1 ms  | 21.2 ms   | 727.3 ms  | 27.2 ms     | 711.8 ms    | 793.2 ms    | 10   |
| Javascript | Yes   |        1,000,000 | 769.8 ms  | 10.6 ms   | 613.5 ms  | 154.6 ms    | 757.0 ms    | 786.2 ms    | 10   |
| Go         | No    |        1,000,000 | 1.6 ms    | 0.1 ms    | 1.2 ms    | 0.3 ms      | 1.5 ms      | 3.6 ms      | 1734 |
| Javascript | No    |        1,000,000 | 26.2 ms   | 1.6 ms    | 16.7 ms   | 9.2 ms      | 23.1 ms     | 29.6 ms     | 108  |
| Go         | Yes   |       10,000,000 | 7.197 s   | 0.079 s   | 7.134 s   | 0.050 s     | 7.077 s     | 7.296 s     | 10   |
| Javascript | Yes   |       10,000,000 | 8.863 s   | 0.245 s   | 6.089 s   | 2.018 s     | 8.679 s     | 9.442 s     | 10   |
| Go         | No    |       10,000,000 | 13.6 ms   | 0.5 ms    | 6.7 ms    | 6.7 ms      | 11.9 ms     | 15.5 ms     | 214  |
| Javascript | No    |       10,000,000 | 65.5 ms   | 1.4 ms    | 51.8 ms   | 13.2 ms     | 62.8 ms     | 68.7 ms     | 45   |

## Conclusion

This benchmark was very insightful into the performance of Bun and getting deeper than just the standard benchmarks we've grown accustomed to of dependency install times, server startup times, etc.

While Go is most likely going to remain the clear option for those building performance critical network/backend services, it's great to see a new contender revolutionizing the Javascript space with performance that *can* rival performance optimized languages in terms of speed even with the overhead of a runtime.

The biggest takeaway I had, however, is to pre-allocate wherever possible, because who knows how Bun would've performed if I had used a dynamic array with no instantiation.

I'm likely going to do this same benchmark but with Bun vs. Node and update this to see how that turns out as well.

Thank you!