# CookieJar

[![Build Status](https://travis-ci.org/cookiejars/cookiejar.svg?branch=master)](https://travis-ci.org/cookiejars/cookiejar)

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

   _Note: please consider sharing your implementation as an open-source project and help this baby grow
   (check [Contributing](github.com/cookiejars/cookiejar#contributing))._

3. Setup your backoff strategy which defines how the worker pool behaves in the interval of processing work, by either:

   a) Using the `ConstantBackoff` or `ExponentialBackoff` implementations provided.

   b) Creating your own implementation that follows the `Backoff` interface.

4. Initialize the pool:

    ```golang
    digester := cookiejar.NewDigester(workers, jar, backoff)
    ```

5. Start the pool:

    Here you have to pass the function used to process the work, in the form of: `func(cookie cookiejar.Cookie) error`.

    Also, you can pass a list of `os.Signal` which will make the pool wait for to shutdown gracefully.
    If you don't pass any signals the `Start` function will exit immediately, leaving the pool working on the
    background. It will be up to you to define how to wait for the work to be complete and
    call `digester.Stop()` to stop the pool.

    ```golang
    digester.Start(digestFn, signals...)
    ```

## Contributing

This project aims to be generic and fit as much cases as possible. This will only be possible if you share your
specific usecase to help identify where the project is still lacking.

I'd love to have [github.com/cookiejars](https://github.com/cookiejars) as the main place to go to find the existent
implementations, so if you wish to contribute feel free to get in touch with me by opening an issue in the
[CookieJar](https://github.com/cookiejars/cookiejar) project.
