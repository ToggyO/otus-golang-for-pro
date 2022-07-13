package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, s := range stages {
		if s == nil {
			continue
		}

		if out == nil {
			return nil
		}

		stage := s
		out = stage(handleStage(done, out))
	}

	return out
}

func handleStage(done In, in In) Out {
	tempChan := make(Bi)

	go func(stop In, input In) {
		defer close(tempChan)
		for {
			select {
			case <-stop:
				return
			case item, ok := <-input:
				if !ok {
					return
				}

				// Добавил select после занятия по разбору ДЗ
				select {
				case <-stop:
					return
				case tempChan <- item:
				}
			}
		}
	}(done, in)

	return tempChan
}
