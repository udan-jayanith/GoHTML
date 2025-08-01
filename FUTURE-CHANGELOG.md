## v0.0.0-beta.3
-bug: some css and js code get selected as comments.

## v0.0.0
Delete deletes the branch without connecting sibling nodes.
* Delete

QuerySelector takes attribute name and regexp for the value and returns the first node that matches the regexp.  
* QuerySelector

QuerySelectorAll takes two regexps and returns all nodes that matches the regexps in attribute and value in order.
* QuerySelectorAll

Closest returns the closest node that matches the className. 
* Closest

## v0.0.1
AddClass add the given class name to the node.
* AddClass

RemoveClass removes the specified class name from the node.
* RemoveClass

HasClass returns a boolean value specifying whether the node has the specified class name or not.
* HasClass

GetClassList returns a map of class names in the specified node.
* GetClassList

## v0.0.2
* GetElementById
* GetElementByClassName
* GetElementByTagName
* GetElementsById
* GetElementsByClassName
* GetElementsByTagName

## v0.0.3
DecodeHeader only serializes only up to head. And return a node with only head and it's child nodes.
* DecodeOnly
* DecodeOnlyByClassName
* DecodeHeader