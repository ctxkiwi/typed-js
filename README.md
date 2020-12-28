
# Typed-js (WIP)

## Example

```js
def object window;

// Macros
#if env in development staging
include "./libs/vue.js" // Include plain .js code (no type checking)
#else
include "./libs/vue.min.js"
#endif

include "./globals.tjs" // Include other .tjs file

Struct aMessage {
    string message
}

func removeElement = function(object el, aMessage msg) void {

    el.parentNode.removeChild(el);

    if(showMessage){
        window.alert(msg.message);
        // Note:
        alert("..."); // Will fail unless you use: def func alert (string) void
    }

    #if debug eq 1
    console.log("An element was deleted")
    #endif
}

object myButton = document.getElementById("my-button");

aMessage msg = {
    message: "Element has been deleted"
};

removeElement(myButton, msg);

include "./window-ready.tjs"
```

```
tjs compile src/main.tjs lib/main.js --vars "env=development|debug=1"
```