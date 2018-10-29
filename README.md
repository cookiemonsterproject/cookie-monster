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
   }

   type Cookie interface {
        Content() (interface{}, error) // provide information needed to process the work
        Done() error // mark the work as done (e.g., delete a message from a queue after it's been processed)
   }
   ```

   _Please consider sharing your implementation as an open-source project and help this baby grow
   (check [Contributing](#contributing))._

3. Setup your backoff strategy which defines how the worker pool behaves in the interval of processing work, by either:

   a) Using the `ConstantBackoff` or `ExponentialBackoff` implementations provided.

   b) Creating your own implementation that implements the `Backoff` interface.

4. Initialize the pool:

    ```golang
    digester := cookiejar.NewDigester(workers, jar, backoff)
    ```

5. Start the pool:

    ```golang
    digester.Start(digestFn, signals...)
    ```

    Here you have to pass the function used to process the work, in the form of: `func(cookie cookiejar.Cookie) error`.

    Also, you can pass a list of `os.Signal` which will make the pool wait for to shutdown gracefully.
    If you don't pass any signals the `Start` function will exit immediately, leaving the pool working in the
    background. It will be up to you to define how to wait for the work to be complete and
    call `digester.Stop()` to stop the pool.

6. Stop the pool (optional):

    ```golang
    digester.Stop()
    ```

   If you didn't pass any signals to the `Start` function, you'll need to explicitly call the `Stop` function for the
   pool to gracefully shutdown.

   On the other hand, if you did pass them, this will be automatically called once any of them is triggered.

## Contributing

This project aims to be generic and fit as much cases as possible. This will only be possible if you share your
specific usecase to help identify where the project is still lacking.

To ease discovery, I'd love to have [github.com/cookiejars](https://github.com/cookiejars) as the main place to go to
find the existent implementations, so if you wish to contribute feel free to [open an issue](/issues/new) or
[DM me](https://gophers.slack.com/team/U6FQ0K82K) on the Gophers' Slack.
