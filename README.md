
# Typed-js (WIP)

## Example

```js
def object window;

// Macros
#if env in development staging
include "./libs/vue.js"
#else
include "./libs/vue.min.js"
#endif

include "./globals.tjs" // Includes other .tjs file

Struct aMessage {
    string message
}

object myButton = document.getElementById("my-button");
func removeElement = function(object el, aMessage msg) void {

    el.parentNode.removeChild(el);

    if(showMessage){
        window.alert(msg.message);
        // Note:
        alert("..."); // Will fail unless you use: def func alert (string) void
    }
}

aMessage msg = {
    message: "Element has been deleted"
};

removeElement(myButton, msg);

include "./window-ready.tjs"
```