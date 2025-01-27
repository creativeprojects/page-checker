package lib

import "context"

func RunWithContext(ctx context.Context, run func() error) error {
	job := make(chan error, 1)
	go func() {
		job <- run()
	}()

	select {
	case err := <-job:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
