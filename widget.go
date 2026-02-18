package exene

import (
	"encoding/json"
	"fmt"
)


type Widget interface {
	Type () string
	//Dispatcher() Dispatcher
}


var id int = 0

func NewId() int {
	newId := id
	id += 1
	return newId
}

/*
   ************************************************************
   
     Widget library

   ************************************************************
*/


type Button struct {
	Id string `json:"id"`
	Label string `json:"label"`
	Style map[string]string `json:"style"`
	dispatcher Dispatcher
}

func (w *Button) MarshalJSON() ([]byte, error) {
	w2 := struct{Button; Type string `json:"type"`}{*w, w.Type()}
	j, err := json.Marshal(w2)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (w *Button) Type() string {
	return "button"
}

func (w *Button) WithStyle(style string, value string) *Button {
	w.Style[style] = value
	return w
}

func NewButton(d Dispatcher, label string, act func()) *Button {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	style := make(map[string]string)
	eventDispatch := make(chan bool)
	button := &Button{strId, label, style, d}
	go func() {
		for {
			// Also: handle destroy messages?
			select {
			case <- eventDispatch:
				act()
			}
		}
	}()
	d.RegisterEvent(strId, eventDispatch)
	return button
}




type Text struct {
	Id string `json:"id"`
	Text string `json:"text"`
	Style map[string]string `json:"style"`
	dispatcher Dispatcher
}

func (w *Text) MarshalJSON() ([]byte, error) {
	w2 := struct{Text; Type string `json:"type"`}{*w, w.Type()}
	j, err := json.Marshal(w2)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (w *Text) Type() string {
	return "text"
}

func (w *Text) WithStyle(style string, value string) *Text {
	w.Style[style] = value
	return w
}

func (w *Text) UpdateLabel(text string) {
	w.Text = text
	w.dispatcher.PutUpdate() <- map[string]any{"target": w.Id, "type": "update", "text": text}
}

func NewText(d Dispatcher, text string) *Text {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	style := make(map[string]string)
	textWidget := &Text{strId, text, style, d}
	return textWidget
}





type Gap struct {
	Id string `json:"id"`
	Size string `json:"size"`
	Style map[string]string `json:"style"`
	dispatcher Dispatcher
}

func (w *Gap) MarshalJSON() ([]byte, error) {
	w2 := struct{Gap; Type string `json:"type"`}{*w, w.Type()}
	j, err := json.Marshal(w2)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (w *Gap) Type() string {
	return "gap"
}

func (w *Gap) WithStyle(style string, value string) *Gap {
	w.Style[style] = value
	return w
}

func NewGap(d Dispatcher, size string) *Gap {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	style := make(map[string]string)
	gap := &Gap{strId, size, style, d}
	return gap
}



type Box struct {
	Id string `json:"id"`
	Direction string `json:"direction"`
	Widgets []Widget `json:"widgets"`
	Style map[string]string `json:"style"`
	dispatcher Dispatcher
}

func (w *Box) MarshalJSON() ([]byte, error) {
	w2 := struct{Box; Type string `json:"type"`}{*w, w.Type()}
	j, err := json.Marshal(w2)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (w *Box) Type() string {
	return "layout"
}

func (w *Box) WithStyle(style string, value string) *Box {
	w.Style[style] = value
	return w
}

func NewBox(d Dispatcher, direction string, widgets []Widget) *Box {
	id := NewId()
	strId := fmt.Sprintf("%d", id)
	style := make(map[string]string)
	return &Box{strId, direction, widgets, style, d}
}
