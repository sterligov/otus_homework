package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		intermediate := make(Bi)

		go func(intermediate Bi, out Out) {
			defer close(intermediate)

			for {
				select {
				case v, ok := <-out:
					if !ok {
						return
					}
					intermediate <- v
				case <-done:
					return
				}
			}
		}(intermediate, out)

		out = stage(intermediate)
	}

	return out
}
