
# Typed-js (WIP)

## Example

```jsx

include "./libs/vue.js" // Include plain .js code (no type checking)
include "./libs/vue-defs.tjs" { // Include the tjs definitions for Vue
    // Optional: Rename a struct in case there is already a struct named Component in your code or other library
    Component: VueComponent
    Node: VueNode
}

include "./libs/vue-router.js"
include "./libs/vue-router-defs.tjs" {
    Component: VueComponent // tell it to use VueComponent when looking for the struct Component
}

include "./globals.tjs" // Include other .tjs files

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

func removeElement = function(Element el, aMessage|null msg) void {

    el.parentNode.removeChild(el);

    if(msg){
        window.alert(msg.message);
        // Note:
        alert("..."); // Will fail unless you use: define func alert (string) void
    }

    // Macros
    #if debug eq 1
    console.log("An element was deleted")
    #endif
};

// Element is one of the default structs provided by tjs and follows the definitions found on the MDN website
Element myButton = document.getElementById("my-button");

aMessage msg = {
    message: "Element has been deleted"
};

removeElement(myButton, msg);

// Store a string into a variable, the word "HMTL" is just for the IDE to know it's HTML
#STRING:HTML myTemplate
<div>Hello world!</div>
#ENDSTRING // result in var myTemplate = "...content...";

include "./window-ready.tjs"

// In case u want to export struct & variable definitions, so others can use it in their code
export structs App aMessage VueComponent|Component // Rename VueComponent to Component in our export
export vars msg // results in "define aMessage msg;"
```

```
tjs compile src/main.tjs dest/main.js --vars "env=development|debug=1"
```

## Rules

- No union types, types of objects cannot be checked on runtime, so because u cant garantee a type, your code becomes unpredictable.
- No general types such as object or any
