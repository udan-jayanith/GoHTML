## Stage 1
- bug fix: Comments get parsed if html, css and js.
- minor-change: html="true" fix.
- feature add: RemoveAttribute

## Stage 2
QuerySelector takes attribute name and regexp for the value and returns the first node that matches the regexp.  
* QuerySelector

QuerySelectorAll takes two regexps and returns all nodes that matches the regexps in attribute and value in order.
* QuerySelectorAll

Closest returns the closest node that matches the className. 
* Closest

## Stage 3
AddClass add the given class name to the node.
* AddClass

RemoveClass removes the specified class name from the node.
* RemoveClass

HasClass returns a boolean value specifying whether the node has the specified class name or not.
* HasClass

GetClassList returns a map of class names in the specified node.
* GetClassList

## Stage 4
* GetElementById
* GetElementByClassName
* GetElementByTagName
* GetElementsById
* GetElementsByClassName
* GetElementsByTagName