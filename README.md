<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/go-compiles"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

go-compiles
===========
`go-compiles` verifies that all of the go packages that are part of a project compile properly. This is similar to the 
check done by `go build ./...`, but goes further by also verifying that test files (both those that are part of a 
package and those that are part of a `_test` package) also compile and build without errors.
