# Concurrent Response Benchmark
Access a URL multiple times at once and outputs statistics about response time.

Written in Go.


## Installation
1. run **go install**

## Usage
crb -url **URL** [-count **COUNT**] [-loops **LOOPS**] [-cooldown **COOLDOWN**] [-verbose]
- **URL** : URL to request
- **COUNT** (optional, default 1) : How many concurrent requests to make at once
- **LOOPS** (optional, default 1) : How many times to repeat the process for extra accuracy
- **COOLDOWN** (optional, default 1000) : How many milliseconds to wait between loops
- **VERBOSE** (optional, default false) : If to display extra information about loops

### Examples
* crb -url "http://localhost"
* crb -url "http://localhost" -count 5 -verbose
* crb -url "http://localhost" -count 10 -loops 5 -cooldown 3000 -verbose
