# Load Tester

This is a demonstration of the [Coding Challenges Load Tester](https://codingchallenges.fyi/challenges/challenge-load-tester) problem. 

# Step 1
read a URL from the command line and make one request to it.
https://gobyexample.com/command-line-arguments

# Step 2
allow the user to specify how many requests to make as well as the URL to test
https://gobyexample.com/command-line-flags

# Step 3
allow concurrent requests using `-c` flag and report summary of result
https://gobyexample.com/worker-pools

# Step 4
Collect stats for the requests:
- total time
- time to first byte
- time to last byte
- number of responses

https://go.dev/blog/http-tracing
https://pkg.go.dev/net/http/httptrace#example-package
https://gist.github.com/cep21/86ddbaf4e66977fc2b67be84c17989f1

# Step 5
report the stats gathered in step 5 in a summary format that is useful to the user

