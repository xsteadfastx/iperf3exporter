# logginghandler

[![Build Status](https://ci.xsfx.dev/api/badges/xsteadfastx/logginghandler/status.svg)](https://ci.xsfx.dev/xsteadfastx/logginghandler)
[![Go Reference](https://pkg.go.dev/badge/go.xsfx.dev/logginghandler.svg)](https://pkg.go.dev/go.xsfx.dev/logginghandler)
[![Go Report Card](https://goreportcard.com/badge/go.xsfx.dev/logginghandler)](https://goreportcard.com/report/go.xsfx.dev/logginghandler)

Just a simple zerolog based request logging http middleware. It also sets a `X-Request-ID` in the request and response headers.

## Install

        go get -v go.xsfx.dev/logginghandler

## Usage

        handler := logginghandler.Handler(http.HandlerFunc(myHandler))
        http.Handle("/", handler)
        log.Fatal().Msg(http.ListenAndServe(":5000", nil).Error())

In other handlers you can access the UUID:

        func anotherHandler(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintf(w, "your uuid is: %s", logginghandler.GetUUID(r))
        }

The already prepared logger is also available:

        l := logginghandler.Logger(r)
        l.Info().Msg("foo bar")
