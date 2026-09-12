package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

./godelw verify failed after updating godel plugins and assets

Command that caused error:
./godelw lint --fix

Output:
panic: file requires newer Go version go1.27 (application built with go1.26) [recovered, repanicked]

goroutine 1299 [running]:
go/types.(*Checker).handleBailout(0x33cea16bcc00, 0x33cea407f9f0)
	/home/runner/go/go-dists/go1.26.0/src/go/types/check.go:473 +0x91
panic({0x14dd560?, 0x33cea39cb630?})
	/home/runner/go/go-dists/go1.26.0/src/runtime/panic.go:860 +0x13a
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).convertError(0x33ce9956a180, {0x1b10b80?, 0x33ce9c324380})
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:482 +0x6c7
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).loadFromSource.func2({0x1b10b80?, 0x33ce9c324380?})
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:206 +0x37
go/types.(*Checker).handleError(0x33cea16bcc00, 0x0, {0x1b11a60, 0x33cea09494e0}, 0x97, {0x33ce9f708730, 0x45}, 0x0)
	/home/runner/go/go-dists/go1.26.0/src/go/types/errors.go:221 +0x442
go/types.(*error_).report(0x33cea407f770)
	/home/runner/go/go-dists/go1.26.0/src/go/types/errors.go:148 +0x2cc
go/types.(*Checker).errorf(0x33ce9931f6b0?, {0x1b11a60, 0x33cea09494e0}, 0x6?, {0x186de64?, 0x1b11660?}, {0x33cea407f878?, 0x33cea2925ef0?, 0x33cea0dea210?})
	/home/runner/go/go-dists/go1.26.0/src/go/types/errors.go:243 +0x150
go/types.(*Checker).initFiles(0x33cea16bcc00, {0x33ce9f4c99a0, 0x9, 0x28?})
	/home/runner/go/go-dists/go1.26.0/src/go/types/check.go:421 +0x797
go/types.(*Checker).checkFiles(0x33cea16bcc00, {0x33ce9f4c99a0?, 0x6f4b65?, 0x16472a0?})
	/home/runner/go/go-dists/go1.26.0/src/go/types/check.go:522 +0x18d
go/types.(*Checker).Files(0x1594c40?, {0x33ce9f4c99a0?, 0x33cea39fea80?, 0x5?})
	/home/runner/go/go-dists/go1.26.0/src/go/types/check.go:491 +0x75
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).loadFromSource(0x33ce9956a180, 0x2)
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:212 +0x73a
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).loadImportedPackageWithFacts(0x33ce9956a180, 0x2)
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:339 +0xe5
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).loadWithFacts(0x33ce9b997da8?, 0x33ce9b997cf0?)
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:319 +0x128
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).analyze(0x33ce9956a180, {0x1b21288, 0x33ce9aa0b7c0}, 0x33ce9a10b9a0, 0x2, 0x33ce9a95a5b0)
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:78 +0x145
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).analyzeRecursive.func1()
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:61 +0x1ff
sync.(*Once).doSlow(0x0?, 0x0?)
	/home/runner/go/go-dists/go1.26.0/src/sync/once.go:78 +0xac
sync.(*Once).Do(...)
	/home/runner/go/go-dists/go1.26.0/src/sync/once.go:69
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).analyzeRecursive(0x0?, {0x1b21288?, 0x33ce9aa0b7c0?}, 0x0?, 0x0?, 0x0?)
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:45 +0x59
github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).analyzeRecursive.func1.1(0x0?)
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:53 +0x30
created by github.com/golangci/golangci-lint/v2/pkg/goanalysis.(*loadingPackage).analyzeRecursive.func1 in goroutine 1035
	/tmp/custom-gcl3055823771/golangci-lint/pkg/goanalysis/runner_loadingpackage.go:52 +0xcb

*/
