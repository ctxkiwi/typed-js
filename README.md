
# Typed-js (WIP)

## Full feature list

```jsx

include "./libs/vue.js" // Include plain .js code (no type checking)
include "./libs/vue-defs.tjs" // include structs from Vue (just an example, this file doesnt exists)

include "./libs/vue-router.js"
include "./libs/vue-router-defs.tjs"

import Vue Component:VueComponent from Vue // Load structs, load Vue as Vue & Component as VueComponent
import Router from VueRouter

Router myRouter;

include "./libs/jQuery1.0.tjs"
include "./libs/jQuery2.0.tjs" as jQuery2 // jQuery namespace already exists, rename
include "./libs/jqeury-ajax-lib.tjs" as AjaxLib alias jQuery:jQuery2 // make it use jQuery2 when it imports from jQuery

import Ajax from AjaxLib

import "./some-structs.tjs" // Include other .tjs files (type checked)

struct App {
    loading: bool = false // default: false (optional)
    storage: object<string> // { key: string-value }
    request_queue: array<Request>

    // functions in structs cant have a default value
    // This is to avoid bad practices, use "class" instead
    isLoading: func () bool
};

App app = {
    loading: true,
    isLoading: function(){ return this.isLoading; }
};

// Structs are just a layout for objects, if u want something like an object with functions, u create a class
class App {
    loading: bool = false
    storage: object<string>
    request_queue: array<Request>
    title: string // default = ""

    constructor: func (string title) void {
        this.title = title;
    }

    isLoading: func () bool {
        return this.loading;
    }
}

App app = new App("Hello world");

// Extending an existing struct
extend struct Window { // by default tjs creates a "Window" struct, you can extend it like this
    app: App
};

// window is a default variable and uses the struct Window
window.app = app;


struct aMessage {
    string? message // default value null
}

func removeElement = function(Element el, aMessage msg) void {

    el.parentNode.removeChild(el);

    if not null msg.message { // Now the compiler knows it's not null
        window.alert(msg.message);
        // Note:
        alert("..."); // Will fail unless you use: define func (string) void alert
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

//
func doRequest = function() void {

    struct ResponseData {
        success: bool
    };
    func responseHandler = function(string jsonData) ResponseData {
        // JSON.parse has return type "any" (not recommended)
        // "any" return type must be assigned to a variable with a known type
        // in this case type: ResponseData
        ResponseData responsData = JSON.parse(response.data);
        return responseData;
    }
};

// Store a string into a variable, the word "HMTL" is just for the IDE to know it's HTML
#string:html myTemplate
<div>Hello world!</div>
#end // results into: string myTemplate = "<div>Hello world!</div>"; 

include "./window-ready.tjs"

// exporting has no real function except for when u want to share your own structs/classes/functions/... with others
// by using: tjs compile src/main.tjs dist/my-package.js --export dest/my-package-defs.tjs
export namespace MyPlugin // When someone includes your export file, it will store the structs under this name (optional)
// Export classes & structs
// Rename VueComponent to Component in our export
export types App aMessage VueComponent:Component
export values msg // results in "define aMessage msg;"

// Then someone else can do
include "./libs/my-plugin-defs.tjs"
import aMessage from MyPlugin
alert(msg.message);
```

```
tjs compile src/main.tjs dest/main.js --vars "env=development|debug=1"
```

## Rules

- No union types, types of objects cannot be checked on runtime, so because u cant garantee a type, your code becomes unpredictable.
- No general types such as object or any
