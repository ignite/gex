package integration_test

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/require"
)

var (
	// gexApp hold the location of the gex binary used in the integration
	// tests. The binary is compiled the first time the env.New() function is
	// invoked.
	gexApp = path.Join(os.TempDir(), "gex-tests", "gex")

	compileBinaryOnce sync.Once
)

// ensureBinary ensure gex binary was compiled.
func ensureBinary(t *testing.T, ctx context.Context) {
	t.Helper()
	ctx, cancel := context.WithCancel(ctx)

	t.Cleanup(cancel)
	compileBinaryOnce.Do(func() {
		if err := compileBinary(ctx); err != nil {
			panic(err)
		}
	})
}

// compileBinary compile the gex binary.
func compileBinary(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Errorf("unable to get working dir: %v", err)
	}
	pkgs, err := gocmd.List(ctx, wd, []string{"-m", "-f={{.Dir}}", "github.com/ignite/gex"})
	if err != nil {
		return errors.Errorf("unable to list gex package: %v", err)
	}
	if len(pkgs) != 1 {
		return errors.Errorf("expected only one package, got %d", len(pkgs))
	}
	var (
		appPath        = pkgs[0]
		output, binary = filepath.Split(gexApp)
	)

	if err = gocmd.BuildPath(ctx, output, binary, appPath, nil); err != nil {
		return errors.Errorf("error while building binary: %v", err)
	}
	return err
}

func TestGexExplorer(t *testing.T) {
	var (
		require     = require.New(t)
		env         = envtest.New(t)
		app         = env.Scaffold("github.com/ignite/mars")
		servers     = app.RandomizeServerPorts()
		execErr     = &bytes.Buffer{}
		ctx, cancel = context.WithCancel(env.Ctx())
	)
	ensureBinary(t, ctx)

	steps := step.NewSteps(
		step.New(
			step.Stderr(execErr),
			step.Workdir(app.SourcePath()),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),
			step.Exec(gexApp, "explorer", servers.RPC),
			step.InExec(func() error {
				time.Sleep(10 * time.Second)
				cancel()
				return nil
			}),
		),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		env.Must(env.Exec("run gex", steps, envtest.ExecRetry(), envtest.ExecCtx(ctx)))
		wg.Done()
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))
	wg.Wait()

	require.Empty(execErr.String())
}
