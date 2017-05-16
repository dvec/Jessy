package interception

import (
	"sync"
	"main/web/vk"
)

type Indications struct {
	sync.Mutex
	InterceptedMessage map[int64]chan vk.Message
}

func (indications *Indications) Add(id int64) {
	indications.Lock()
	defer indications.Unlock()
	indications.InterceptedMessage[id] = make(chan vk.Message)
}

func (indications *Indications) Delete(id int64) {
	indications.Lock()
	defer indications.Unlock()
	delete(indications.InterceptedMessage, id)
}

func (indications *Indications) Init() {
	indications.InterceptedMessage = make(map[int64]chan vk.Message)
}
