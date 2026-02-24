package exene

import (
	"fmt"
	"encoding/json"
)

type Html struct {
	id string
	tag string
	attrs map[string]string
	style map[string]string
	text string
	children []*Html
	events []string
}

type HtmlJson struct {
	Id string `json:"id"`
	Tag string `json:"tag"`
	Attrs map[string]string `json:"attrs"`
	Style map[string]string `json:"style"`
	Text string `json:"text"`
	Children []HtmlJson `json:"children"`
	Events []string `json:"events"`
}

func NewHtml(tag string) *Html {
	return &Html{tag:tag}
}

func (h *Html) Id(id string) *Html {
	h.id = id
	return h
}

func (h *Html) Attr(name string, value any) *Html {
	attrs := h.attrs
	if attrs == nil {
		attrs = make(map[string]string)
	}
	attrs[name] = fmt.Sprintf("%v", value)
	h.attrs = attrs
	return h
}

func (h *Html) GetStyle(name string) string {
	return h.style[name]
}

func (h *Html) Style(name string, value any) *Html {
	styles := h.style
	if styles == nil {
		styles = make(map[string]string)
	}
	styles[name] = fmt.Sprintf("%v", value)
	h.style = styles
	return h
}

func (h *Html) Styles(m map[string]string) *Html {
	styles := h.style
	if styles == nil {
		styles = make(map[string]string)
	}
	for k, v := range m {
		styles[k] = v
	}
	h.style = styles
	return h
}

func (h *Html) Text(t string) *Html {
	h.text = t
	return h
}

func (h *Html) Append(h2 *Html) *Html {
	children := h.children
	if children == nil {
		children = make([]*Html, 0, 1)
	}
	h.children = append(children, h2)
	return h
}

func (h *Html) AppendAll(hs []*Html) *Html {
	children := h.children
	if children == nil {
		children = make([]*Html, 0, 1)
	}
	for _, h2 := range hs {
		children = append(children, h2)
	}
	h.children = children
	return h
}

func (h *Html) Event(evt string) *Html {
	events := h.events
	if events == nil {
		events = make([]string, 0, 1)
	}
	h.events = append(events, evt)
	return h
}

func (h *Html) ToExport() HtmlJson {
	children := make([]HtmlJson, len(h.children))
	for i, c := range h.children {
		children[i] = c.ToExport()
	}
	return HtmlJson{h.id, h.tag, h.attrs, h.style, h.text, children, h.events}
}

func (h *Html) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.ToExport())
}
