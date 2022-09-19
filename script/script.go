package script

import "context"

func prepareExecutor(filename string, modules ...string) (string, Executor, error) {
	const id = "__default__"
	conf := DefaultConfig()
	conf.Modules = modules
	man := NewManager()
	exe := NewExecutor(man, Unlimited)
	if err := man.Add(id, filename, conf); err != nil {
		return id, nil, err
	}
	return id, exe, nil
}

// RunFile execute tengo script file using default context
func RunFile(filename string, modules ...string) error {
	id, exe, err := prepareExecutor(filename, modules...)
	if err != nil {
		return err
	}
	_, err = exe.Exec(id, nil)
	return err
}

// RunFileContext execute tengo script file with given context
func RunFileContext(ctx context.Context, filename string, modules ...string) error {
	id, exe, err := prepareExecutor(filename, modules...)
	if err != nil {
		return err
	}

	_, err = exe.ExecContext(ctx, id, nil)
	return err
}
