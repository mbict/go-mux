[![wercker status](https://app.wercker.com/status/2e576950a64845c6c52d29511f9c0bac/s "wercker status")](https://app.wercker.com/project/bykey/2e576950a64845c6c52d29511f9c0bac)
[![Build Status](https://travis-ci.org/mbict/go-mux.png?branch=master)](https://travis-ci.org/mbict/go-mux)
[![GoDoc](https://godoc.org/github.com/mbict/go-mux?status.png)](http://godoc.org/github.com/mbict/go-mux)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-mux)](http://gocover.io/github.com/mbict/go-mux)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-mux)](http://goreportcard.com/report/mbict/go-mux)

# Mux Router

A simplified / stripdown version of the golang http.ServeMux.

## Why i created this

I created this simplified version to be able to match on partial paths.
The muxer is mainly used for preselection of handlers based on the beginning of the path.

This is something i do when creating go-kit services.

## How it works

The matcher will start at the longest possible path an tries to (partial) match it with the request path.
If it found a match the http.Handler will be invoked and matching stops. 
When the path pattern does not match it will try the next path pattern.
If none of the path patterns match, the default not found handler will be invoked. 

## Parts i removed

- Redirect explicit slash
- Exact match on routes not ending with a slash
- Hostname matching
- Mutex locks (this muxer is initialized once at startup)