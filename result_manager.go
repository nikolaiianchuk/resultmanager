package result_manager

// ResultManager используется при необходимости обрабатывать ответы
// от нескольких Uploader-горутин и возвращать единичный ответ в handler,
// вызвавший обработку
type ResultManager struct {
	JobCount   int
	ReturnOK   chan bool
	ReturnErr  chan error
	ReceiveOK  chan bool
	ReceiveErr chan error
}

// Manage запускает обработку данных из принимающих каналов
func (r *ResultManager) Manage() {
	for i := 0; i < r.JobCount; i++ {
		select {
		case <-r.ReceiveOK:
			continue
		case err := <-r.ReceiveErr:
			r.ReturnErr <- err
			return
		}
	}
	r.ReturnOK <- true
}

// NewResultManager создает инстанс ResultManager и возвращает каналы, через который будут передаваться данные
func NewResultManager(uploadCount int, ReturnOK chan bool, ReturnErr chan error) (chan bool, chan error) {
	b := make(chan chan bool)
	e := make(chan chan error)
	go func(uploadCount int, b chan chan bool, e chan chan error) {
		r := ResultManager{
			JobCount:   uploadCount,
			ReturnOK:   ReturnOK,
			ReturnErr:  ReturnErr,
			ReceiveErr: make(chan error, uploadCount),
			ReceiveOK:  make(chan bool, uploadCount),
		}
		b <- r.ReceiveOK
		e <- r.ReceiveErr
		r.Manage()
	}(uploadCount, b, e)
	return <-b, <-e
}
