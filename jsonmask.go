package jsonmask

import (
	"encoding/json"
	"fmt"
)

const (
	eRaw = iota
	eObj
	eAry
	eOther
)

var jsonNull = []byte("null")

// Mask selects the specific parts of an JSON string, according to the mask "fields".
func Mask(doc []byte, fields string) ([]byte, error) {
	sl, err := compile(fields)
	if err != nil {
		return nil, err
	}
	if len(doc) == 0 || len(sl) == 0 {
		return doc, nil
	}

	raw := json.RawMessage(doc)
	src := newLazyNode(&raw)
	dst := newLazyNode(nil)
	if err := copyLazyNode(dst, src, sl); err != nil {
		return nil, err
	}
	return json.Marshal(dst)
}

func newLazyNode(raw *json.RawMessage) *lazyNode {
	return &lazyNode{raw: raw, obj: nil, ary: nil, which: eRaw}
}

type lazyNode struct {
	raw   *json.RawMessage
	obj   *partialObj
	ary   partialArray
	which int
}

func (n *lazyNode) MarshalJSON() ([]byte, error) {
	switch n.which {
	case eObj:
		return json.Marshal(n.obj)
	case eAry:
		return json.Marshal(n.ary)
	default:
		if n.raw != nil {
			return *n.raw, nil
		}
		return jsonNull, nil
	}
}

func (n *lazyNode) UnmarshalJSON(data []byte) error {
	dest := make(json.RawMessage, len(data))
	copy(dest, data)
	n.raw = &dest
	n.which = eRaw
	return nil
}

func (n *lazyNode) unmarshal() error {
	n.which = eOther
	if n.raw == nil {
		return nil
	}
	switch checkWhich(*n.raw) {
	case eObj:
		err := json.Unmarshal(*n.raw, &n.obj)
		if err != nil {
			return err
		}
		n.which = eObj
	case eAry:
		err := json.Unmarshal(*n.raw, &n.ary)
		if err != nil {
			return nil
		}
		n.which = eAry
	}
	return nil
}

type partialArray []*lazyNode

type partialObj struct {
	obj map[string]*lazyNode
}

func (n *partialObj) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.obj)
}

func (n *partialObj) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &n.obj)
}

func checkWhich(buf []byte) int {
	for _, c := range buf {
		switch c {
		case ' ':
		case '\n':
		case '\t':
		case '[':
			return eAry
		case '{':
			return eObj
		default:
			return eOther
		}
	}
	return eOther
}

func copyLazyNode(dst, src *lazyNode, sl selection) error {
	err := src.unmarshal()
	if err != nil {
		return err
	}

	dst.which = src.which
	switch src.which {
	case eObj:
		if len(sl) == 0 {
			dst.obj = src.obj
			return nil
		}

		dst.obj = &partialObj{obj: make(map[string]*lazyNode)}
		for sk, sv := range sl {
			if sk == "*" {
				for nk, nv := range src.obj.obj {
					dst.obj.obj[nk] = newLazyNode(nil)
					if err := copyLazyNode(dst.obj.obj[nk], nv, sv); err != nil {
						return err
					}
				}
			} else if nv := src.obj.obj[sk]; nv != nil {
				dst.obj.obj[sk] = newLazyNode(nil)
				if err := copyLazyNode(dst.obj.obj[sk], nv, sv); err != nil {
					return err
				}
			}
		}
	case eAry:
		if len(sl) == 0 {
			dst.ary = src.ary
			return nil
		}

		dst.ary = make([]*lazyNode, len(sl))
		for i := range src.ary {
			dst.ary[i] = newLazyNode(nil)
			if err := copyLazyNode(dst.ary[i], src.ary[i], sl); err != nil {
				return err
			}
		}
	default:
		if len(sl) == 0 {
			dst.raw = src.raw
			return nil
		}
		return fmt.Errorf("can not select: not a object or array")
	}
	return nil
}
