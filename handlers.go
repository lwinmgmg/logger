package logger

func HandleFunc(ch <-chan Message, confirm chan<- struct{}, done <-chan struct{}, writer LoggerWriter, pattern, tFormat string) {
	defer func() {
		confirm <- struct{}{}
	}()
	stopSig := false
	for !stopSig {
		select {
		case mesg, ok := <-ch:
			if !ok {
				stopSig = true
			}
			writer.Write([]byte(mesg.Format(pattern, tFormat)))
		case <-done:
			if len(ch) == 0 {
				stopSig = true
			}
		}
	}
	writer.Close()
}
