
# Typed-js (WIP)

## Example

```jsx

import "./libs/vue.js" // Include plain .js code (no type checking)
import Vue Component:VueComponent "./libs/vue-defs.tjs" // Import Vue & Component struct
// use VueComponent as alias for Component incase we already have a struct named Component

import "./libs/vue-router.js"
import Router from "./libs/vue-router-defs.tjs" pass { VueComponent:Component }
// Because VueRouter is dependent on Component from the Vue Package, we pass it through

import "./some-structs.tjs" // Include other .tjs files (type checked)

struct App {
    loading: bool
    storage: object<string> // { key: string-value }
    request_queue: array<Request>
};

extend struct Window { // by default tjs creates a "Window" struct, you can extend it like this
    app: App
};

struct aMessage {
    string message
}

func removeElement = function(Element el, aMessage? msg) void {

    el.parentNode.removeChild(el);

    if(msg){
        window.alert(msg.message);
        // Note:
        alert("..."); // Will fail unless you use: define func alert (string) void
    }

    // Macros
    #if debug eq 1
    console.log("An element was deleted")
    #end
};

// Element is one of the default structs provided by tjs and follows the definitions found on the MDN website
Element myButton = document.getElementById("my-button");

aMessage msg = {
    message: "Element has been deleted"
};

removeElement(myButton, msg);

// Store a string into a variable, the word "HMTL" is just for the IDE to know it's HTML
#string:html myTemplate
<div>Hello world!</div>
#end // results into: string myTemplate = "<div>Hello world!</div>"; 

include "./window-ready.tjs"

// exporting has no real function except for when u want to share your own structs with others
// by using: tjs compile src/main.tjs dist/my-package.js --export dest/my-package-defs.tjs
export structs App aMessage VueComponent:Component // Rename VueComponent to Component in our export
export vars msg // results in "define aMessage msg;"
```

```
tjs compile src/main.tjs dest/main.js --vars "env=development|debug=1"
```

## Rules

- No union types, types of objects cannot be checked on runtime, so because u cant garantee a type, your code becomes unpredictable.
- No general types such as object or any
