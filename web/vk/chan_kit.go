package vk

type Request struct {
	Name   string
	Params map[string]string
}

type Answer struct {
	Output map[string]interface{}
	Error  error
}

type ChanKit struct {
	RequestChan chan Request
	AnswerChan  chan Answer
}

func (chanKit ChanKit)MakeRequest(name string, params map[string]string) Answer {
	chanKit.RequestChan <- Request{name, params}
	return <- chanKit.AnswerChan
}