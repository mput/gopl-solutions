package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type Node interface{}  // CharData or *Element

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func elString(e *Element, lvl int) string {
	atrs := ""
	for _, atr := range e.Attr {
		atrs += fmt.Sprintf(" %s='%s'", atr.Name.Local, atr.Value)
	}

	s := fmt.Sprintf("%s<%s%s>\n",strings.Repeat(" ", lvl), e.Type.Local, atrs)
	for _, el := range e.Children {
		switch el := el.(type) {
		case string:
			s += strings.Repeat(" ", lvl + 2) + el + "\n"
		case *Element:
			s += strings.Repeat(" ", lvl) + elString(el, lvl + 2)
		}
	}
	s += fmt.Sprintf("%s</%s>\n",strings.Repeat(" ", lvl), e.Type.Local)
	return s
}

func (e *Element) String() string  {
	return elString(e, 0)
}

// RECURSION SOLUTION

func BuildNode(d *xml.Decoder, root *Element) (*Element, error) {
	for {
		token, err := d.Token()
		if err == io.EOF {
			return root, nil
		}
		if err != nil {
			return root, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			nextNodePtr, err := BuildNode(d, &Element{
				Type: token.Name,
				Attr: token.Attr,
			})
			if err != nil {
				return &Element{}, err
			}
			root.Children = append(root.Children, nextNodePtr)
		case xml.CharData:
			s := strings.TrimSpace(string(token))
			if s != "" {
				root.Children = append(root.Children, s)
			}
		case xml.EndElement:
			return root, nil
		}
	}
}

func Parse(r io.Reader) (*Element, error) {
	d := xml.NewDecoder(r)
	rootNode, err := BuildNode(d, &Element{
		Type: xml.Name{Local: "root-service-element"},
		Children: make([]Node, 0),
	})
	if err != nil {
		return rootNode, err
	}

	return extractFirstElmnt(rootNode)
}

func extractFirstElmnt (e *Element) (*Element, error) {
	for _, v := range e.Children {
		v, ok := v.(*Element)
		if ok {
			return v, nil
		}
	}
	return &Element{}, fmt.Errorf("Miss any nodes")
}

// STACK SOLUTION

func ParseStack(r io.Reader) (*Element, error) {
	d := xml.NewDecoder(r)
	st := make(elmStack, 0)
	st.push(&Element{
		Type: xml.Name{Local: "root-service-element"},
	})

	for {
		token, err := d.Token()
		if err == io.EOF {
			if rootElm, ok := st.pop(); ok {
				return extractFirstElmnt(rootElm)
			}
			return &Element{}, fmt.Errorf("wrong xml document structure")
		}
		if err != nil {
			return &Element{}, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			nextNodePtr := &Element{
				Type: token.Name,
				Attr: token.Attr,
			}
			st.addChild(nextNodePtr)
			st.push(nextNodePtr)
		case xml.CharData:
			s := strings.TrimSpace(string(token))
			if s != "" {
				st.addChild(s)
			}
		case xml.EndElement:
			st.pop()
		}
	}
}

type elmStack []*Element

func (st *elmStack) push (el *Element) {
	*st = append(*st, el)
}

func (st *elmStack) pop () (*Element, bool) {
	l := len(*st)
	if l < 1 {
		return &Element{}, false
	}
	elp := (*st)[l-1]
	*st = (*st)[:l-1]
	return elp, true
}

func (st *elmStack) addChild (el Node) {
	l := len(*st)
	(*st)[l-1].Children = append((*st)[l-1].Children, el)
}

func main() {
	// res, err := Parse(os.Stdin)
	var in io.Reader
	in = strings.NewReader(blob)
	// in = os.Stdin
	resReq, err := Parse(in)
	if err != nil {
		fmt.Printf("Error Req Solution: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf(
		"Recursion solution: \n%s\n------------------\n",
		resReq,
	)

	in = strings.NewReader(blob)

	resStack, err := ParseStack(in)
	if err != nil {
		fmt.Printf("Error Stack Solution: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf(
		"Stack solution: \n%s\n------------------\n",
		resStack,
	)
}


var blob string = `
<animals help='4'>
	<animal>gopher</animal>
	<animal>armadillo</animal>
	<animal>zebra</animal>
	<animal>unknown</animal>
	<animal>gopher</animal>
	<animal>bee</animal>
	<animal>gopher</animal>
	<animal>zebra</animal>
</animals>`
