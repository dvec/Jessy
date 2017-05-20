package vk

//API request struct
type Request struct {
	Name   string //Method name
	Params map[string]string //Request parameters
}

//API answer struct
type Answer struct {
	Output map[string]interface{} //API response
	Error  error //API error
}

//Struct that allows you to create safety API requests
type ChanKit struct {
	RequestChan chan Request //Chan for API request
	AnswerChan  chan Answer //Chan to get the API response
}

//Function that performs the request
func (chanKit ChanKit) MakeRequest(name string, params map[string]string) Answer {
	chanKit.RequestChan <- Request{name, params} //Making the request
	return <- chanKit.AnswerChan //Getting an answer
}