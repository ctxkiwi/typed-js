
# Typed-js (WIP)

## Example

```js

namespace App;

struct App {
    loading: bool
    storage: object<string> // { key: string-value }
    request_queue: array<Request>
};

extend struct window { // by default tjs creates a "window" struct, you can extend it like this
    app: App
};

// Macros
include "./globals.tjs" // Include other .tjs file (type checked)
include "./libs/vue.js" // Include plain .js code (no type checking)
include "./libs/vue-structs.tjs" // tjs must know the structs from Vue.. so either vue provides this via their github, but most likely you'll have to make it yourself

struct aMessage {
    string message
}

func removeElement = function(object el, aMessage|null msg) void {

    el.parentNode.removeChild(el);

    if(msg){
        window.alert(msg.message);
        // Note:
        alert("..."); // Will fail unless you use: def func alert (string) void
    }

    #if debug eq 1
    console.log("An element was deleted")
    #endif
};

Element myButton = document.getElementById("my-button");

aMessage msg = {
    message: "Element has been deleted"
};

removeElement(myButton, msg);

include "./window-ready.tjs"

export structs App
export vars msg
```

```
tjs compile src/main.tjs lib/main.js --vars "env=development|debug=1"
```

## Rules

- No union types, types of objects cannot be checked on runtime, so because u cant garantee a type, your code becomes unpredictable.
- No general types such as object or any
