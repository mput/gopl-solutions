package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type Node interface{}

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func elString(e Element, lvl int) string {
	atrs := ""
	for _, atr := range e.Attr {
		atrs += fmt.Sprintf(" %s='%s'", atr.Name.Local, atr.Value)
	}

	s := fmt.Sprintf("%s<%s%s>\n",strings.Repeat(" ", lvl), e.Type.Local, atrs)
	for _, el := range e.Children {
		switch el := el.(type) {
		case string:
			s += strings.Repeat(" ", lvl + 2) + el + "\n"
		case Element:
			s += strings.Repeat(" ", lvl) + elString(el, lvl + 2)
		}
	}
	s += fmt.Sprintf("%s</%s>\n",strings.Repeat(" ", lvl), e.Type.Local)
	return s
}

func (e Element) String() string  {
	return elString(e, 0)
}

func BuildNode(d *xml.Decoder, root Element) (Element, error) {
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
			nextNode := Element{
				Type: token.Name,
				Attr: token.Attr,
			}
			nextNode, err = BuildNode(d, nextNode)
			if err != nil {
				return Element{}, err
			}
			root.Children = append(root.Children, nextNode)
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

func Parse(r io.Reader) (Element, error) {
	d := xml.NewDecoder(r)
	rootNode := Element{
		Type: xml.Name{Local: "root-service-element"},
		Children: make([]Node, 0),
	}
	rootNode, err := BuildNode(d,rootNode)
	if err != nil {
		return rootNode, err
	}
	for _, v := range rootNode.Children {
		v, ok := v.(Element)
		if ok {
			return v, nil
		}
	}

	return Element{}, fmt.Errorf("Miss any nodes")
}

func main() {
	res, err := Parse(os.Stdin)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	fmt.Println(res.String())
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
