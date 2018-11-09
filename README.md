# CookieJar

[![Build Status](https://travis-ci.org/cookiejars/cookiejar.svg?branch=master)](https://travis-ci.org/cookiejars/cookiejar)
[![GolangCI](https://golangci.com/badges/github.com/cookiejars/cookiejar.svg)](https://golangci.com/r/github.com/cookiejars/cookiejar)
[![Go Report Card](https://goreportcard.com/badge/github.com/cookiejars/cookiejar)](https://goreportcard.com/report/github.com/cookiejars/cookiejar)

## Purpose

CookieJar is a flexible worker pool that can be adapted to fit (almost) any usecase.

This is accomplished with the use of _Jars_ and _Cookies_. These are simply interfaces that represent a work provider
and the information needed to process that work, respectively.

## Usage

1. Add CookieJar to your project:

   `go get -u github.com/cookiejars/cookiejar`

2. Setup your Jar, by either:

   a) [Checking](https://github.com/cookiejars) if there's already an implementation that fits your usecase.

   b) Implementing your own Jar. For this you'll need to create implementations that fit the Jar and Cookie interfaces.

    ```golang
    type Jar interface {
        Retrieve() ([]Cookie, error) // generate work to distribute amongst the various workers
        Retire(Cookie) error // mark the work as done (e.g., delete a message from a queue after it's been processed)
    }

    type Cookie interface {
        ID() string // unique work identifier
        Content() (interface{}, error) // return metadata needed to process the work
    }
    ```

   _Please consider sharing your implementation as an open-source project and help this baby grow
   (check [Contributing](#contributing))._

3. Initialize the pool:

    ```golang
    digester := cookiejar.NewDigester(jar)
    ```

    Optionally, you can also pass a list of options to tweak how the digester works internally.

    List of options:

    Function | Description | Default behaviour (if not called)
    --- | --- | ---
    `SetWorkers` | Define the number of goroutines processing work concurrently. | Sets to `runtime.NumCPU()`.
    `SetBackoff` | Define how the digester behaves in the interval of processing work. You can use the `ConstantBackoff` or `ExponentialBackoff` implementations provided or use your own implementation of the `Backoff` interface. | Sets to `ConstantBackoff` of 10 seconds.
    `SetInfoLog` | Define info logger where internal information is written to. | No info logs.
    `SetErrorLog` | Define error logger where errors are written to. | No error logs.
    `SetStopSignals` | Define a list of `os.Signal`s for the digester to wait for to call `Stop()` automatically, once you call `Start()`. | `Start()` will exit immediately, leaving the pool working in the background. It will be up to you to define how to wait for the work to be complete and call `digester.Stop()` to stop the pool.

4. Start the pool:

    ```golang
    digester.Start(digestFn)
    ```

    Here you have to pass the function used to process the work, in the form of: `func(cookie cookiejar.Cookie) error`.

5. Stop the pool (optional):

    ```golang
    digester.Stop()
    ```

   If you didn't pass any signals, you'll need to explicitly call the `Stop` function for the pool to gracefully shutdown.

   On the other hand, if you did pass them, this will be automatically called once any of them is triggered.

## Contributing

This project aims to be generic and fit as much cases as possible. This will only be possible if you share your
specific usecase to help identify where the project is still lacking.

To ease discovery, I'd love to have [github.com/cookiejars](https://github.com/cookiejars) as the main place to go to
find the existent implementations, so if you wish to contribute feel free to [open an issue](/issues/new) or
[DM me](https://gophers.slack.com/team/U6FQ0K82K) on the Gophers' Slack.
