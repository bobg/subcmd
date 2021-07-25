# Subcmd - command-line interfaces with subcommands and flags

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/subcmd.svg)](https://pkg.go.dev/github.com/bobg/subcmd)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/subcmd)](https://goreportcard.com/report/github.com/bobg/subcmd)
![Tests](https://github.com/bobg/subcmd/actions/workflows/go.yml/badge.svg)

This is subcmd,
a Go package that is a layer on top of the standard `flag` package.
It simplifies the creation of command-line interfaces
that require flag parsing and that have subcommands that also require flag parsing.
