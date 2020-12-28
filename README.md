
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

object myButton = document.getElementById("my-button");
func removeElement = function(object el, bool showMessage) void {

    el.parentNode.removeChild(el);

    if(showMessage){
        window.alert("Element has been deleted");
        // Note:
        alert("..."); // Will fail unless you use: def func alert (string) void
    }
}

removeElement(myButton, true);

include "./window-ready.tjs"
```