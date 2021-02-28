
# Typed-js (Shutdown)

In the making of this package it looked like js development was going to stick with webpack and tools alike. But now with the rise of snowpack and vite, for the first time, js development has been improved instead of being mutilated by bad node.js developers. This brings me hope. After testing vite and seeing that the browser support is doing barely ok, i decided to go their route. Not that their idea is that great.. i mean, u still need strange plugins for no reason and there's still alot of bloat in things like vue. But.. im going to give it a shot. And it saves me the trouble of dealing with mindless github complainers. If they fuck it up again.. ill continue working on this.

## Full feature list

```jsx

// Include plain javascript (no type checking)
include "./libs/vue.js"
include "./libs/vue-router.js"

import Vue Component:VueComponent from "./libs/vue.tjs" // include structs, but give Component a new name 
import Router from VueRouter // alternative: link VueRouter to a file in the typedjs.json config 

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

    local struct ResponseData {
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
// Export classes & structs
// Rename VueComponent to Component in our export
export types App aMessage VueComponent:Component
export values msg // results in "define aMessage msg;"
// If your package is based on another package, you may want to import the components from that package
// Instead of exporting them yourself
// This reduces code and makes sure that there are no version problems
// Because you might export structs from jQuery2.0, but the user of your package might have included jQuery1.6
export imports {
    Component Vue from Vue // import Component & Vue from Vue
}

// Then someone else can do
import someClass orSomeStruct from "./libs/my-plugin.tjs"
```

```
tjs compile src/main.tjs dest/main.js --vars "env=development|debug=1"
```

## Example for a package

```jsx
include "./libs/vue-defs.tjs"

import Component from Vue;

class Route {
    path: string
    component: Component

    constructor: func (string path, Component c){
        this.path = path;
        this.component = c;
    }
}

class Router {
    routes: array<Route>

    addRoute: func (string path, Component c) {
        var r = new Route(path, c);
        this.routes.push(r);
    }

    match: func(string path) Route? {
        for(var k in this.routes){
            Route r = this.routes[k];
            if (r.path == path) {
                return r;
            }
        }
        return null;
    }
}

export namespace VueRouter
export types Route Router
```

## Example for a website

```jsx
// main.tjs
include "./libs/vue.js"
include "./libs/vue-defs.tjs"

include "./components/world.tjs"

import Vue from Vue;

#string:html mytemplate
<div>{{ message }} <world></world></div>
#end

Vue app = new Vue({
    el: '#app',
    template: mytemplate,
    data: {
        message: 'Hello'
    }
});

// components/world.tjs
import Vue from Vue;

#string:html mytemplate
<div>world!</div>
#end

Vue.component('world', {
    template: mytemplate
});
```

## Rules

- No union types, types of objects cannot be checked on runtime, so because u cant garantee a type, your code becomes unpredictable.
- No general types such as object or any
