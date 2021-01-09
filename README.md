
# Typed-js (WIP)

## Full feature list

```jsx

include "./libs/vue.js" // Include plain .js code (no type checking)
include "./libs/vue-defs.tjs" // Import structs from Vue

include "./libs/vue-router.js"
include "./libs/vue-router-defs.tjs"

import Vue Component:VueComponent from Vue // Load structs, load Vue as Vue & Component as VueComponent
import Router from VueRouter

include "./libs/jQuery1.0.tjs"
include "./libs/jQuery2.0.tjs" as jQuery2 // jQuery namespace already exists, rename
include "./libs/jqeury-ajax-lib.tjs" as AjaxLib alias jQuery:jQuery2 // make it use jQuery2 when it imports from jQuery

import Ajax from AjaxLib

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
export namespace MyPlugin // When someone includes your export file, it will store the structs under this name (optional)
export structs App aMessage VueComponent:Component // Rename VueComponent to Component in our export
export vars msg // results in "define aMessage msg;"

// Then someone else can do
include "./libs/my-plugin-defs.tjs"
import aMessage from MyPlugin
```

```
tjs compile src/main.tjs dest/main.js --vars "env=development|debug=1"
```

## Rules

- No union types, types of objects cannot be checked on runtime, so because u cant garantee a type, your code becomes unpredictable.
- No general types such as object or any
