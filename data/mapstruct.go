package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"fyksz/yaml"
	"github.com/pkg/errors"
)

type Visitor interface {
	OnKey(*KeyNode)

	BeforeMap(node *MapNode)
	AfterMap(node *MapNode)
	BeforeMapItem(node *MapNode, key string, index int)
	AfterMapItem(node *MapNode, key string, index int)

	BeforeList(node *ListNode)
	AfterList(node *ListNode)
	BeforeListItem(node *ListNode, item Node, index int)
	AfterListItem(node *ListNode, item Node, index int)
}

type Node interface {
	Accept(Visitor)
	GetPath() Path
}

// ----------------- KEY NODE --------------
type KeyNode struct {
	Value interface{}
	Path  Path
}

func NewKeyNode(path Path, value interface{}) KeyNode {
	return KeyNode{
		Value: value,
		Path:  path,
	}
}

func (node *KeyNode) GetPath() Path {
	return node.Path
}

func (node *KeyNode) Accept(v Visitor) {
	v.OnKey(node)
}

// ----------------- MAP NODE --------------
type MapNode struct {
	keys     []string
	children map[string]Node
	Path     Path
}

func NewMapNode(path Path) MapNode {
	m := MapNode{
		Path:     path,
		children: make(map[string]Node),
	}
	return m
}

func (node *MapNode) GetPath() Path {
	return node.Path
}

func (node *MapNode) Put(key string, value Node) {
	node.children[key] = value
	for _, indexedKey := range node.keys {
		if indexedKey == key {
			//the key is already indexed
			return
		}
	}
	node.keys = append(node.keys, key)
}

func (node *MapNode) Accept(v Visitor) {
	v.BeforeMap(node)
	idx := 0
	for _, key := range node.keys {
		v.BeforeMapItem(node, key, idx)
		idx = idx + 1
		value := node.children[key]
		value.Accept(v)
		v.AfterMapItem(node, key, idx)

	}
	v.AfterMap(node)
}

func (node *MapNode) ToString() (string, error) {
	converted := ConvertToYaml(node)
	bytes, err := yaml.Marshal(converted)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (node MapNode) Get(s string) Node {
	result := node.children[s]
	return result
}

func (node MapNode) HasKey(s string) bool {
	_, ok := node.children[s]
	return ok
}

func (node *MapNode) Len() int {
	if node.children == nil {
		return 0
	} else {
		return len(node.children)
	}
}
func (node *MapNode) Keys() []string {
	return node.keys
}
func (node *MapNode) PutValue(key string, value interface{}) {
	node.Put(key, &KeyNode{Path: node.Path.Extend(key), Value: value})
}
func (node *MapNode) CreateMap(key string) *MapNode {
	mapNode := NewMapNode(node.Path.Extend(key))
	node.Put(key, &mapNode)
	return &mapNode
}

func (node *MapNode) CreateList(key string) *ListNode {
	listNode := NewListNode(node.Path.Extend(key))
	node.Put(key, &listNode)
	return &listNode
}

func (node *MapNode) Remove(key string) {
	delete(node.children, key)
	newKeys := make([]string, 0)
	for _, indexedKey := range node.keys {
		if indexedKey != key {
			newKeys = append(newKeys, indexedKey)
		}
	}
	node.keys = newKeys
}

func (node *MapNode) GetStringValue(s string) string {
	return node.Get(s).(*KeyNode).Value.(string)
}

// ----------------- LIST NODE --------------
type ListNode struct {
	Children []Node
	Path     Path
}

func NewListNode(path Path) ListNode {
	l := ListNode{
		Children: make([]Node, 0),
		Path:     path,
	}
	return l
}

func (node *ListNode) GetPath() Path {
	return node.Path
}

func (node *ListNode) Append(value Node) {
	node.Children = append(node.Children, value)
}

func (node *ListNode) Accept(v Visitor) {
	v.BeforeList(node)
	for ix, value := range node.Children {
		v.BeforeListItem(node, value, ix)
		value.Accept(v)
		v.AfterListItem(node, value, ix)
	}
	v.AfterList(node)

}
func (node *ListNode) Len() int {
	return len(node.Children)
}

func (node *ListNode) AddValue(value string) {
	node.Append(&KeyNode{Path: node.Path.Extend(strconv.Itoa(len(node.Children))), Value: value})
}

func (node *ListNode) CreateMap() *MapNode {
	mapNode := NewMapNode(node.Path.Extend(strconv.Itoa(len(node.Children))))
	node.Children = append(node.Children, &mapNode)
	return &mapNode
}

// ----------------- VISITORS --------------

type PrintVisitor struct {
	DefaultVisitor
}

func (PrintVisitor) OnKey(node *KeyNode) {
	fmt.Printf("%s %s\n", node.Path.ToString(), node.Value)
}
func (PrintVisitor) BeforeMap(node *MapNode) {
	fmt.Printf("%s [map]\n", node.Path.ToString())
}
func (PrintVisitor) BeforeList(node *ListNode) {
	fmt.Printf("%s [list]\n", node.Path.ToString())
}

type DefaultVisitor struct{}

func (DefaultVisitor) OnKey(*KeyNode)                                      {}
func (DefaultVisitor) BeforeMap(node *MapNode)                             {}
func (DefaultVisitor) AfterMap(node *MapNode)                              {}
func (DefaultVisitor) BeforeMapItem(node *MapNode, key string, index int)  {}
func (DefaultVisitor) AfterMapItem(node *MapNode, key string, index int)   {}
func (DefaultVisitor) BeforeList(node *ListNode)                           {}
func (DefaultVisitor) AfterList(node *ListNode)                            {}
func (DefaultVisitor) BeforeListItem(node *ListNode, item Node, index int) {}
func (DefaultVisitor) AfterListItem(node *ListNode, item Node, index int)  {}

type Apply struct {
	DefaultVisitor
	Path     Path
	Function func(interface{}) interface{}
}

func (visitor *Apply) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) {
		node.Value = visitor.Function(node.Value)
	}
}

type GetKeys struct {
	DefaultVisitor
	Result []GetAllResult
}

type GetKeysResult struct {
	Path  Path
	Value Node
}

func (visitor *GetKeys) OnKey(node *KeyNode) {
	visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})
}

type Yamlize struct {
	DefaultVisitor
	Path       Path
	Serialize  bool
	parsed     bool
	parsedPath Path
}

func (visitor *Yamlize) BeforeMapItem(node *MapNode, key string, index int) {
	if !visitor.Serialize {
		//deserialize phase
		if visitor.parsed {
			return
		}
		matchLimited, _ := visitor.Path.MatchLimited(node.Path.Extend(key))
		matchFull := visitor.Path.Match(node.Path.Extend(key))
		if matchLimited || matchFull {
			switch value := node.Get(key).(type) {
			case *KeyNode:
				yamlDoc := yaml.MapSlice{}

				content := value.Value.(string)
				err := yaml.Unmarshal([]byte(content), &yamlDoc)
				if err != nil {
					panic(err)
				}

				newnode, err := ConvertToNode(yamlDoc, node.Path.Extend(key))
				if err != nil {
					panic(err)
				}
				node.Put(key, newnode)
				visitor.parsed = true
				visitor.parsedPath = node.Path.Extend(key)
				break

			}
		}
	} else {
		if node.Path.Extend(key).Equal(visitor.parsedPath) {
			content, err := node.Get(key).(*MapNode).ToString()
			if err != nil {
				panic(err)
			}
			node.Put(key, &KeyNode{content, node.Path.Extend(key)})
		}
	}

}

type FixPath struct {
	DefaultVisitor
	CurrentPath Path
}

func (visitor *FixPath) OnKey(node *KeyNode) {
	node.Path = visitor.CurrentPath
}

func (visitor *FixPath) BeforeMap(node *MapNode) {
	node.Path = visitor.CurrentPath
}
func (visitor *FixPath) AfterMap(node *MapNode) {

}
func (visitor *FixPath) BeforeList(node *ListNode) {
	node.Path = visitor.CurrentPath
}
func (visitor *FixPath) AfterList(node *ListNode) {

}

func (visitor *FixPath) BeforeMapItem(node *MapNode, key string, index int) {
	visitor.CurrentPath = visitor.CurrentPath.Extend(key)
}
func (visitor *FixPath) AfterMapItem(node *MapNode, key string, index int) {
	visitor.CurrentPath = visitor.CurrentPath.Parent()
}
func (visitor *FixPath) BeforeListItem(node *ListNode, item Node, index int) {
	subKeyName := strconv.Itoa(index)
	if mapItem, convertable := item.(*MapNode); convertable {
		if mapItem.HasKey("name") {
			name := mapItem.Get("name").(*KeyNode).Value.(string)
			subKeyName = name
		}
	}
	visitor.CurrentPath = visitor.CurrentPath.Extend(subKeyName)
}

func (visitor *FixPath) AfterListItem(node *ListNode, item Node, index int) {
	visitor.CurrentPath = visitor.CurrentPath.Parent()
}

func NodeFromPathValue(path Path, value interface{}) Node {
	root := NewMapNode(RootPath())
	current := &root
	for _, segment := range path.Parent().segments {
		current = current.CreateMap(segment)
	}
	current.PutValue(path.Last(), value)
	return &root
}

func (m *MapNode) ToMap() map[string]interface{} {
	return ToMap(m).(map[string]interface{})
}

func ToMap(node Node) interface{} {
	switch object := node.(type) {
	case *MapNode:
		result := make(map[string]interface{})
		for key, value := range object.children {
			result[key] = ToMap(value)
		}
		return result
	case *ListNode:
		result := make([]interface{}, 0)
		for _, item := range object.Children {
			result = append(result, ToMap(item))
		}
		return result
	case *KeyNode:
		return object.Value
	}
	return nil
}

func (keyNode *KeyNode) MarshalJSON() ([]byte, error) {
	if keyNode.Value == nil {
		return []byte("null"), nil
	}
	switch value := keyNode.Value.(type) {
	case string:
		safeValue := strings.ReplaceAll(value, "\"", "\\\"")
		safeValue = strings.Trim(safeValue, "\n")
		return []byte("\"" + safeValue + "\""), nil
	case int:
		return []byte(strconv.Itoa(value)), nil
	case bool:
		if value {
			return []byte("true"), nil
		} else {
			return []byte("false"), nil
		}
	default:
		return []byte("\"" + fmt.Sprintf("%T", keyNode.Value) + "\""), nil
	}
}

func (listNode *ListNode) MarshalJSON() ([]byte, error) {
	res := "["
	for i, childNode := range listNode.Children {
		child, err := json.Marshal(childNode)
		if err != nil {
			return []byte{}, errors.Wrap(err, "Can't marshall "+childNode.GetPath().ToString())
		}
		if i > 0 {
			res += ","
		}
		res += string(child)
	}
	res += "]"
	return []byte(res), nil
}

func (mapNode *MapNode) MarshalJSON() ([]byte, error) {
	res := "{"
	for i, key := range mapNode.Keys() {
		child, err := json.Marshal(mapNode.children[key])
		if err != nil {
			return []byte{}, errors.Wrap(err, "Can't marshall "+mapNode.children[key].GetPath().ToString())
		}
		if i > 0 {
			res += ","
		}
		res += "\"" + key + "\":" + string(child)
	}
	res += "}"
	return []byte(res), nil
}
